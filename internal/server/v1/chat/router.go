package chat

import (
	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/internal/llm"
)

func AddRoutes(r *mux.Router, llm *llm.LLM) *mux.Router {
	r.HandleFunc("/completions", ChatCompletionsHandler).Methods("POST")
	return r
}
