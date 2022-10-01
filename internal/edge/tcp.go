package edge

import (
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/anantadwi13/gorong2/internal/common"
)

type tcpEdge struct {
	id         common.Id
	tunnelInfo common.TunnelInfo
	tcpConn    *net.TCPConn
	eh         EdgeHandler
	edh        EdgeDataHandler
	writeChan  chan []byte
	ctx        context.Context
	cancelCtx  context.CancelFunc
	isServed   atomic.Bool
}

func tcpDial(edgeId common.Id, tunnelInfo common.TunnelInfo, eh EdgeHandler, edh EdgeDataHandler) error {
	if tunnelInfo == nil {
		return ErrEdgeInvalid
	}

	conn, err := net.Dial("tcp", tunnelInfo.EdgeAgentAddr())
	if err != nil {
		return err
	}
	defer conn.Close()

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return ErrEdgeInvalid
	}

	if edgeId == nil {
		edgeId = common.GenerateId()
	}

	eCtx, eCancelCtx := context.WithCancel(context.Background())
	defer eCancelCtx()

	te := &tcpEdge{
		id:         edgeId,
		tcpConn:    tcpConn,
		tunnelInfo: tunnelInfo,
		eh:         eh,
		edh:        edh,
		ctx:        eCtx,
		cancelCtx:  eCancelCtx,
		writeChan:  make(chan []byte, 1024),
	}

	return te.Serve()
}

func tcpListen(
	ctx context.Context, tunnelInfo common.TunnelInfo, ote OnTunnelEstablished, eh EdgeHandler,
	edh EdgeDataHandler,
) error {
	if ctx == nil {
		return ErrEdgeInvalid
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	if tunnelInfo == nil {
		return ErrEdgeInvalid
	}

	listener, err := net.Listen("tcp", tunnelInfo.EdgeServerAddr())
	if err != nil {
		return err
	}
	defer listener.Close()

	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		return ErrEdgeInvalid
	}

	go func() {
		select {
		case <-ctx.Done():
			_ = listener.Close()
		}
	}()

	if ote != nil {
		ote(tunnelInfo)
	}

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			// todo handle other cases
			return err
		}
		go func(tcpConn *net.TCPConn) {
			defer tcpConn.Close()

			eCtx, eCancelCtx := context.WithCancel(ctx)
			defer eCancelCtx()

			te := &tcpEdge{
				id:         common.GenerateId(),
				tunnelInfo: tunnelInfo,
				tcpConn:    tcpConn,
				eh:         eh,
				edh:        edh,
				ctx:        eCtx,
				cancelCtx:  eCancelCtx,
				writeChan:  make(chan []byte, 1024),
			}

			err := te.Serve()
			if err != nil {
				log.Println(err)
				return
			}
		}(tcpConn)
	}
}

func (t *tcpEdge) Id() common.Id {
	return t.id
}

func (t *tcpEdge) TunnelInfo() common.TunnelInfo {
	return t.tunnelInfo
}
func (t *tcpEdge) SendData(data []byte) error {
	if t.IsClosed() {
		return ErrEdgeClosed
	}
	t.writeChan <- data
	return nil
}

func (t *tcpEdge) IsClosed() bool {
	return t.ctx.Err() != nil
}

func (t *tcpEdge) Close() error {
	t.cancelCtx()
	_ = t.tcpConn.Close()
	return nil
}

func (t *tcpEdge) Serve() error {
	if t.isServed.CompareAndSwap(false, true) == false {
		return ErrIllegalState
	}

	var (
		wg = &sync.WaitGroup{}
	)

	wg.Add(2)
	t.safeEdgeHandler().OnEdgeCreated(t)
	defer func() {
		t.safeEdgeHandler().OnEdgeClosed(t)
		close(t.writeChan)
	}()

	// Reader
	go func() {
		defer wg.Done()

		for {
			if t.ctx.Err() != nil {
				return
			}
			buff := make([]byte, 4096)
			n, err := t.tcpConn.Read(buff)
			if err != nil {
				t.cancelCtx()
				return
			}

			t.safeEdgeDataHandler().OnEdgeDataReceived(t.Id(), buff[:n])
		}
	}()

	// Writer
	go func() {
		defer wg.Done()

		for {
			select {
			case <-t.ctx.Done():
				return
			case b, ok := <-t.writeChan:
				if !ok {
					// channel is closed
					t.cancelCtx()
					return
				}
				if b == nil {
					// todo log this
					continue
				}

				n, err := t.tcpConn.Write(b)
				if err != nil {
					t.cancelCtx()
					return
				}

				if n != len(b) {
					// todo log this
				}
			}
		}
	}()

	wg.Wait()
	return nil
}

func (t *tcpEdge) OnEdgeCreated(edge Edge) {
	// do nothing, to support safeEdgeHandler
}

func (t *tcpEdge) OnEdgeClosed(edge Edge) {
	// do nothing, to support safeEdgeHandler
}

func (t *tcpEdge) OnEdgeDataReceived(edgeId common.Id, data []byte) {
	// do nothing, to support safeEdgeDataHandler
}

func (t *tcpEdge) safeEdgeDataHandler() EdgeDataHandler {
	if t.edh == nil {
		return t
	}
	return t.edh
}

func (t *tcpEdge) safeEdgeHandler() EdgeHandler {
	if t.eh == nil {
		return t
	}
	return t.eh
}
