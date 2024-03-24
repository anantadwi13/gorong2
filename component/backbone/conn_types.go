package backbone

import (
	"net"

	"github.com/anantadwi13/gorong2/pkg/graceful"
)

var (
	ErrConnClosed = net.ErrClosed
)

type ConnectionFactory interface {
	Listen(addr string) (ServerListener, error)
	DialController(addr string) (ControllerConn, error)
	DialWorker(addr string) (WorkerConn, error)
}

type ServerListener interface {
	graceful.Service
	net.Listener
	AcceptController() (ControllerConn, error)
	AcceptWorker() (WorkerConn, error)
}

type ControllerConn interface {
	net.Conn
	ReadMessage() (Message, error)
	WriteMessage(msg Message) error
}

type WorkerConn interface {
	net.Conn
}
