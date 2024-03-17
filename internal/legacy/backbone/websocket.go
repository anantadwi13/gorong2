package backbone

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anantadwi13/gorong2/internal/legacy/common"
	"github.com/gorilla/websocket"
)

const HeaderToken = "x-token"

type messageConnect struct {
	Id string `json:"id"`
}

type messagePacket struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
	Data []byte `json:"data"`
}

func wsDial(address BackboneAddress, bh BackboneHandler, ph PacketHandler) error {
	addr, ok := address.(*wsAddress)
	if !ok {
		return ErrAddressInvalid
	}

	tlsConfig := &tls.Config{}

	certPools, err := x509.SystemCertPool()
	if err != nil {
		return err
	}

	if addr.rootCAPath != "" {
		rootPem, err := os.ReadFile(addr.rootCAPath)
		if err != nil {
			return err
		}
		ok := certPools.AppendCertsFromPEM(rootPem)
		if !ok {
			return ErrCertRootInvalid
		}
		tlsConfig.RootCAs = certPools
	}

	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  tlsConfig,
	}

	header := http.Header{}
	header.Set(HeaderToken, addr.token)

	conn, _, err := dialer.Dial(addr.String(), header)
	if err != nil {
		return err
	}
	defer conn.Close()

	bbCtx, bbCancelCtx := context.WithCancel(context.Background())
	defer bbCancelCtx()

	wsBb := &wsBackbone{
		id:        common.GenerateId(),
		wsConn:    conn,
		bh:        bh,
		ph:        ph,
		ctx:       bbCtx,
		cancelCtx: bbCancelCtx,
		// todo set max write buffer
		writeChan: make(chan Packet, 1024),
	}

	// initialize connection

	err = conn.WriteJSON(&messageConnect{Id: wsBb.id.String()})
	if err != nil {
		return err
	}

	log.Println("backbone Id", wsBb.Id())

	return wsBb.Serve()
}

func wsListen(address BackboneAddress, bh BackboneHandler, ph PacketHandler) error {
	addr, ok := address.(*wsAddress)
	if !ok {
		return ErrAddressInvalid
	}

	tlsConfig := &tls.Config{}

	certPools, err := x509.SystemCertPool()
	if err != nil {
		return err
	}

	if addr.rootCAPath != "" {
		rootPem, err := os.ReadFile(addr.rootCAPath)
		if err != nil {
			return err
		}
		ok := certPools.AppendCertsFromPEM(rootPem)
		if !ok {
			return ErrCertRootInvalid
		}
		tlsConfig.RootCAs = certPools
	}

	// todo change buffer size
	upgrader := websocket.Upgrader{ReadBufferSize: 1 << 20, WriteBufferSize: 1 << 20}

	http.HandleFunc(addr.path, func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HeaderToken)
		if token != addr.token {
			w.WriteHeader(http.StatusOK)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// initialize connection
		msgConnect := &messageConnect{}
		err = conn.ReadJSON(&msgConnect)
		if err != nil {
			log.Println(err)
			return
		}

		id, err := common.ParseId(msgConnect.Id)
		if err != nil {
			log.Println(err)
			return
		}

		bbCtx, bbCancelCtx := context.WithCancel(context.Background())
		defer bbCancelCtx()

		wsBb := &wsBackbone{
			id:        id,
			wsConn:    conn,
			bh:        bh,
			ph:        ph,
			ctx:       bbCtx,
			cancelCtx: bbCancelCtx,
			// todo set max write buffer
			writeChan: make(chan Packet, 1024),
		}

		log.Println("backbone Id", wsBb.Id())

		err = wsBb.Serve()
		if err != nil {
			log.Println(err)
			return
		}
	})

	server := &http.Server{
		Addr:      addr.networkAddress,
		TLSConfig: tlsConfig,
	}

	if addr.pubKeyPath != "" && addr.privKeyPath != "" {
		err = server.ListenAndServeTLS(addr.pubKeyPath, addr.privKeyPath)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		return err
	}

	return nil
}

type wsAddress struct {
	token          string
	networkAddress string
	secured        bool
	path           string
	rootCAPath     string
	privKeyPath    string
	pubKeyPath     string
}

