package common

import "errors"

type TunnelType string

const (
	TunnelTypeHttp TunnelType = "http"
	TunnelTypeTcp  TunnelType = "tcp"
	TunnelTypeUdp  TunnelType = "udp"
)

type TunnelInfo interface {
	Id() Id
	Type() TunnelType
	EdgeServerAddr() string
	EdgeAgentAddr() string
}

type Id interface {
	Raw() interface{}
	Equals(other interface{}) bool
	String() string
}

var (
	ErrTunnelInfoInvalid = errors.New("tunnel info is invalid")
)
