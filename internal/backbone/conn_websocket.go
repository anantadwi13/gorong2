package backbone

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anantadwi13/gorong2/component/backbone"
	"github.com/anantadwi13/gorong2/pkg/log"
	"github.com/anantadwi13/gorong2/pkg/pool"
	"github.com/gorilla/websocket"
)

const (
	connTypeController = byte(0x01)
	connTypeWorker     = byte(0x02)
)

type WebsocketConnectionFactory struct {
	messageFactory backbone.MessageFactory
}

func NewWebsocketConnectionFactory(msgFactory backbone.MessageFactory) (*WebsocketConnectionFactory, error) {
	return &WebsocketConnectionFactory{
		messageFactory: msgFactory,
	}, nil
}

func (w *WebsocketConnectionFactory) Listen(addr string) (backbone.ServerListener, error) {
	listener := &websocketServerListener{
		addr:           addr,
		messageFactory: w.messageFactory,
		controllerConn: make(chan backbone.ControllerConn, 1024),
		workerConn:     make(chan backbone.WorkerConn, 1024),
	}
	err := listener.looper()
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (w *WebsocketConnectionFactory) DialController(addr string) (backbone.ControllerConn, error) {
	return w.dial(addr, connTypeController)
}

func (w *WebsocketConnectionFactory) DialWorker(addr string) (backbone.WorkerConn, error) {
	return w.dial(addr, connTypeWorker)
}

func (w *WebsocketConnectionFactory) dial(addr string, connType byte) (*websocketConn, error) {
	var err error

	wsConn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	conn := &websocketConn{
		connType:       connType,
		ctx:            ctx,
		messageFactory: w.messageFactory,
		wsConn:         wsConn,
	}
	defer func() {
		if err != nil {
			log.Trace(ctx, "close websocket conn", wsConn.RemoteAddr())
			_ = conn.Close()
		}
	}()

	n, err := conn.Write([]byte{connType})
	if err != nil {
		return nil, err
	}
	if n != 1 {
		err = errors.New("invalid write length")
		return nil, err
	}

	return conn, nil
}

type websocketServerListener struct {
	wg         sync.WaitGroup
	isShutdown atomic.Uint32
	addr       string
	server     *http.Server
	ctx        context.Context
	cancel     context.CancelFunc

	messageFactory backbone.MessageFactory

	controllerConn chan backbone.ControllerConn
	workerConn     chan backbone.WorkerConn
	listener       net.Listener
}

func (w *websocketServerListener) looper() error {
	wsUpgrader := &websocket.Upgrader{}
	handler := http.NewServeMux()
	w.server = &http.Server{
		Addr:    w.addr,
		Handler: handler,
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())

	handler.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		var (
			ctx = context.Background()
			err error
		)
		if !websocket.IsWebSocketUpgrade(r) {
			rw.WriteHeader(http.StatusOK)
			return
		}
		wsConn, err := wsUpgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Debugf(ctx, "couldn't upgrade websocket. err: %s", err)
			return
		}

		conn := &websocketConn{
			ctx:            ctx,
			messageFactory: w.messageFactory,
			wsConn:         wsConn,
		}
		defer func() {
			if err != nil {
				log.Trace(ctx, "close websocket conn", wsConn.RemoteAddr())
				_ = conn.Close()
			}
		}()

		buf := make([]byte, 1)
		// read connection type
		n, err := io.ReadFull(conn, buf)
		if err != nil {
			return
		}
		buf = buf[:n]
		if len(buf) != 1 {
			err = errors.New("invalid message length")
			return
		}
		log.Trace(ctx, "ws: success read connection type")
		switch buf[0] {
		case connTypeController:
			conn.connType = connTypeController
			w.controllerConn <- conn
		case connTypeWorker:
			conn.connType = connTypeWorker
			w.workerConn <- conn
		default:
			err = errors.New("ws: unknown connection type")
			return
		}
	})

	ln, err := net.Listen("tcp", w.addr)
	if err != nil {
		log.Error(w.ctx, "error listening address", err)
		return err
	}
	w.listener = ln

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer w.cancel()
		err := w.server.Serve(ln)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error(w.ctx, "error server closed", err)
			}
			return
		}
	}()
	return nil
}

func (w *websocketServerListener) Shutdown(ctx context.Context) error {
	if !w.isShutdown.CompareAndSwap(0, 1) {
		return nil
	}
	defer log.Trace(ctx, "wsListener.Shutdown complete")
	defer func() {
		close(w.controllerConn)
		close(w.workerConn)
		log.Trace(ctx, "wsListener.Shutdown channel closed")
	}()
	defer w.wg.Wait()

	log.Trace(ctx, "wsListener.Shutdown")

	err := w.server.Shutdown(ctx)
	if err != nil {
		log.Trace(ctx, "wsListener.Shutdown error", err)
		return err
	}
	return nil
}

func (w *websocketServerListener) Accept() (net.Conn, error) {
	if w.listener == nil {
		return nil, errors.New("listener is not initialized yet")
	}
	return w.listener.Accept()
}

