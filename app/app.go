package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/logger"
	"github.com/teagan42/snidemind/pipeline"
	"github.com/teagan42/snidemind/server"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	BindAddress *config.Host `name:"bindAddress"`
	PortNumber  *config.Port `name:"port"`
	ConfigPath  string       `name:"configPath"`
}

func ParseArgs() Result {
	parser := argparse.NewParser("snidemind", "Snidemind AI Tool CLI Daemon")
	bind := parser.String("b", "bind", &argparse.Options{
		Required: false,
		Help:     "Address to bind the server to.",
	})
	port := parser.Int("p", "port", &argparse.Options{
		Required: false,
		Help:     "Port to run the server on.",
	})
	configPath := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     "Path to the configuration file.",
		Default:  "config.yaml",
	})
	if err := parser.Parse(os.Args); err != nil {
		log.Fatalf("[Argparse] Failed to parse arguments: %v", err)
		panic(err)
	}

	var bindAddress *config.Host
	var portNumber *config.Port
	if bind != nil && *bind != "" {
		if bindHost, err := config.NewHost(*bind); err != nil {
			log.Fatalf("[Bind] Invalid bind address: %v", err)
			panic(err)
		} else {
			bindAddress = &bindHost
		}
	}
	if port != nil && *port > 0 {
		if portNum, err := config.NewPort(*port); err != nil {
			log.Fatalf("[Port] Invalid port number: %v", err)
			panic(err)
		} else {
			portNumber = &portNum
		}
	}
	return Result{
		BindAddress: bindAddress,
		PortNumber:  portNumber,
		ConfigPath:  *configPath,
	}
}

var Module = fx.Module(
	"app",
	fx.Provide(
		ParseArgs,
	),
)

func main() {
	app := fx.New(
		Module,
		logger.Module,
		config.Module,
		pipeline.Module,
		server.Module,
	)

	app.Run()
}
