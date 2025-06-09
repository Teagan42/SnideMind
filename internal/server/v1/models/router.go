package models

import (
	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/internal/llm"
)

func AddRoutes(r *mux.Router, llmClient *llm.LLM) *mux.Router {
	r.HandleFunc("/", ListModelsHandler).Methods("GET")
	r.HandleFunc("/{modelID}", GetModelHandler).Methods("GET")
	return r
}
