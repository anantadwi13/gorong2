package agent

import (
	"context"
	"errors"
)

type Agent interface {
	Start() error
	Shutdown(ctx context.Context) error
}

var (
	ErrConfigNotValid     = errors.New("config is not valid")
	ErrTunnelInfoNotFound = errors.New("tunnel info is not found")
	ErrEdgeNotRegistered  = errors.New("edge is not registered")
)
