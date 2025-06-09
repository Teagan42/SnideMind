package models

import (
	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/", ListModelsHandler).Methods("GET")
	r.HandleFunc("/{modelID}", GetModelHandler).Methods("GET")
	return r
}
