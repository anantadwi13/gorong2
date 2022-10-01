package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/anantadwi13/gorong2/internal/backbone"
	"github.com/anantadwi13/gorong2/internal/common"
	"github.com/anantadwi13/gorong2/internal/edge"
)

type edgeBackbone struct {
	bbId    common.Id
	isReady chan struct{}
}

type Config struct {
	Token                  string // todo change it
	BackboneAddr           string
	BackboneRootCAPath     string
	BackbonePublicKeyPath  string
	BackbonePrivateKeyPath string
}

type server struct {
	config             Config
	errChan            chan error
	br                 backbone.Registry
	er                 edge.Registry
	tunnelToBackboneLB int
	tunnelToBackbone   map[string][]common.Id
	tunnelToBackboneMu sync.RWMutex
	edgeToBackbone     map[string]*edgeBackbone
	edgeToBackboneMu   sync.RWMutex
	bbEdges            map[string][]common.Id
	bbTunnels          map[string][]common.Id
	bbMu               sync.Mutex
}

func NewServer(config Config) (Server, error) {
	if config.BackboneAddr == "" {
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
	return &server{
		config:           config,
		errChan:          make(chan error, 1),
		br:               br,
		er:               er,
		tunnelToBackbone: map[string][]common.Id{},
		edgeToBackbone:   map[string]*edgeBackbone{},
		bbEdges:          map[string][]common.Id{},
		bbTunnels:        map[string][]common.Id{},
	}, nil
}

func (s *server) Start() error {
	defer func() {
		close(s.errChan)
	}()

	var (
		err error
		wg  sync.WaitGroup
	)

	if s.config.Token == "" {
		s.config.Token = common.GenerateId().String()
	}

	log.Println("Token", s.config.Token)

	err = s.br.SetAddress(backbone.NewWebsocketAddress(s.config.Token, s.config.BackboneAddr, "", false, s.config.BackboneRootCAPath, s.config.BackbonePrivateKeyPath, s.config.BackbonePublicKeyPath))
	if err != nil {
		return err
	}

	err = s.br.SetBackboneHandler(s)
	if err != nil {
		return err
	}

	err = s.br.SetPacketHandler(s)
	if err != nil {
		return err
	}

	err = s.er.RegisterEdgeHandler(s)
	if err != nil {
		return err
	}

	err = s.er.RegisterEdgeDataHandler(s)
	if err != nil {
		return err
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		err := s.br.Listen()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	wg.Wait()

	select {
	case err, ok := <-s.errChan:
		if ok && err != nil {
			return err
		}
	default:
	}

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	_ = s.br.Shutdown(ctx)
	_ = s.er.Shutdown(ctx)
	return nil
}

func (s *server) fatal(err error) {
	s.errChan <- err
	_ = s.Shutdown(context.Background())
}

func (s *server) OnCreate(b backbone.Backbone) {
	// do nothing
}

func (s *server) OnClosed(b backbone.Backbone) {
	// todo
	s.unregisterBackbone(b.Id())
	log.Println("connection closed", b.Id())
}

func (s *server) OnReceive(backboneId common.Id, packet backbone.Packet) {
	switch packet.Type() {
	case backbone.PacketTypeOpen:
		// todo handle authentication
		sendResponse := func(err error) {
			log.Println(err, string(packet.Data()))
			var (
				packetType = backbone.PacketTypeOpenAckSuccess
				data       = packet.Data()
			)
			if err != nil {
				packetType = backbone.PacketTypeOpenAckFailed
				data = []byte(err.Error())
			}
			bb, err := s.br.GetById(backboneId)
			if err != nil {
				return
			}

			ackPacket, err := backbone.NewPacket(packet.Id(), packetType, data)
			if err != nil {
				return
			}
			err = bb.Send(ackPacket)
			if err != nil {
				return
			}
		}
		tunnelInfo, err := common.UnmarshalTunnelInfo(packet.Data())
		if err != nil {
			sendResponse(err)
			return
		}
		s.tunnelToBackboneMu.RLock()
		if _, found := s.tunnelToBackbone[tunnelInfo.Id().String()]; found {
			s.tunnelToBackboneMu.RUnlock()
			s.appendTunnelAndBackbone(tunnelInfo.Id(), backboneId)
			sendResponse(nil)
			return
		}
		s.tunnelToBackboneMu.RUnlock()
		go func() {
			err := s.er.Listen(tunnelInfo, func(tunnelInfo common.TunnelInfo) {
				s.appendTunnelAndBackbone(tunnelInfo.Id(), backboneId)

				sendResponse(nil)
			})
			if err != nil {
				sendResponse(err)
				return
			}
		}()
		return
	case backbone.PacketTypeEdgeCreateAckSuccess:
		data := strings.Split(string(packet.Data()), ",")
		if len(data) != 2 {
			return
		}
		edgeId, err := common.ParseId(data[0])
		if err != nil {
			return
		}
		es, err := s.er.GetById(edgeId)
		if err != nil {
			return
		}
		edgeBB, err := s.getEdgeBackbone(es.Id())
		if err != nil {
			return
		}
		close(edgeBB.isReady)
		log.Println("edge creation success")
	case backbone.PacketTypeEdgeCreateAckFailed:
		log.Println("edge creation fail")
		edgeId, err := common.ParseId(string(packet.Data()))
		if err != nil {
			return
		}
		defer s.unregisterEdge(edgeId)
		es, err := s.er.GetById(edgeId)
		if err != nil {
			return
		}
		_ = es.Close()
	case backbone.PacketTypeEdgeClose:
		data := strings.Split(string(packet.Data()), ",")
		if len(data) != 2 {
			return
		}
		edgeId, err := common.ParseId(data[0])
		if err != nil {
			return
		}
		ea, err := s.er.GetById(edgeId)
		if err != nil {
			return
		}
		_ = ea.Close()
		s.unregisterEdge(edgeId)
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
		err = s.er.SendData(edgeId, packetData[1])
		if err != nil {
			return
		}
	}
}

func (s *server) OnEdgeCreated(edge edge.Edge) {
	defer log.Println("on edge created", edge.Id(), edge.TunnelInfo())

	ti := edge.TunnelInfo()
	if ti == nil {
		_ = edge.Close()
		return
	}

	bb, err := s.getBackboneByTunnelInfoId(ti.Id())
	if err != nil {
		_ = edge.Close()
		return
	}

	if edgeBB, err := s.getEdgeBackbone(edge.Id()); err == nil && edgeBB != nil {
		if bbTemp, err := s.er.GetById(edgeBB.bbId); err == nil && bbTemp != nil {
			// edge was established before
			return
		}
	}

	packet, err := backbone.NewPacket(nil, backbone.PacketTypeEdgeCreate, []byte(fmt.Sprintf("%v,%v", edge.Id().String(), edge.TunnelInfo().Id().String())))
	if err != nil {
		_ = edge.Close()
		return
	}

	err = bb.Send(packet)
	if err != nil {
		_ = edge.Close()
		return
	}

	s.registerEdgeToBackbone(edge.Id(), bb.Id())

	log.Println("creating edge", edge.Id())
}

func (s *server) OnEdgeClosed(edge edge.Edge) {
	defer log.Println("on edge closed", edge.Id())

	edgeBB, err := s.getEdgeBackbone(edge.Id())
	if err != nil {
		return
	}

	bb, err := s.br.GetById(edgeBB.bbId)
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

	s.unregisterEdge(edge.Id())

	log.Println("closing edge", edge.Id())
}

func (s *server) OnEdgeDataReceived(edgeId common.Id, data []byte) {
	es, err := s.er.GetById(edgeId)
	if err != nil {
		return
	}
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	edgeBB, err := s.getEdgeBackbone(edgeId)
	if err != nil {
		return
	}

	select {
	case <-timer.C:
		_ = es.Close()
		return
	case <-edgeBB.isReady:
	}

	bb, err := s.br.GetById(edgeBB.bbId)
	if err != nil {
		return
	}
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
	// todo handle rerouting if backbone is nil
}

func (s *server) appendTunnelAndBackbone(tunnelId, backboneId common.Id) {
	if tunnelId == nil || backboneId == nil {
		return
	}

	s.tunnelToBackboneMu.Lock()
	defer s.tunnelToBackboneMu.Unlock()
	for _, bId := range s.tunnelToBackbone[tunnelId.String()] {
		if bId == backboneId {
			return
		}
	}
	s.tunnelToBackbone[tunnelId.String()] = append(s.tunnelToBackbone[tunnelId.String()], backboneId)
	s.bbMu.Lock()
	defer s.bbMu.Unlock()
	s.bbTunnels[backboneId.String()] = append(s.bbTunnels[backboneId.String()], tunnelId)
}

func (s *server) registerEdgeToBackbone(edgeId, backboneId common.Id) {
	s.edgeToBackboneMu.Lock()
	defer s.edgeToBackboneMu.Unlock()
	s.edgeToBackbone[edgeId.String()] = &edgeBackbone{
		bbId:    backboneId,
		isReady: make(chan struct{}),
	}
	s.bbMu.Lock()
	defer s.bbMu.Unlock()
	s.bbEdges[backboneId.String()] = append(s.bbEdges[backboneId.String()], edgeId)
}

func (s *server) getEdgeBackbone(edgeId common.Id) (*edgeBackbone, error) {
	s.edgeToBackboneMu.RLock()
	defer s.edgeToBackboneMu.RUnlock()
	edgeBB, found := s.edgeToBackbone[edgeId.String()]
	if !found || edgeBB == nil {
		return nil, ErrEdgeNotRegistered
	}
	return edgeBB, nil
}

func (s *server) getBackboneByTunnelInfoId(tunnelInfoId common.Id) (backbone.Backbone, error) {
	s.tunnelToBackboneMu.RLock()
	defer s.tunnelToBackboneMu.RUnlock()
	bbIds, found := s.tunnelToBackbone[tunnelInfoId.String()]
	if !found || len(bbIds) <= 0 {
		return nil, ErrTunnelInfoNotRegistered
	}
	s.tunnelToBackboneLB %= len(bbIds)

	var (
		prev = s.tunnelToBackboneLB
		bb   backbone.Backbone
		err  error
	)
	for {
		s.tunnelToBackboneLB++
		s.tunnelToBackboneLB %= len(bbIds)
		log.Println(s.tunnelToBackboneLB, prev)
		if s.tunnelToBackboneLB == prev {
			break
		}
		bb, err = s.br.GetById(bbIds[s.tunnelToBackboneLB])
		if err != nil {
			continue
		}
		if bb.IsClosed() {
			continue
		}
		break
	}
	if bb == nil {
		bb, err = s.br.GetById(bbIds[s.tunnelToBackboneLB])
		if err != nil {
			return nil, ErrBackboneNotFound
		}
		if bb.IsClosed() {
			return nil, ErrBackboneNotFound
		}
	}
	return bb, nil
}

func (s *server) unregisterBackbone(bbId common.Id) {
	defer log.Println("done unregistering backbone")
	s.bbMu.Lock()
	defer s.bbMu.Unlock()
	if edgeIds, found := s.bbEdges[bbId.String()]; found {
		for _, edgeId := range edgeIds {
			s.unregisterEdge(edgeId)
			if es, err := s.er.GetById(edgeId); err == nil {
				_ = es.Close()
			}
		}
	}
	delete(s.bbEdges, bbId.String())
	if tunnelIds, found := s.bbTunnels[bbId.String()]; found {
		for _, tunnelId := range tunnelIds {
			s.unregisterTunnelInfoAndBackbone(tunnelId, bbId)
		}
	}
	delete(s.bbTunnels, bbId.String())
}

func (s *server) unregisterTunnelInfoAndBackbone(tunnelInfoId common.Id, bbId common.Id) {
	s.tunnelToBackboneMu.Lock()
	defer s.tunnelToBackboneMu.Unlock()
	data, ok := s.tunnelToBackbone[tunnelInfoId.String()]
	log.Println("unregisterTunnelInfoAndBackbone", data)
	if !ok {
		return
	}
	log.Println("unregisterTunnelInfoAndBackbone", data)
	found := 0
	for i := range data {
		for {
			if i+found >= len(data) {
				break
			}
			if data[i+found].Equals(bbId) {
				found++
				continue
			}
			break
		}
		if i+found >= len(data) {
			break
		}
		data[i] = data[i+found]
	}
	s.tunnelToBackbone[tunnelInfoId.String()] = data[:len(data)-found]
}

func (s *server) unregisterEdge(edgeId common.Id) {
	s.edgeToBackboneMu.Lock()
	defer s.edgeToBackboneMu.Unlock()
	delete(s.edgeToBackbone, edgeId.String())
}
