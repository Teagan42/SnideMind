package v1

import (
	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/internal/llm"
	"github.com/teagan42/snidemind/internal/server/v1/chat"
	"github.com/teagan42/snidemind/internal/server/v1/models"
)

func AddRoutes(r *mux.Router, llmClient *llm.LLM) *mux.Router {
	chat.AddRoutes(r.PathPrefix("/chat").Subrouter(), llmClient)
	models.AddRoutes(r.PathPrefix("/models").Subrouter(), llmClient)

	return r
}
