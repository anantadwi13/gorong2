package agent

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/anantadwi13/gorong2/internal/legacy/backbone"
	"github.com/anantadwi13/gorong2/internal/legacy/common"
	"github.com/anantadwi13/gorong2/internal/legacy/edge"
)

type Config struct {
	Token              string // todo change it
	BackboneAddr       string
	BackboneSecured    bool
	BackboneRootCAPath string
	TunnelId           string
	TunnelType         common.TunnelType
	EdgeServerAddr     string
	EdgeAgentAddr      string
}

type agent struct {
	config               Config
	errChan              chan error
	br                   backbone.Registry
	er                   edge.Registry
	tunnelInfoRegistry   map[string]common.TunnelInfo
	tunnelInfoRegistryMu sync.RWMutex
	edgeToBackbone       map[string]common.Id
	edgeToBackboneMu     sync.RWMutex
}

func NewAgent(config Config) (Agent, error) {
	if config.BackboneAddr == "" || config.EdgeServerAddr == "" || config.EdgeAgentAddr == "" {
		return nil, ErrConfigNotValid
	}

	br, err := backbone.NewRegistry()
	if err != nil {
		return nil, err
	}

	er, err := edge.NewRegistry()
	if err != nil {
		return nil, err
	}

	return &agent{
		config:             config,
		errChan:            make(chan error, 1),
		br:                 br,
		er:                 er,
		tunnelInfoRegistry: map[string]common.TunnelInfo{},
		edgeToBackbone:     map[string]common.Id{},
	}, nil
}

