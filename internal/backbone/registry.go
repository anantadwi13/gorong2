package backbone

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/anantadwi13/gorong2/internal/common"
)

type registry struct {
	isRunning atomic.Bool
	addr      BackboneAddress
	addrMu    sync.RWMutex
	pl        PacketHandler
	plMu      sync.RWMutex
	bl        BackboneHandler
	blMu      sync.RWMutex

	backbones sync.Map
}

func NewRegistry() (Registry, error) {
	return &registry{}, nil
}

func (r *registry) SetAddress(address BackboneAddress) error {
	if r.IsRunning() {
		return ErrIllegalState
	}
	r.addrMu.Lock()
	r.addr = address
	r.addrMu.Unlock()
	return nil
}

func (r *registry) SetPacketHandler(pl PacketHandler) error {
	if r.IsRunning() {
		return ErrIllegalState
	}
	r.plMu.Lock()
	r.pl = pl
	r.plMu.Unlock()
	return nil
}

func (r *registry) SetBackboneHandler(bl BackboneHandler) error {
	if r.IsRunning() {
		return ErrIllegalState
	}
	r.blMu.Lock()
	r.bl = bl
	r.blMu.Unlock()
	return nil
}

func (r *registry) Start() error {
	//TODO implement me
	panic("implement me")
}

func (r *registry) IsRunning() bool {
	return r.isRunning.Load()
}

func (r *registry) Shutdown(ctx context.Context) error {
	backbones, err := r.GetAll()
	if err != nil {
		return err
	}

	for _, b := range backbones {
		if ctx != nil && ctx.Err() != nil {
			return err
		}
		// todo handle error
		_ = b.Close()
	}

	// todo stop listener
	return nil
}

func (r *registry) Listen() error {
	r.isRunning.Store(true)
	defer r.isRunning.Store(false)

	if r.addr == nil {
		return ErrAddressNil
	}

	switch r.addr.(type) {
	case *wsAddress:
		return wsListen(r.addr, r, r)
	}

	return ErrUnableListen
}

func (r *registry) Dial() error {
	r.isRunning.Store(true)
	defer r.isRunning.Store(false)

	if r.addr == nil {
		return ErrAddressNil
	}

	switch r.addr.(type) {
	case *wsAddress:
		return wsDial(r.addr, r, r)
	}

	return ErrUnableDial
}

func (r *registry) GetAll() ([]Backbone, error) {
	var backbones []Backbone

	r.backbones.Range(func(key, value any) bool {
		b, ok := value.(Backbone)
		if ok && b != nil {
			backbones = append(backbones, b)
		}
		return true
	})

	return backbones, nil
}

func (r *registry) GetById(backboneId common.Id) (Backbone, error) {
	bRaw, ok := r.backbones.Load(backboneId.String())
	if !ok || bRaw == nil {
		return nil, ErrBackboneNotFound
	}
	b, ok := bRaw.(Backbone)
	if !ok || bRaw == nil {
		return nil, ErrBackboneInvalid
	}
	return b, nil
}

func (r *registry) Send(backboneId common.Id, packet Packet) error {
	bb, err := r.GetById(backboneId)
	if err != nil {
		return err
	}

	err = bb.Send(packet)
	if err != nil {
		return err
	}

	return nil
}

func (r *registry) OnReceive(backboneId common.Id, packet Packet) {
	r.plMu.RLock()
	if r.pl != nil {
		r.pl.OnReceive(backboneId, packet)
	}
	r.plMu.RUnlock()
}

func (r *registry) OnCreate(b Backbone) {
	if b != nil {
		r.backbones.Store(b.Id().String(), b)
	}

	r.blMu.RLock()
	if r.bl != nil {
		r.bl.OnCreate(b)
	}
	r.blMu.RUnlock()
}

func (r *registry) OnClosed(b Backbone) {
	r.blMu.RLock()
	if r.bl != nil {
		r.bl.OnClosed(b)
	}
	r.blMu.RUnlock()

	if b == nil {
		return
	}
	r.backbones.Delete(b.Id().String())
}
