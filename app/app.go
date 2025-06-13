package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/logger"
	"github.com/teagan42/snidemind/pipeline"
	"github.com/teagan42/snidemind/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, pipeline *pipeline.Pipeline, logger *zap.Logger, server *server.Server) {
			err := server.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
				pathTemplate, err := route.GetPathTemplate()
				if err == nil {
					fmt.Println("ROUTE:", pathTemplate)
				}
				pathRegexp, err := route.GetPathRegexp()
				if err == nil {
					fmt.Println("Path regexp:", pathRegexp)
				}
				queriesTemplates, err := route.GetQueriesTemplates()
				if err == nil {
					fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
				}
				queriesRegexps, err := route.GetQueriesRegexp()
				if err == nil {
					fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
				}
				methods, err := route.GetMethods()
				if err == nil {
					fmt.Println("Methods:", strings.Join(methods, ","))
				}
				fmt.Println()
				return nil
			})

			if err != nil {
				fmt.Println(err)
			}
			lc.Append(fx.Hook{
				OnStart: func(context context.Context) error {
					logger.Info("Starting Snidemind server...")
					return nil
				},
				OnStop: func(context context.Context) error {
					logger.Info("Stopping Snidemind server...")
					return nil
				},
			})
		}),
	)

	app.Run()
}
