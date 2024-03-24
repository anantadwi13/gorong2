package backbone

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anantadwi13/gorong2/component/backbone"
	"github.com/anantadwi13/gorong2/pkg/pool"
	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	times := 3
	pingCount := &atomic.Int32{}
	pongCount := &atomic.Int32{}

	messageFactory := &ProtobufMessageFactory{}

	connFactory, err := NewWebsocketConnectionFactory(messageFactory)
	if err != nil {
		t.Error("error initializing connection factory", err)
		return
	}

	listen, err := connFactory.Listen("localhost:12345")
	if err != nil {
		t.Error(err)
		return
	}
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		defer listen.Close()

		t.Log("listening")
		for {
			conn, err := listen.AcceptController()
			if err != nil {
				if !errors.Is(err, backbone.ErrConnClosed) {
					t.Error(err)
				}
				return
			}
			go func() {
				defer listen.Close()
				defer conn.Close()
				defer t.Log("conn closed", conn.RemoteAddr())

				for x := times; x >= 0; x-- {
					msg, err := conn.ReadMessage()
					assert.NoError(t, err)

					message, ok := msg.(*backbone.PingMessage)
					assert.True(t, ok)
					if assert.False(t, message.Time.IsZero()) {
						pingCount.Add(1)
					}

					err = conn.WriteMessage(&backbone.PongMessage{
						Error: nil,
						Time:  time.Now(),
					})
					assert.NoError(t, err)
				}
			}()
		}
	}()

	go func() {
		defer wg.Done()

		conn, err := connFactory.DialController(fmt.Sprintf("ws://%s", listen.Addr()))
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		for x := times; x >= 0; x-- {
			err := conn.WriteMessage(&backbone.PingMessage{Time: time.Now()})
			assert.NoError(t, err)

			msg, err := conn.ReadMessage()
			assert.NoError(t, err)

			message, ok := msg.(*backbone.PongMessage)
			assert.True(t, ok)
			if assert.Nil(t, message.Error) && assert.False(t, message.Time.IsZero()) {
				pongCount.Add(1)
			}
		}
	}()
	wg.Wait()
	assert.EqualValues(t, pingCount.Load(), pongCount.Load())
}

func TestWorker(t *testing.T) {
	const writeSize = 1 << 20
	readSize := &atomic.Int32{}

	messageFactory := &ProtobufMessageFactory{}

	connFactory, err := NewWebsocketConnectionFactory(messageFactory)
	if err != nil {
		t.Error("error initializing connection factory", err)
		return
	}

	listen, err := connFactory.Listen("localhost:12345")
	if err != nil {
		t.Error(err)
		return
	}
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		defer listen.Close()

		t.Log("listening")
		for {
			conn, err := listen.AcceptWorker()
			if err != nil {
				if !errors.Is(err, backbone.ErrConnClosed) {
					t.Error(err)
				}
				return
			}
			go func() {
				defer listen.Close()
				defer conn.Close()
				defer t.Log("conn closed", conn.RemoteAddr())

				pipeR, pipeW := io.Pipe()
				defer pipeW.Close()

				go func() {
					size := 4096
					buf := pool.GetBuffer(size)
					defer pool.PutBuffer(buf)
					for {
						buf.Buf = buf.Buf[:size]
						n, err := pipeR.Read(buf.Buf)
						if err != nil {
							if !errors.Is(err, io.EOF) {
								t.Error("error read", err)
							}
							return
						}
						t.Log("read msg", string(buf.Buf[:n]))
					}
				}()

				defer func() {
					t.Log("total read", readSize.Load())
				}()

				bufSize := 256
				buf := pool.GetBuffer(bufSize)
				defer pool.PutBuffer(buf)

				for {
					t.Log("copying")

					// method 1
					//data, err := io.ReadAll(conn)
					//total += int64(len(data))
					//log.Info(ctx, "read", len(data), string(data))
					//if err != nil {
					//	return
					//}
					//return

					// method 2
					n, err := io.CopyBuffer(io.Discard, conn, buf.Buf[:bufSize])
					readSize.Add(int32(n))
					if err != nil {
						if errors.Is(err, backbone.ErrConnClosed) {
							return
						}
						t.Error("error read msg", n, err)
						return
					}
					t.Log("retry copy", n)
					time.Sleep(1 * time.Second)

					// method 3
					//buf.Buf = buf.Buf[:bufSize]
					//n2, err := conn.Read(buf.Buf)
					//total += int64(n2)
					//if err != nil {
					//	if errors.Is(err, backbone.ErrConnClosed) {
					//		log.Error(ctx, "closed", string(buf.Buf[:n2]), n2, err)
					//		return
					//	}
					//	log.Error(ctx, "unknown error", string(buf.Buf[:n2]), n2, err)
					//	return
					//}
					//log.Info(ctx, "read manual", string(buf.Buf[:n2]), n2, err)
				}
			}()
		}
	}()

	go func() {
		defer wg.Done()

		conn, err := connFactory.DialWorker(fmt.Sprintf("ws://%s", listen.Addr()))
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		bufSize := 512
		buf := pool.GetBuffer(bufSize)
		defer pool.PutBuffer(buf)

		buf.Buf = buf.Buf[:bufSize]
		n, err := io.CopyBuffer(conn, io.LimitReader(rand.Reader, writeSize), buf.Buf)
		t.Log("writing", n)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()
	assert.EqualValues(t, writeSize, readSize.Load())
}
