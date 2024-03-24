package server

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anantadwi13/gorong2/component/backbone"
	backboneImpl "github.com/anantadwi13/gorong2/internal/backbone"
	"github.com/anantadwi13/gorong2/pkg/log"
)

type Server struct {
	cfg Config

	wg        sync.WaitGroup
	shutdown  chan struct{}
	isRunning atomic.Uint32

	connFactory       backbone.ConnectionFactory
	messageReadWriter backbone.MessageReadWriter
	messageFactory    backbone.MessageFactory

	edgeLock sync.RWMutex
	edge     map[string]string
}

func NewServer(cfg Config) (*Server, error) {
	server := &Server{
		shutdown: make(chan struct{}),
		cfg:      cfg,
	}
	return server, nil
}

func (s *Server) Run(ctx context.Context) (err error) {
	if !s.isRunning.CompareAndSwap(0, 1) {
		log.Debug(ctx, "server already running")
		return nil
	}

	// todo retry listen with backoff

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	// todo
	msgFactory := &backboneImpl.ProtobufMessageFactory{}
	s.messageFactory = msgFactory
	s.messageReadWriter = msgFactory

	switch s.cfg.Backbone.ListenType {
	case ListenTypeWebsocket:
		s.connFactory, err = backboneImpl.NewWebsocketConnectionFactory(s.messageFactory)
	case ListenTypeTcp:
		err = errors.New("unimplemented yet")
	}
	if err != nil {
		return
	}
	if s.connFactory == nil {
		err = errors.New("unable to initialize connection")
		return
	}

	listener, err := s.connFactory.Listen(s.cfg.Backbone.ListenAddr)
	if err != nil {
		return
	}

	s.wg.Add(3)
	go func() {
		defer s.wg.Done()
		select {
		case <-ctx.Done():
		case <-s.shutdown:
		}
		cancel()
		_ = listener.Close()
	}()
	go s.acceptController(ctx, listener)
	go s.acceptWorker(ctx, listener)
	return
}

func (s *Server) acceptController(ctx context.Context, listener backbone.ServerListener) {
	defer s.wg.Done()

	localWg := &sync.WaitGroup{}
	defer localWg.Wait()

	for {
		ctrlConn, err := listener.AcceptController()
		if err != nil {
			if !errors.Is(err, backbone.ErrClosed) {
				log.Errorf(ctx, "error accepting controller. err: %v", err)
			}
			return
		}
		log.Debug(ctx, "got new controller", ctrlConn.RemoteAddr())

		localWg.Add(1)
		go s.handleController(ctx, localWg, ctrlConn)
	}
}

func (s *Server) handleController(ctx context.Context, wg *sync.WaitGroup, conn backbone.Conn) {
	defer wg.Done()
	defer conn.Close()

	err := s.handleHandshake(conn)
	if err != nil {
		log.Errorf(ctx, "error handshaking. %s", err)
		return
	}

	var (
		lastPingLock sync.RWMutex
		lastPing     = time.Now()
		localWg      = &sync.WaitGroup{}
	)
	defer localWg.Wait()

	ctx = initConnID(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	defer func() {
		// unregister all service related to this controller
	}()

	localWg.Add(1)
	go func() {
		// health check handler
		defer cancel()
		defer localWg.Done()
		defer conn.Close() // agent is not responding after several times, then close the connection

		ticker := time.NewTicker(10 * time.Second) // todo make a setting of this value
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Trace(ctx, "health check stopped normally")
				return
			case <-ticker.C:
				lastPingLock.RLock()
				if time.Now().Sub(lastPing) > 30*time.Second {
					log.Error(ctx, "health check error: agent is not responding", lastPing)
					lastPingLock.RUnlock()
					return
				}
				lastPingLock.RUnlock()
			}
		}
	}()

	defer cancel()
	for {
		message, err := s.messageReadWriter.ReadMessage(conn)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Debugf(ctx, "error network closed. %s", err)
			} else {
				log.Errorf(ctx, "error read message. %s", err)
			}
			return
		}
		log.Trace(ctx, "got new message", conn.RemoteAddr(), message.MessageType())

		switch msg := message.(type) {
		case *backbone.PingMessage:
			localWg.Add(1)
			go func() {
				defer localWg.Done()
				lastPingLock.Lock()
				if msg.Time.After(lastPing) {
					lastPing = msg.Time
				}
				lastPingLock.Unlock()
				err := s.messageReadWriter.WriteMessage(conn, &backbone.PongMessage{
					Error: nil,
					Time:  time.Now(),
				})
				if err != nil {
					log.Errorf(ctx, "error sending pong message. %s", err)
					return
				}
			}()
		case *backbone.RegisterEdgeMessage:
			localWg.Add(1)
			go func() {
				defer localWg.Done()
				// todo
			}()
		case *backbone.UnregisterEdgeMessage:
			localWg.Add(1)
			go func() {
				defer localWg.Done()
				// todo
			}()
		default:
			continue
		}
	}
}

func (s *Server) handleHandshake(conn backbone.Conn) error {
	message, err := s.messageReadWriter.ReadMessage(conn)
	if err != nil {
		return err
	}
	handshakeMsg, ok := message.(*backbone.HandshakeMessage)
	if !ok {
		return errors.New("not a handshake message")
	}

	// todo check auth & version
	_ = handshakeMsg

	return s.messageReadWriter.WriteMessage(conn, &backbone.HandshakeResMessage{
		Error:   nil,
		Version: "1.2.3", // todo change
	})
}

func (s *Server) acceptWorker(ctx context.Context, listener backbone.ServerListener) {
	defer s.wg.Done()

	localWg := &sync.WaitGroup{}
	defer localWg.Wait()

	for {
		workerConn, err := listener.AcceptWorker()
		if err != nil {
			if !errors.Is(err, backbone.ErrClosed) {
				log.Errorf(ctx, "error accepting worker. err: %v", err)
			}
			return
		}
		log.Debug(ctx, "got new worker", workerConn.RemoteAddr())

		localWg.Add(1)
		go s.handleWorker(ctx, localWg, workerConn)
	}
}

func (s *Server) handleWorker(ctx context.Context, wg *sync.WaitGroup, conn backbone.Conn) {

}

func (s *Server) Shutdown(ctx context.Context) error {
	close(s.shutdown)
	s.wg.Wait()
	return nil
}
