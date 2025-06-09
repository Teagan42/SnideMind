package chat

import (
	"encoding/json"
	"net/http"

	"github.com/teagan42/snidemind/internal/server/middleware"
	"github.com/teagan42/snidemind/internal/types"
)

func ChatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, err := middleware.GetValidatedBody[*types.ChatCompletionRequest](r); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// TODO: Implement the logic to filter Tools/Prompts/Resources based on query and MCP blacklist
	// Forward request to the LLM

	var response = types.ChatCompletionResponse{}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
