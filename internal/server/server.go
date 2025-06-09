package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/internal/models"
	"github.com/teagan42/snidemind/internal/server/middleware"
	v1 "github.com/teagan42/snidemind/internal/server/v1"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"

	"github.com/teagan42/snidemind/internal/llm"
)

type Server struct {
	Config models.ServerConfig
}

func NewServer(cfg models.ServerConfig) *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Start(llmClient *llm.LLM) error {
	fmt.Printf("Starting server on %s:%d\n", s.Config.Bind, s.Config.Port)
	addr := fmt.Sprintf("%s:%d", s.Config.Bind, s.Config.Port)
	r := mux.NewRouter()

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
	r.Use(middleware.OpenAPIValidationMiddleware(openAPIRouter))
	v1.AddRoutes(r.PathPrefix("/v1").Subrouter(), llmClient)

	err = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
		log.Fatalf("failed to walk routes: %v", err)
		panic(err)
	}

	fmt.Println("Server is running...")
	return http.ListenAndServe(addr, r)
}
