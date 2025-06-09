package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/teagan42/snidemind/internal/server/v1/chat"
	"github.com/teagan42/snidemind/internal/server/v1/models"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Mount("/chat", chat.Router())
	r.Mount("/models", models.Router())

	return r
}
