package backbone

import (
	"errors"
	"testing"
	"time"

	"github.com/anantadwi13/gorong2/component/backbone"
	"github.com/stretchr/testify/assert"
)

func TestProtobufMessageFactory_MarshalUnmarshalMessage(t *testing.T) {
	p := &ProtobufMessageFactory{}
	type args struct {
		msg     backbone.Message
		msgResp backbone.Message
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handshake",
			args: args{
				msg: &backbone.HandshakeMessage{
					Version: "1.2.3",
					Token:   []byte("zkxjchqwiuhdas"),
				},
				msgResp: &backbone.HandshakeMessage{},
			},
		},
		{
			name: "handshake res",
			args: args{
				msg: &backbone.HandshakeResMessage{
					Error:   errors.New("error test"),
					Version: "1.2.3",
				},
				msgResp: &backbone.HandshakeResMessage{},
			},
		},
		{
			name: "ping",
			args: args{
				msg: &backbone.PingMessage{
					Time: time.Now().UTC(),
				},
				msgResp: &backbone.PingMessage{},
			},
		},
		{
			name: "pong",
			args: args{
				msg: &backbone.PongMessage{
					Error: errors.New("error test"),
					Time:  time.Now().UTC(),
				},
				msgResp: &backbone.PongMessage{},
			},
		},
		{
			name: "register edge",
			args: args{
				msg: &backbone.RegisterEdgeMessage{
					EdgeName:        "app1",
					EdgeType:        "tcp",
					LoadBalancerKey: "askdljasd",
					EdgeServerAddr:  ":9090",
				},
				msgResp: &backbone.RegisterEdgeMessage{},
			},
		},
		{
			name: "register edge res",
			args: args{
				msg: &backbone.RegisterEdgeResMessage{
					Error:        errors.New("error test"),
					EdgeRunnerId: "qwerty",
					EdgeName:     "app1",
				},
				msgResp: &backbone.RegisterEdgeResMessage{},
			},
		},
		{
			name: "unregister edge",
			args: args{
				msg: &backbone.UnregisterEdgeMessage{
					EdgeRunnerId: "qwerty",
				},
				msgResp: &backbone.UnregisterEdgeMessage{},
			},
		},
		{
			name: "unregister edge res",
			args: args{
				msg: &backbone.UnregisterEdgeResMessage{
					Error:        errors.New("error test"),
					EdgeRunnerId: "qwerty",
					EdgeName:     "app1",
				},
				msgResp: &backbone.UnregisterEdgeResMessage{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.MarshalMessage(tt.args.msg)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			err = p.UnmarshalMessage(got, tt.args.msgResp)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.msg, tt.args.msgResp)
		})
	}
}
