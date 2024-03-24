package backbone

import (
	"errors"

	"github.com/anantadwi13/gorong2/component/backbone"
	"github.com/anantadwi13/gorong2/pkg/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProtobufMessageFactory struct {
}

func (p *ProtobufMessageFactory) NewMessage(msgType backbone.MessageType) (msg backbone.Message, err error) {
	switch msgType {
	case backbone.MessageTypeHandshake:
		msg = &backbone.HandshakeMessage{}
	case backbone.MessageTypeHandshakeRes:
		msg = &backbone.HandshakeResMessage{}
	case backbone.MessageTypePing:
		msg = &backbone.PingMessage{}
	case backbone.MessageTypePong:
		msg = &backbone.PongMessage{}
	case backbone.MessageTypeRegisterEdge:
		msg = &backbone.RegisterEdgeMessage{}
	case backbone.MessageTypeRegisterEdgeRes:
		msg = &backbone.RegisterEdgeResMessage{}
	case backbone.MessageTypeUnregisterEdge:
		msg = &backbone.UnregisterEdgeMessage{}
	case backbone.MessageTypeUnregisterEdgeRes:
		msg = &backbone.UnregisterEdgeResMessage{}
	default:
		err = backbone.ErrUnknownMessageType
	}
	return
}

func (p *ProtobufMessageFactory) MarshalMessage(msg backbone.Message) ([]byte, error) {
	var protoMsg proto.Message
	switch m := msg.(type) {
	case *backbone.HandshakeMessage:
		protoMsg = &Handshake{
			Version: m.Version,
			Token:   m.Token,
		}
	case *backbone.HandshakeResMessage:
		protoMsg = &HandshakeRes{
			Error:   utils.ErrorToString(m.Error),
			Version: m.Version,
		}
	case *backbone.PingMessage:
		protoMsg = &Ping{
			Time: timestamppb.New(m.Time),
		}
	case *backbone.PongMessage:
		protoMsg = &Pong{
			Error: utils.ErrorToString(m.Error),
			Time:  timestamppb.New(m.Time),
		}
	case *backbone.RegisterEdgeMessage:
		protoMsg = &RegisterEdge{
			EdgeName:        m.EdgeName,
			EdgeType:        m.EdgeType,
			LoadBalancerKey: m.LoadBalancerKey,
			EdgeServerAddr:  m.EdgeServerAddr,
		}
	case *backbone.RegisterEdgeResMessage:
		protoMsg = &RegisterEdgeRes{
			Error:        utils.ErrorToString(m.Error),
			EdgeRunnerId: m.EdgeRunnerId,
			EdgeName:     m.EdgeName,
		}
	case *backbone.UnregisterEdgeMessage:
		protoMsg = &UnregisterEdge{
			EdgeRunnerId: m.EdgeRunnerId,
		}
	case *backbone.UnregisterEdgeResMessage:
		protoMsg = &UnregisterEdgeRes{
			Error:        utils.ErrorToString(m.Error),
			EdgeRunnerId: m.EdgeRunnerId,
			EdgeName:     m.EdgeName,
		}
	}
	if protoMsg == nil {
		return nil, backbone.ErrUnknownMessageType
	}
	res, err := proto.Marshal(protoMsg)
	if err != nil {
		return nil, errors.Join(backbone.ErrMarshalUnmarshal, err)
	}
	return res, nil
}

func (p *ProtobufMessageFactory) UnmarshalMessage(buf []byte, msg backbone.Message) error {
	if buf == nil || msg == nil {
		return errors.Join(backbone.ErrMarshalUnmarshal, errors.New("nil parameter"))
	}

	switch m := msg.(type) {
	case *backbone.HandshakeMessage:
		protoMsg := &Handshake{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.Token = protoMsg.Token
		m.Version = protoMsg.Version
	case *backbone.HandshakeResMessage:
		protoMsg := &HandshakeRes{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.Error = utils.StringToError(protoMsg.Error)
		m.Version = protoMsg.Version
	case *backbone.PingMessage:
		protoMsg := &Ping{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.Time = protoMsg.Time.AsTime()
	case *backbone.PongMessage:
		protoMsg := &Pong{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.Error = utils.StringToError(protoMsg.Error)
		m.Time = protoMsg.Time.AsTime()
	case *backbone.RegisterEdgeMessage:
		protoMsg := &RegisterEdge{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.EdgeName = protoMsg.EdgeName
		m.EdgeType = protoMsg.EdgeType
		m.LoadBalancerKey = protoMsg.LoadBalancerKey
		m.EdgeServerAddr = protoMsg.EdgeServerAddr
	case *backbone.RegisterEdgeResMessage:
		protoMsg := &RegisterEdgeRes{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.Error = utils.StringToError(protoMsg.Error)
		m.EdgeRunnerId = protoMsg.EdgeRunnerId
		m.EdgeName = protoMsg.EdgeName
	case *backbone.UnregisterEdgeMessage:
		protoMsg := &UnregisterEdge{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.EdgeRunnerId = protoMsg.EdgeRunnerId
	case *backbone.UnregisterEdgeResMessage:
		protoMsg := &UnregisterEdgeRes{}
		err := proto.Unmarshal(buf, protoMsg)
		if err != nil {
			return err
		}
		m.Error = utils.StringToError(protoMsg.Error)
		m.EdgeRunnerId = protoMsg.EdgeRunnerId
		m.EdgeName = protoMsg.EdgeName
	default:
		return backbone.ErrUnknownMessageType
	}
	return nil
}
