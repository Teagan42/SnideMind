package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"

	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/server/middleware"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

type Server struct {
	Config     config.ServerConfig
	HttpServer *http.Server
	Router     *mux.Router
}

type Params struct {
	fx.In
	Config    *config.Config
	Lifecycle fx.Lifecycle
}

type Result struct {
	fx.Out
	Server *Server     `name:"server"`
	Router *mux.Router `name:"rootRouter"`
}

func NewServer(p Params) (Result, error) {
	serverConfig := p.Config.Server
	addr := fmt.Sprintf("%s:%d", serverConfig.Bind, serverConfig.Port)
	server := &Server{
		Config: serverConfig,
		HttpServer: &http.Server{
			Addr: addr,
		},
		Router: mux.NewRouter(),
	}
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Starting server on %s", addr)
			go func() {
				if err := server.Start(ctx); err != nil {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Stopping server on %s", addr)
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down server: %v", err)
			}
			return nil
		},
	})

	return Result{
		Server: server,
		Router: server.Router,
	}, nil
}

func (s *Server) GetOpenAPIRouter() (*routers.Router, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("spec/openapi.yaml")
	if err != nil {
		log.Fatalf("failed to load OpenAPI spec: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		log.Fatalf("invalid OpenAPI spec: %v", err)
	}

	openAPIRouter, err := gorillamux.NewRouter(doc)
	if err != nil {
		log.Fatalf("failed to create OpenAPI router: %v", err)
	}

	return &openAPIRouter, nil
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.Config.Bind, s.Config.Port)
	log.Printf("Starting server on %s", addr)
	if openAPIRouter, err := s.GetOpenAPIRouter(); err != nil {
		log.Fatalf("failed to get OpenAPI router: %v", err)
	} else {
		s.Router.Use(middleware.OpenAPIValidationMiddleware(*openAPIRouter))
	}

	s.Router.PathPrefix("/v1").Subrouter()
	return s.HttpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("Stopping server on %s", s.HttpServer.Addr)
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
		return err
	}
	return nil
}
