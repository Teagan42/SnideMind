package models

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", ListModelsHandler)
	r.Get("/{modelID}", GetModelHandler)
	return r
}
