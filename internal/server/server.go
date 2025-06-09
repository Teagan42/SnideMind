package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teagan42/snidemind/internal/models"
	"github.com/teagan42/snidemind/internal/server/middleware"
	v1 "github.com/teagan42/snidemind/internal/server/v1"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/legacy"
)

type Server struct {
	Config models.ServerConfig
}

func NewServer(cfg models.ServerConfig) *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Start() error {
	fmt.Printf("Starting server on %s:%d\n", s.Config.Bind, s.Config.Port)
	addr := fmt.Sprintf("%s:%d", s.Config.Bind, s.Config.Port)
	r := chi.NewRouter()

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("spec/openapi.yaml")
	if err != nil {
		log.Fatalf("failed to load OpenAPI spec: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		log.Fatalf("invalid OpenAPI spec: %v", err)
	}

	openAPIRouter, err := legacy.NewRouter(doc)
	if err != nil {
		log.Fatalf("failed to create OpenAPI router: %v", err)
	}
	r.Use(middleware.OpenAPIValidationMiddleware(openAPIRouter))
	r.Mount("/v1", v1.Router())
	fmt.Println("Server is running...")
	return http.ListenAndServe(addr, r)
}