func (w *websocketServerListener) Close() error {
	return w.Shutdown(context.Background())
}

func (w *websocketServerListener) Addr() net.Addr {
	if w.listener == nil {
		return nil
	}
	return w.listener.Addr()
}

func (w *websocketServerListener) AcceptController() (backbone.ControllerConn, error) {
	conn, ok := <-w.controllerConn
	if !ok {
		return nil, net.ErrClosed
	}
	return conn, nil
}

func (w *websocketServerListener) AcceptWorker() (backbone.WorkerConn, error) {
	conn, ok := <-w.workerConn
	if !ok {
		return nil, net.ErrClosed
	}
	return conn, nil
}

type websocketConn struct {
	connType       byte // 0x01 controller, 0x02 worker
	ctx            context.Context
	messageFactory backbone.MessageFactory
	wsConn         *websocket.Conn

	lastReader io.Reader
	writerBuf  *bufio.Writer
}

func (w *websocketConn) Read(b []byte) (read int, err error) {
	var (
		n      = 0
		reader io.Reader
	)
	log.Trace(w.ctx, "w.Read", len(b))
	defer func() { log.Trace(w.ctx, "w.Read out", err, len(b)) }()
	for read == 0 {
		if w.lastReader == nil {
			_, reader, err = w.wsConn.NextReader()
			if err != nil {
				log.Trace(w.ctx, "error w.wsConn.NextReader", w.connType, err, read, n, reader)
				if websocket.IsUnexpectedCloseError(err) {
					err = errors.Join(net.ErrClosed, err)
				}
				return
			}
			w.lastReader = reader
		}

		if w.lastReader == nil {
			err = errors.New("cannot initialize reader")
			return
		}

		n, err = w.lastReader.Read(b[read:])
		read += n
		log.Trace(w.ctx, "w.lastReader.Read", n, err, read)
		if err != nil {
			log.Trace(w.ctx, "got error", err, read, w.lastReader)
			switch {
			case errors.Is(err, io.EOF) && read == 0:
				// lastReader has reached its end, get next reader
				w.lastReader = nil
			case errors.Is(err, io.EOF):
				w.lastReader = nil
				return read, nil
			default:
				return
			}
		}
	}
	return
}

func (w *websocketConn) Write(b []byte) (wrote int, err error) {
	n := 0
	writer, err := w.wsConn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return
	}
	defer writer.Close()
	log.Trace(w.ctx, "write", w.connType, len(b))
	n, err = writer.Write(b)
	wrote += n
	if err != nil {
		return
	}
	return
}

func (w *websocketConn) ReadMessage() (backbone.Message, error) {
	bufMsgType := make([]byte, 1)
	_, err := io.ReadFull(w, bufMsgType)
	if err != nil {
		return nil, err
	}

	var messageLen int64
	err = binary.Read(w, binary.BigEndian, &messageLen)
	if err != nil {
		return nil, err
	}

	buf := pool.GetBuffer(int(messageLen))
	defer pool.PutBuffer(buf)

	buf.Buf = buf.Buf[:messageLen]
	_, err = io.ReadFull(w, buf.Buf)
	if err != nil {
		return nil, err
	}

	msg, err := w.messageFactory.NewMessage(backbone.MessageType(bufMsgType[0]))
	if err != nil {
		return nil, err
	}
	err = w.messageFactory.UnmarshalMessage(buf.Buf, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (w *websocketConn) WriteMessage(msg backbone.Message) error {
	if w.writerBuf == nil {
		w.writerBuf = bufio.NewWriterSize(w, 4096) // todo change buffer size
	}
	w.writerBuf.Reset(w)

	rawMsg, err := w.messageFactory.MarshalMessage(msg)
	if err != nil {
		return err
	}
	_, err = w.writerBuf.Write([]byte{byte(msg.MessageType())})
	if err != nil {
		return err
	}

	err = binary.Write(w.writerBuf, binary.BigEndian, int64(len(rawMsg)))
	if err != nil {
		return err
	}

	_, err = w.writerBuf.Write(rawMsg)
	if err != nil {
		return err
	}

	return w.writerBuf.Flush()
}

func (w *websocketConn) Close() error {
	log.Trace(w.ctx, "close websocket conn", w.LocalAddr(), w.RemoteAddr(), w.connType)
	var err error
	err = w.wsConn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(30*time.Second),
	)
	if err != nil {
		return err
	}
	err = w.wsConn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (w *websocketConn) LocalAddr() net.Addr {
	return w.wsConn.LocalAddr()
}

func (w *websocketConn) RemoteAddr() net.Addr {
	return w.wsConn.RemoteAddr()
}

func (w *websocketConn) SetDeadline(t time.Time) error {
	return w.wsConn.UnderlyingConn().SetDeadline(t)
}

func (w *websocketConn) SetReadDeadline(t time.Time) error {
	return w.wsConn.SetReadDeadline(t)
}

func (w *websocketConn) SetWriteDeadline(t time.Time) error {
	return w.wsConn.SetWriteDeadline(t)
}
