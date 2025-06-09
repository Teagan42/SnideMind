package chat

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/completions", ChatCompletionsHandler)
	return r
}
