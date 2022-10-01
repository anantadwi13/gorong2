package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anantadwi13/gorong2/internal/agent"
	"github.com/anantadwi13/gorong2/internal/common"
	"github.com/spf13/pflag"
)

func main() {
	var (
		config     = agent.Config{}
		tunnelType string
	)
	pflag.ErrHelp = errors.New("")
	pflag.StringVarP(&config.Token, "token", "t", "", "Token")
	pflag.StringVarP(&config.BackboneAddr, "backbone", "b", "", "Backbone Address")
	pflag.BoolVar(&config.BackboneSecured, "backbone-secured", false, "Backbone Using Secure Protocol (TLS)")
	pflag.StringVar(&config.BackboneRootCAPath, "backbone-root-ca", "", "Backbone Root CA Path")
	pflag.StringVarP(&config.EdgeServerAddr, "server", "s", "", "Edge Server Address")
	pflag.StringVarP(&config.EdgeAgentAddr, "agent", "a", "", "Edge Agent Address")
	pflag.StringVarP(&config.TunnelId, "tunnelid", "i", "", "Tunnel ID")
	pflag.StringVar(&tunnelType, "tunneltype", "tcp", "Tunnel Type (tcp, udp)")
	pflag.Parse()
	switch tunnelType {
	case "tcp":
		config.TunnelType = common.TunnelTypeTcp
	case "udp":
		config.TunnelType = common.TunnelTypeUdp
	default:
		log.Fatalln("tunnel type is not valid")
	}

	if config.Token == "" {
		log.Fatalln("token is required")
	}

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	a, err := agent.NewAgent(config)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer close(stopChan)
		err = a.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	<-stopChan

	err = a.Shutdown(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
