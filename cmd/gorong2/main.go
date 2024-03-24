package main

import (
	"syscall"
	"time"

	"github.com/anantadwi13/gorong2/internal/backbone"
	"github.com/anantadwi13/gorong2/internal/server"
	"github.com/anantadwi13/gorong2/pkg/graceful"
	"github.com/anantadwi13/gorong2/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.InitLogrusLogger(logrusLogger)

	ctx, listener, cleanup := graceful.ListenTermination(5*time.Second, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)
	defer log.Info(ctx, "success")
	defer cleanup(func() {
		log.Info(ctx, "force stopped")
	})

	log.Info(ctx, "running")

	messageFactory := &backbone.ProtobufMessageFactory{}

	connFactory, err := backbone.NewWebsocketConnectionFactory(messageFactory)
	if err != nil {
		log.Error(ctx, "error initializing connection factory", err)
		return
	}

	svr, err := server.NewServer(connFactory)
	if err != nil {
		log.Error(ctx, "error initializing server", err)
		return
	}
	listener.RegisterOnShutdown(svr, nil)

	<-ctx.Done()
	log.Info(nil, "context done")
}
