package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anantadwi13/gorong2/internal/legacy/server"
	"github.com/spf13/pflag"
)

func main() {
	config := server.Config{}
	pflag.ErrHelp = errors.New("")
	pflag.StringVarP(&config.Token, "token", "t", "", "Token (Default: Auto Generated)")
	pflag.StringVarP(&config.BackboneAddr, "backbone", "b", ":19000", "Backbone Address")
	pflag.StringVar(&config.BackboneRootCAPath, "backbone-root-ca", "", "Backbone Root CA Path")
	pflag.StringVar(&config.BackbonePublicKeyPath, "backbone-public-key", "", "Backbone Public Key Path")
	pflag.StringVar(&config.BackbonePrivateKeyPath, "backbone-private-key", "", "Backbone Private Key Path")
	pflag.Parse()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	s, err := server.NewServer(config)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer close(stopChan)
		err = s.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	<-stopChan

	err = s.Shutdown(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
