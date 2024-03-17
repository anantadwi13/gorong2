package server

import (
	"context"
	"errors"
)

type Server interface {
	Start() error
	Shutdown(context.Context) error
}

var (
	ErrConfigNotValid          = errors.New("config is not valid")
	ErrBackboneNotFound        = errors.New("backbone is not found")
	ErrTunnelInfoNotRegistered = errors.New("tunnel info is not registered")
	ErrEdgeNotRegistered       = errors.New("edge is not registered")
)
