package edge

import (
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/anantadwi13/gorong2/internal/common"
)

type udpEdge struct {
	id         common.Id
	tunnelInfo common.TunnelInfo
	udpAddr    *net.UDPAddr
	udpConn    *net.UDPConn
	eh         EdgeHandler
	edh        EdgeDataHandler
	readChan   chan []byte
	writeChan  chan []byte
	ctx        context.Context
	cancelCtx  context.CancelFunc
	isServed   atomic.Bool
}

func udpDial(edgeId common.Id, tunnelInfo common.TunnelInfo, eh EdgeHandler, edh EdgeDataHandler) error {
	if tunnelInfo == nil {
		return ErrEdgeInvalid
	}

	udpAddr, err := net.ResolveUDPAddr("udp", tunnelInfo.EdgeAgentAddr())
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if edgeId == nil {
		edgeId = common.GenerateId()
	}

	eCtx, eCancelCtx := context.WithCancel(context.Background())
	defer eCancelCtx()

	ue := &udpEdge{
		id:         edgeId,
		tunnelInfo: tunnelInfo,
		udpAddr:    nil,
		udpConn:    conn,
		eh:         eh,
		edh:        edh,
		ctx:        eCtx,
		cancelCtx:  eCancelCtx,
		readChan:   make(chan []byte, 1024),
		writeChan:  make(chan []byte, 1024),
	}

	go func() {
		for {
			buff := make([]byte, 4096)
			n, _, err := conn.ReadFromUDP(buff)
			if err != nil {
				_ = ue.Close()
				return
			}
			go func() {
				err = ue.passRead(buff[:n])
				if err != nil {
					_ = ue.Close()
					return
				}
			}()
		}
	}()

	return ue.Serve()
}

func udpListen(
	ctx context.Context, tunnelInfo common.TunnelInfo, ote OnTunnelEstablished, eh EdgeHandler,
	edh EdgeDataHandler,
) error {
	if ctx == nil {
		return ErrEdgeInvalid
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	if tunnelInfo == nil {
		return ErrEdgeInvalid
	}

	addr, err := net.ResolveUDPAddr("udp", tunnelInfo.EdgeServerAddr())
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		}
	}()

	if ote != nil {
		ote(tunnelInfo)
	}

	var (
		udpAddrToEdge   = make(map[string]*udpEdge)
		udpAddrToEdgeMu sync.RWMutex
	)

	defer func() {
		udpAddrToEdgeMu.RLock()
		defer udpAddrToEdgeMu.RUnlock()
		for key, es := range udpAddrToEdge {
			if es == nil {
				continue
			}
			_ = es.Close()
			delete(udpAddrToEdge, key)
		}
	}()

	for {
		buff := make([]byte, 4096)
		n, udpAddr, err := conn.ReadFromUDP(buff)
		if err != nil {
			return err
		}

		go func() {
			udpAddrToEdgeMu.Lock()
			if e, ok := udpAddrToEdge[udpAddr.String()]; ok && e != nil {
				err := e.passRead(buff[:n])
				if err == nil {
					udpAddrToEdgeMu.Unlock()
					return
				}
				log.Println("new conn err", err)
			}

			log.Println("new conn", udpAddr)

			eCtx, eCancelCtx := context.WithCancel(ctx)
			defer eCancelCtx()

			ue := &udpEdge{
				id:         common.GenerateId(),
				tunnelInfo: tunnelInfo,
				udpAddr:    udpAddr,
				udpConn:    conn,
				eh:         eh,
				edh:        edh,
				ctx:        eCtx,
				cancelCtx:  eCancelCtx,
				readChan:   make(chan []byte, 1024),
				writeChan:  make(chan []byte, 1024),
			}

			udpAddrToEdge[udpAddr.String()] = ue
			udpAddrToEdgeMu.Unlock()

			err := ue.passRead(buff[:n])
			if err != nil {
				log.Println(err)
				return
			}

			err = ue.Serve()
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (u *udpEdge) Id() common.Id {
	return u.id
}

func (u *udpEdge) TunnelInfo() common.TunnelInfo {
	return u.tunnelInfo
}

func (u *udpEdge) SendData(data []byte) error {
	if u.IsClosed() {
		return ErrEdgeClosed
	}
	u.writeChan <- data
	return nil
}

func (u *udpEdge) IsClosed() bool {
	return u.ctx.Err() != nil
}

func (u *udpEdge) Close() error {
	log.Println("close connection udp")
	u.cancelCtx()
	return nil
}

func (u *udpEdge) Serve() error {
	if u.isServed.CompareAndSwap(false, true) == false {
		return ErrIllegalState
	}

	var (
		wg = &sync.WaitGroup{}
	)

	wg.Add(2)
	u.safeEdgeHandler().OnEdgeCreated(u)
	defer func() {
		u.safeEdgeHandler().OnEdgeClosed(u)
		close(u.writeChan)
	}()

	// Reader
	go func() {
		defer wg.Done()

		for {
			select {
			case <-u.ctx.Done():
				return
			case data, ok := <-u.readChan:
				if !ok {
					u.cancelCtx()
					return
				}
				if data == nil {
					continue
				}
				u.safeEdgeDataHandler().OnEdgeDataReceived(u.id, data)
			}
		}
	}()

	// Writer
	go func() {
		defer wg.Done()

		for {
			select {
			case <-u.ctx.Done():
				return
			case data, ok := <-u.writeChan:
				if !ok {
					u.cancelCtx()
					return
				}
				if data == nil {
					continue
				}

				var (
					n   int
					err error
				)

				if u.udpAddr == nil {
					n, err = u.udpConn.Write(data)
					if err != nil {
						log.Println(err)
						u.cancelCtx()
						return
					}
				} else {
					n, err = u.udpConn.WriteToUDP(data, u.udpAddr)
					if err != nil {
						log.Println(err)
						u.cancelCtx()
						return
					}
				}

				if n != len(data) {
					// todo log this
				}
			}
		}
	}()

	wg.Wait()
	return nil
}

func (u *udpEdge) passRead(data []byte) error {
	if u.IsClosed() {
		return ErrEdgeClosed
	}
	u.readChan <- data
	return nil
}

func (u *udpEdge) OnEdgeCreated(edge Edge) {
	// do nothing, to support safeEdgeHandler
}

func (u *udpEdge) OnEdgeClosed(edge Edge) {
	// do nothing, to support safeEdgeHandler
}

func (u *udpEdge) OnEdgeDataReceived(edgeId common.Id, data []byte) {
	// do nothing, to support safeEdgeDataHandler
}

func (u *udpEdge) safeEdgeDataHandler() EdgeDataHandler {
	if u.edh == nil {
		return u
	}
	return u.edh
}

func (u *udpEdge) safeEdgeHandler() EdgeHandler {
	if u.eh == nil {
		return u
	}
	return u.eh
}
