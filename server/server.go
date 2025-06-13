package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

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
	Logger     *zap.Logger
}

type Params struct {
	fx.In
	Config    *config.Config
	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

type Result struct {
	fx.Out
	Server *Server
	Router *mux.Router `name:"rootRouter"`
}

func (w *Server) WalkRoutes() {
	err := w.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var fields = []zap.Field{}
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fields = append(fields, zap.String("pathTemplate", pathTemplate))
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fields = append(fields, zap.String("pathRegexp", pathRegexp))
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fields = append(fields, zap.Strings("queriesTemplates", queriesTemplates))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fields = append(fields, zap.Strings("queriesRegexps", queriesRegexps))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fields = append(fields, zap.Strings("methods", methods))
		}
		name := route.GetName()
		if name != "" {
			fields = append(fields, zap.String("name", name))
		}
		w.Logger.Info("Route found", fields...)
		return nil
	})

	if err != nil {
		w.Logger.Error("Error walking routes", zap.Error(err))
	} else {
		w.Logger.Info("Finished walking routes")
	}
}

func NewServer(p Params) (Result, error) {
	serverConfig := p.Config.Server
	bindAddr := ""
	if serverConfig.Bind != nil && *serverConfig.Bind != "" {
		bindAddr = string(*serverConfig.Bind)
	}
	addr := fmt.Sprintf("%s:%d", bindAddr, serverConfig.Port)
	httpServer := http.Server{
		Addr: addr,
	}
	server := &Server{
		Config:     serverConfig,
		HttpServer: &httpServer,
		Router:     mux.NewRouter(),
		Logger:     p.Logger.Named("Server"),
	}
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			server.WalkRoutes()
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

	// server.Router.Use(middleware.LogRequestMiddleware())
	server.HttpServer.Handler = server.Router

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
	if openAPIRouter, err := s.GetOpenAPIRouter(); err != nil {
		log.Fatalf("failed to get OpenAPI router: %v", err)
	} else {
		s.Router.Use(middleware.OpenAPIValidationMiddleware(*openAPIRouter))
	}
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
