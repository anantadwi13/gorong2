package backbone

import (
	"errors"
	"time"
)

var (
	ErrUnknownMessageType = errors.New("unknown message type")
	ErrMarshalUnmarshal   = errors.New("error marshal/unmarshal")
)

type MessageType byte

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeHandshake
	MessageTypeHandshakeRes
	MessageTypePing
	MessageTypePong
	MessageTypeRegisterEdge
	MessageTypeRegisterEdgeRes
	MessageTypeUnregisterEdge
	MessageTypeUnregisterEdgeRes
)

type MessageFactory interface {
	NewMessage(msgType MessageType) (Message, error)
	MarshalMessage(msg Message) ([]byte, error)
	UnmarshalMessage(buf []byte, msg Message) error
}

type Message interface {
	MessageType() MessageType
}

type HandshakeMessage struct {
	Version string
	Token   []byte
}

func (h HandshakeMessage) MessageType() MessageType {
	return MessageTypeHandshake
}

type HandshakeResMessage struct {
	Error   error
	Version string
}

func (h HandshakeResMessage) MessageType() MessageType {
	return MessageTypeHandshakeRes
}

type PingMessage struct {
	Time time.Time
}

func (p PingMessage) MessageType() MessageType {
	return MessageTypePing
}

type PongMessage struct {
	Error error
	Time  time.Time
}

func (p PongMessage) MessageType() MessageType {
	return MessageTypePong
}

type RegisterEdgeMessage struct {
	EdgeName        string // must be unique
	EdgeType        string // "tcp" or "udp"
	LoadBalancerKey string

	EdgeServerAddr string
}

func (r RegisterEdgeMessage) MessageType() MessageType {
	return MessageTypeRegisterEdge
}

type RegisterEdgeResMessage struct {
	Error        error
	EdgeRunnerId string
	EdgeName     string
}

func (r RegisterEdgeResMessage) MessageType() MessageType {
	return MessageTypeRegisterEdgeRes
}

type UnregisterEdgeMessage struct {
	EdgeRunnerId string
}

func (u UnregisterEdgeMessage) MessageType() MessageType {
	return MessageTypeUnregisterEdge
}

type UnregisterEdgeResMessage struct {
	Error        error
	EdgeRunnerId string
	EdgeName     string
}

func (u UnregisterEdgeResMessage) MessageType() MessageType {
	return MessageTypeUnregisterEdgeRes
}
