package backbone

import (
	"context"
	"errors"

	"github.com/anantadwi13/gorong2/internal/common"
)

type Protocol int

const (
	ProtocolWebsocket Protocol = iota
)

type BackboneAddress interface {
	Protocol() Protocol
	String() string
}

type PacketType int

const (
	PacketTypeUnknown PacketType = iota
	PacketTypeOpen
	PacketTypeOpenAckSuccess
	PacketTypeOpenAckFailed
	PacketTypeEdgeCreate
	PacketTypeEdgeCreateAckSuccess
	PacketTypeEdgeCreateAckFailed
	PacketTypeEdgeClose
	PacketTypeMessage
)

type Packet interface {
	Id() common.Id
	Type() PacketType
	Data() []byte
}

type PacketHandler interface {
	OnReceive(backboneId common.Id, packet Packet)
}

type Backbone interface {
	Id() common.Id
	Serve() error
	Send(packet Packet) error
	IsClosed() bool
	Close() error
}

type BackboneHandler interface {
	OnCreate(b Backbone)
	OnClosed(b Backbone)
}

type Registry interface {
	SetAddress(address BackboneAddress) error
	SetPacketHandler(handler PacketHandler) error
	SetBackboneHandler(handler BackboneHandler) error

	Start() error
	IsRunning() bool
	Listen() error
	Dial() error
	Shutdown(ctx context.Context) error

	GetAll() ([]Backbone, error)
	GetById(backboneId common.Id) (Backbone, error)
	Send(backboneId common.Id, packet Packet) error
}

var (
	ErrAddressNil       = errors.New("address is nil")
	ErrAddressInvalid   = errors.New("invalid address")
	ErrCertRootInvalid  = errors.New("root ca is invalid")
	ErrCertPrivInvalid  = errors.New("private key is invalid")
	ErrCertPubInvalid   = errors.New("public key is invalid")
	ErrUnableDial       = errors.New("unable to dial")
	ErrUnableListen     = errors.New("unable to listen")
	ErrIllegalState     = errors.New("illegal state")
	ErrBackboneNotFound = errors.New("backbone is not found")
	ErrBackboneInvalid  = errors.New("backbone is invalid")
	ErrBackboneClosed   = errors.New("backbone is closed")
	ErrPacketInvalid    = errors.New("packet is invalid")
)
