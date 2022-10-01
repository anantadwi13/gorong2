package backbone

import "github.com/anantadwi13/gorong2/internal/common"

type defaultBackboneHandler struct {
	onCreate func(b Backbone)
	onClosed func(b Backbone)
}

func NewBackboneHandler(onCreate func(b Backbone), onClosed func(b Backbone)) BackboneHandler {
	return &defaultBackboneHandler{
		onCreate: onCreate,
		onClosed: onClosed,
	}
}

func (bl *defaultBackboneHandler) OnCreate(b Backbone) {
	if bl.onCreate == nil {
		return
	}
	bl.onCreate(b)
}

func (bl *defaultBackboneHandler) OnClosed(b Backbone) {
	if bl.onClosed == nil {
		return
	}
	bl.onClosed(b)
}

type defaultPacketHandler struct {
	onReceive func(backboneId common.Id, packet Packet)
}

func NewPacketHandler(onReceive func(backboneId common.Id, packet Packet)) PacketHandler {
	return &defaultPacketHandler{onReceive: onReceive}
}

func (pl *defaultPacketHandler) OnReceive(backboneId common.Id, packet Packet) {
	if pl.onReceive == nil {
		return
	}

	pl.onReceive(backboneId, packet)
}