func NewWebsocketAddress(
	token string, networkAddress, path string, secured bool, rootCAPath, privKeyPath, pubKeyPath string,
) BackboneAddress {
	if path == "" {
		path = "/backbone"
	}
	return &wsAddress{
		token:          token,
		networkAddress: networkAddress,
		secured:        secured,
		path:           path,
		rootCAPath:     rootCAPath,
		privKeyPath:    privKeyPath,
		pubKeyPath:     pubKeyPath,
	}
}

func (w *wsAddress) Protocol() Protocol {
	return ProtocolWebsocket
}

func (w *wsAddress) String() string {
	backboneURL := url.URL{}

	backboneURL.Host = w.networkAddress
	backboneURL.Scheme = "ws"
	if w.secured {
		backboneURL.Scheme = "wss"
	}
	backboneURL.Path = strings.Trim(w.path, "/")

	return backboneURL.String()
}

type wsBackbone struct {
	id        common.Id
	wsConn    *websocket.Conn
	bh        BackboneHandler
	ph        PacketHandler
	writeChan chan Packet
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func (w *wsBackbone) Id() common.Id {
	return w.id
}

func (w *wsBackbone) Serve() error {
	var (
		wg = &sync.WaitGroup{}
	)

	wg.Add(2)
	w.safeBackboneHandler().OnCreate(w)
	defer func() {
		w.safeBackboneHandler().OnClosed(w)
		close(w.writeChan)
	}()

	// Reader
	go func() {
		defer wg.Done()

		for {
			if w.ctx.Err() != nil {
				return
			}
			_, r, err := w.wsConn.NextReader()
			if err != nil {
				w.cancelCtx()
				return
			}

			msgPacket := &messagePacket{}
			err = json.NewDecoder(r).Decode(&msgPacket)
			if err != nil {
				// todo log this
				continue
			}

			id, err := common.ParseId(msgPacket.Id)
			if err != nil {
				// todo log this
				return
			}

			packet, err := NewPacket(id, PacketType(msgPacket.Type), msgPacket.Data)
			if err != nil {
				// todo log this
				return
			}

			w.safePacketHandler().OnReceive(w.Id(), packet)
		}
	}()

	// Writer
	go func() {
		defer wg.Done()

		for {
			select {
			case <-w.ctx.Done():
				return
			case p, ok := <-w.writeChan:
				if !ok {
					// channel is closed
					w.cancelCtx()
					return
				}
				if p == nil {
					// todo log this
					continue
				}

				writer, err := w.wsConn.NextWriter(websocket.BinaryMessage)
				if err != nil {
					w.cancelCtx()
					return
				}

				msgPacket := &messagePacket{
					Id:   p.Id().String(),
					Type: int(p.Type()),
					Data: p.Data(),
				}

				err = json.NewEncoder(writer).Encode(&msgPacket)
				if err != nil {
					// todo log this
					_ = writer.Close()
					continue
				}
				_ = writer.Close()
			}
		}
	}()

	wg.Wait()
	return nil
}

func (w *wsBackbone) Send(packet Packet) error {
	if w.IsClosed() {
		return ErrBackboneClosed
	}
	w.writeChan <- packet
	return nil
}

func (w *wsBackbone) IsClosed() bool {
	return w.ctx.Err() != nil
}

func (w *wsBackbone) Close() error {
	w.cancelCtx()
	_ = w.wsConn.Close()
	return nil
}

func (w *wsBackbone) OnCreate(b Backbone) {
	// do nothing, to support safeBackboneHandler
}

func (w *wsBackbone) OnClosed(b Backbone) {
	// do nothing, to support safeBackboneHandler
}

func (w *wsBackbone) OnReceive(backboneId common.Id, packet Packet) {
	// do nothing, to support safePacketHandler
}

func (w *wsBackbone) safeBackboneHandler() BackboneHandler {
	if w.bh == nil {
		return w
	}
	return w.bh
}

func (w *wsBackbone) safePacketHandler() PacketHandler {
	if w.ph == nil {
		return w
	}
	return w.ph
}
