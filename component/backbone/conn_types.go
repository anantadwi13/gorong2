package backbone

import (
	"net"

	"github.com/anantadwi13/gorong2/pkg/graceful"
)

var (
	ErrClosed = net.ErrClosed
)

type ConnectionFactory interface {
	Listen(addr string) (ServerListener, error)
	DialController(addr string) (Conn, error)
	DialWorker(addr string) (Conn, error)
}

type ServerListener interface {
	graceful.Service
	net.Listener
	AcceptController() (Conn, error)
	AcceptWorker() (Conn, error)
}

type Conn interface {
	net.Conn
}
