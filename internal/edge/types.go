package edge

import (
	"context"
	"errors"

	"github.com/anantadwi13/gorong2/internal/common"
)

type Edge interface {
	Id() common.Id
	TunnelInfo() common.TunnelInfo
	SendData(data []byte) error
	IsClosed() bool
	Close() error
	Serve() error
}

type OnTunnelEstablished func(tunnelInfo common.TunnelInfo)

type EdgeHandler interface {
	OnEdgeCreated(edge Edge)
	OnEdgeClosed(edge Edge)
}

type EdgeDataHandler interface {
	OnEdgeDataReceived(edgeId common.Id, data []byte)
}

type Registry interface {
	RegisterEdgeHandler(edgeHandlers ...EdgeHandler) error
	RegisterEdgeDataHandler(edgeDataHandlers ...EdgeDataHandler) error

	IsRunning() bool
	Listen(tunnelInfo common.TunnelInfo, onTunnelEstablished OnTunnelEstablished) error
	Dial(edgeId common.Id, tunnelInfo common.TunnelInfo) error
	Shutdown(ctx context.Context) error

	GetAll() ([]Edge, error)
	GetById(edgeId common.Id) (Edge, error)
	SendData(edgeId common.Id, data []byte) error
}

var (
	ErrAddressInvalid    = errors.New("address is invalid")
	ErrTunnelInfoInvalid = errors.New("tunnel info is invalid")
	ErrUnableDial        = errors.New("unable to dial")
	ErrUnableListen      = errors.New("unable to listen")
	ErrIllegalState      = errors.New("illegal state")
	ErrEdgeNotFound      = errors.New("edge is not found")
	ErrEdgeInvalid       = errors.New("edge is invalid")
	ErrEdgeClosed        = errors.New("edge is closed")
)