func (a *agent) Start() error {
	defer func() {
		close(a.errChan)
	}()

	wg := sync.WaitGroup{}

	err := a.br.SetAddress(backbone.NewWebsocketAddress(a.config.Token, a.config.BackboneAddr, "", a.config.BackboneSecured, a.config.BackboneRootCAPath, "", ""))
	if err != nil {
		return err
	}

	err = a.br.SetBackboneHandler(a)
	if err != nil {
		return err
	}

	err = a.br.SetPacketHandler(a)
	if err != nil {
		return err
	}

	err = a.er.RegisterEdgeHandler(a)
	if err != nil {
		return err
	}

	err = a.er.RegisterEdgeDataHandler(a)
	if err != nil {
		return err
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		err := a.br.Dial()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	wg.Wait()

	select {
	case err, ok := <-a.errChan:
		if ok && err != nil {
			return err
		}
	default:
	}

	return nil
}

func (a *agent) Shutdown(ctx context.Context) error {
	_ = a.br.Shutdown(ctx)
	_ = a.er.Shutdown(ctx)
	return nil
}

func (a *agent) fatal(err error) {
	a.errChan <- err
	_ = a.Shutdown(context.Background())
}

func (a *agent) OnCreate(b backbone.Backbone) {
	var (
		tunnelId common.Id
		err      error
	)

	if a.config.TunnelId != "" {
		tunnelId, err = common.ParseId(a.config.TunnelId)
		if err != nil {
			a.fatal(err)
			return
		}
	}

	tunnelInfo, err := common.GenerateTunnelInfo(tunnelId, a.config.TunnelType, a.config.EdgeServerAddr, a.config.EdgeAgentAddr)
	if err != nil {
		a.fatal(err)
		return
	}
	rawTunnelInfo, err := common.MarshalTunnelInfo(tunnelInfo)
	if err != nil {
		a.fatal(err)
		return
	}
	packet, err := backbone.NewPacket(common.GenerateId(), backbone.PacketTypeOpen, rawTunnelInfo)
	if err != nil {
		a.fatal(err)
		return
	}
	err = b.Send(packet)
	if err != nil {
		a.fatal(err)
		return
	}
}

func (a *agent) OnClosed(b backbone.Backbone) {
	// todo
	log.Println("connection closed", b.Id())
}

func (a *agent) OnReceive(backboneId common.Id, packet backbone.Packet) {
	switch packet.Type() {
	case backbone.PacketTypeOpenAckSuccess:
		bb, err := a.br.GetById(backboneId)
		if err != nil {
			return
		}

		tunnelInfo, err := common.UnmarshalTunnelInfo(packet.Data())
		if err != nil {
			return
		}

		a.tunnelInfoRegistryMu.Lock()
		a.tunnelInfoRegistry[tunnelInfo.Id().String()] = tunnelInfo
		a.tunnelInfoRegistryMu.Unlock()

		log.Println("tunnel Id", tunnelInfo.Id().String())
		log.Println("connection established", bb.Id())
	case backbone.PacketTypeOpenAckFailed:
		a.fatal(errors.New(string(packet.Data())))
		return
	case backbone.PacketTypeEdgeCreate:
		sendErrResponse := func(err error) {
			bb, err := a.br.GetById(backboneId)
			if err != nil {
				log.Println(err)
				return
			}

			ackPacket, err := backbone.NewPacket(nil, backbone.PacketTypeEdgeCreateAckFailed, packet.Data())
			if err != nil {
				log.Println(err)
				return
			}

			err = bb.Send(ackPacket)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(err)
		}

		data := strings.Split(string(packet.Data()), ",")
		if len(data) != 2 {
			sendErrResponse(backbone.ErrPacketInvalid)
			return
		}
		edgeId, err := common.ParseId(data[0])
		if err != nil {
			sendErrResponse(err)
			return
		}
		tunnelInfoId, err := common.ParseId(data[1])
		if err != nil {
			sendErrResponse(err)
			return
		}
		tunnelInfo, err := a.getTunnelInfo(tunnelInfoId)
		if err != nil {
			sendErrResponse(err)
			return
		}

		go func() {
			a.edgeToBackboneMu.Lock()
			a.edgeToBackbone[edgeId.String()] = backboneId
			a.edgeToBackboneMu.Unlock()

			err = a.er.Dial(edgeId, tunnelInfo)
			if err != nil {
				a.edgeToBackboneMu.Lock()
				delete(a.edgeToBackbone, edgeId.String())
				a.edgeToBackboneMu.Unlock()
				sendErrResponse(backbone.ErrPacketInvalid)
				return
			}
		}()
	case backbone.PacketTypeEdgeClose:
		data := strings.Split(string(packet.Data()), ",")
		if len(data) != 2 {
			return
		}
		edgeId, err := common.ParseId(data[0])
		if err != nil {
			return
		}
		ea, err := a.er.GetById(edgeId)
		if err != nil {
			return
		}
		_ = ea.Close()
		a.edgeToBackboneMu.Lock()
		delete(a.edgeToBackbone, edgeId.String())
		a.edgeToBackboneMu.Unlock()
	case backbone.PacketTypeMessage:
		packetData := bytes.SplitN(packet.Data(), []byte(","), 2)
		if len(packetData) != 2 {
			// drop packet
			return
		}
		edgeId, err := common.ParseId(string(packetData[0]))
		if err != nil {
			return
		}
		err = a.er.SendData(edgeId, packetData[1])
		if err != nil {
			return
		}
	}
}

func (a *agent) getTunnelInfo(tunnelInfoId common.Id) (common.TunnelInfo, error) {
	a.tunnelInfoRegistryMu.RLock()
	defer a.tunnelInfoRegistryMu.RUnlock()
	ti, found := a.tunnelInfoRegistry[tunnelInfoId.String()]
	if !found || ti == nil {
		return nil, ErrTunnelInfoNotFound
	}
	return ti, nil
}

func (a *agent) getBackboneByEdgeId(edgeId common.Id) (backbone.Backbone, error) {
	a.edgeToBackboneMu.RLock()
	a.edgeToBackboneMu.RUnlock()
	bbId, found := a.edgeToBackbone[edgeId.String()]
	if !found || bbId == nil {
		return nil, ErrEdgeNotRegistered
	}
	bb, err := a.br.GetById(bbId)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func (a *agent) OnEdgeCreated(edge edge.Edge) {
	tunnelInfoId := edge.TunnelInfo().Id()

	bb, err := a.getBackboneByEdgeId(edge.Id())
	if err != nil {
		return
	}

	ackPacket, err := backbone.NewPacket(nil, backbone.PacketTypeEdgeCreateAckSuccess, []byte(fmt.Sprintf("%v,%v", edge.Id().String(), tunnelInfoId.String())))
	if err != nil {
		log.Println(err)
		return
	}

	err = bb.Send(ackPacket)
	if err != nil {
		log.Println(err)
		return
	}
}

func (a *agent) OnEdgeClosed(edge edge.Edge) {
	defer log.Println("on edge closed", edge.Id())

	bb, err := a.getBackboneByEdgeId(edge.Id())
	if err != nil {
		return
	}

	packet, err := backbone.NewPacket(nil, backbone.PacketTypeEdgeClose, []byte(fmt.Sprintf("%v,%v", edge.Id().String(), edge.TunnelInfo().Id().String())))
	if err != nil {
		_ = edge.Close()
		return
	}

	err = bb.Send(packet)
	if err != nil {
		_ = edge.Close()
		return
	}

	a.edgeToBackboneMu.Lock()
	delete(a.edgeToBackbone, edge.Id().String())
	a.edgeToBackboneMu.Unlock()

	log.Println("closing edge", edge.Id())
}

func (a *agent) OnEdgeDataReceived(edgeId common.Id, data []byte) {
	es, err := a.er.GetById(edgeId)
	if err != nil {
		return
	}
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	for {
		bb, err := a.getBackboneByEdgeId(edgeId)
		if err == nil {
			if bb.IsClosed() {
				// todo handle rerouting
				_ = es.Close()
				return
			}
			packetData := make([]byte, len(edgeId.String())+len(data)+1)
			copy(packetData, edgeId.String()+",")
			copy(packetData[len(edgeId.String())+1:], data)
			packet, err := backbone.NewPacket(nil, backbone.PacketTypeMessage, packetData)
			if err != nil {
				return
			}
			err = bb.Send(packet)
			if err != nil {
				// todo handle rerouting
				_ = es.Close()
				return
			}
			return
		}
		// todo handle rerouting if backbone is nil

		select {
		case <-timer.C:
			_ = es.Close()
			return
		case <-ticker.C:
		}
	}
}
