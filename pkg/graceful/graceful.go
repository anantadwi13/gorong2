package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func ListenTermination(graceTimeout time.Duration, signals ...os.Signal) (
	runningCtx context.Context, onShutdownListener *OnShutdownListener, cleanup func(onForceShutdown func()),
) {
	runningCtx, cancel := context.WithCancel(context.Background())
	waitChan := make(chan struct{})
	stopChan := make(chan os.Signal, 1)
	onShutdownListener = &OnShutdownListener{}
	forceShutdown := int32(0)
	cleanup = func(onForceShutdown func()) {
		cancel()
		<-waitChan
		if atomic.LoadInt32(&forceShutdown) == 1 && onForceShutdown != nil {
			onForceShutdown()
		}
	}

	if len(signals) == 0 {
		signals = append(signals, syscall.SIGTERM)
	}
	signal.Notify(stopChan, signals...)

	go func() {
		select {
		case <-stopChan:
		case <-runningCtx.Done():
		}

		cancel()

		if graceTimeout <= 0 {
			graceTimeout = 30 * time.Second
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), graceTimeout)
		defer cancel()

		done := onShutdownListener.execute(shutdownCtx)

		select {
		case <-done:
		case <-stopChan: // force shutdown
			atomic.CompareAndSwapInt32(&forceShutdown, 0, 1)
		case <-shutdownCtx.Done(): // timeout reached
			select {
			case <-done:
			case <-time.After(graceTimeout / 10): // wait longer time to make sure all resources are freed up
			case <-stopChan: // force shutdown after grace timeout
				atomic.CompareAndSwapInt32(&forceShutdown, 0, 1)
			}
		}

		close(waitChan)
	}()

	return
}

type Service interface {
	Shutdown(ctx context.Context) error
}

type OnShutdownFunc func(ctx context.Context) error

func (o OnShutdownFunc) Shutdown(ctx context.Context) error {
	return o(ctx)
}

type service struct {
	s       Service
	onError func(err error)
}

type OnShutdownListener struct {
	services []service
	mu       sync.RWMutex
	executed int32
}

func (l *OnShutdownListener) RegisterOnShutdown(svc Service, onError func(err error)) {
	if atomic.LoadInt32(&l.executed) == 1 {
		return
	}
	if svc == nil {
		return
	}
	if onError == nil {
		onError = func(err error) {}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.services = append(l.services, service{
		s:       svc,
		onError: onError,
	})
}

func (l *OnShutdownListener) execute(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	if !atomic.CompareAndSwapInt32(&l.executed, 0, 1) {
		close(done)
		return done
	}

	wg := sync.WaitGroup{}
	go func() {
		wg.Wait()
		close(done)
	}()

	l.mu.RLock()
	defer l.mu.RUnlock()

	wg.Add(len(l.services))
	for _, svc := range l.services {
		go func(svc service) {
			defer wg.Done()
			err := svc.s.Shutdown(ctx)
			if err != nil {
				svc.onError(err)
				return
			}
		}(svc)
	}

	return done
}
