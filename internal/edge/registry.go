package edge

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/anantadwi13/gorong2/internal/common"
)

type registry struct {
	ctx                context.Context
	ctxCancel          context.CancelFunc
	isRunning          atomic.Int32
	edgeHandlers       []EdgeHandler
	edgeHandlersMu     sync.RWMutex
	edgeDataHandlers   []EdgeDataHandler
	edgeDataHandlersMu sync.RWMutex

	edges sync.Map
}

func NewRegistry() (Registry, error) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return &registry{
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}, nil
}

func (r *registry) RegisterEdgeHandler(edgeHandlers ...EdgeHandler) error {
	if r.IsRunning() {
		return ErrIllegalState
	}

	r.edgeHandlersMu.Lock()
	defer r.edgeHandlersMu.Unlock()

	for _, handler := range edgeHandlers {
		if handler == nil {
			continue
		}
		r.edgeHandlers = append(r.edgeHandlers, handler)
	}

	return nil
}

func (r *registry) RegisterEdgeDataHandler(edgeDataHandlers ...EdgeDataHandler) error {
	if r.IsRunning() {
		return ErrIllegalState
	}

	r.edgeDataHandlersMu.Lock()
	defer r.edgeDataHandlersMu.Unlock()

	for _, handler := range edgeDataHandlers {
		if handler == nil {
			continue
		}
		r.edgeDataHandlers = append(r.edgeDataHandlers, handler)
	}

	return nil
}

func (r *registry) IsRunning() bool {
	return r.isRunning.Load() > 0
}

func (r *registry) Listen(tunnelInfo common.TunnelInfo, ote OnTunnelEstablished) error {
	r.isRunning.Add(1)
	defer r.isRunning.Add(-1)

	if tunnelInfo == nil {
		return ErrTunnelInfoInvalid
	}

	switch tunnelInfo.Type() {
	case common.TunnelTypeTcp:
		return tcpListen(r.ctx, tunnelInfo, ote, r, r)
	case common.TunnelTypeUdp:
		return udpListen(r.ctx, tunnelInfo, ote, r, r)
	}

	return ErrUnableListen
}

func (r *registry) Dial(edgeId common.Id, tunnelInfo common.TunnelInfo) error {
	r.isRunning.Add(1)
	defer r.isRunning.Add(-1)

	if tunnelInfo == nil {
		return ErrTunnelInfoInvalid
	}

	switch tunnelInfo.Type() {
	case common.TunnelTypeTcp:
		return tcpDial(edgeId, tunnelInfo, r, r)
	case common.TunnelTypeUdp:
		return udpDial(edgeId, tunnelInfo, r, r)
	}

	return ErrUnableListen
}

func (r *registry) Shutdown(ctx context.Context) error {
	defer r.ctxCancel()

	edges, err := r.GetAll()
	if err != nil {
		return err
	}

	for _, edge := range edges {
		if ctx != nil && ctx.Err() != nil {
			return err
		}
		// todo handle error
		_ = edge.Close()
	}

	// todo stop listener
	return nil
}

func (r *registry) GetAll() ([]Edge, error) {
	var edges []Edge

	r.edges.Range(func(key, value any) bool {
		e, ok := value.(Edge)
		if ok && e != nil {
			edges = append(edges, e)
		}
		return true
	})

	return edges, nil
}

func (r *registry) GetById(edgeId common.Id) (Edge, error) {
	eRaw, ok := r.edges.Load(edgeId.String())
	if !ok || eRaw == nil {
		return nil, ErrEdgeNotFound
	}
	e, ok := eRaw.(Edge)
	if !ok || e == nil {
		return nil, ErrEdgeInvalid
	}
	return e, nil
}

func (r *registry) SendData(edgeId common.Id, data []byte) error {
	e, err := r.GetById(edgeId)
	if err != nil {
		return err
	}
	err = e.SendData(data)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) OnEdgeCreated(edge Edge) {
	if edge != nil {
		r.edges.Store(edge.Id().String(), edge)
	}

	r.edgeHandlersMu.RLock()
	for _, handler := range r.edgeHandlers {
		if handler == nil {
			continue
		}
		handler.OnEdgeCreated(edge)
	}
	r.edgeHandlersMu.RUnlock()
}

func (r *registry) OnEdgeClosed(edge Edge) {
	r.edgeHandlersMu.RLock()
	for _, handler := range r.edgeHandlers {
		if handler == nil {
			continue
		}
		handler.OnEdgeClosed(edge)
	}
	r.edgeHandlersMu.RUnlock()

	if edge != nil {
		r.edges.Delete(edge.Id().String())
	}
}

func (r *registry) OnEdgeDataReceived(edgeId common.Id, data []byte) {
	r.edgeDataHandlersMu.RLock()
	for _, handler := range r.edgeDataHandlers {
		if handler == nil {
			continue
		}
		handler.OnEdgeDataReceived(edgeId, data)
	}
	r.edgeDataHandlersMu.RUnlock()
}
