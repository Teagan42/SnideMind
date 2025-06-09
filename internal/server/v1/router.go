package v1

import (
	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/internal/server/v1/chat"
	"github.com/teagan42/snidemind/internal/server/v1/models"
)

func AddRoutes(r *mux.Router) *mux.Router {
	chat.AddRoutes(r.PathPrefix("/chat").Subrouter())
	models.AddRoutes(r.PathPrefix("/models").Subrouter())

	return r
}
