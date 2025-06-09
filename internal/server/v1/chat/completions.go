package chat

import (
	"encoding/json"
	"net/http"

	"github.com/teagan42/snidemind/internal/server/middleware"
)

func ChatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, err := middleware.GetValidatedBody[*ChatCompletionRequest](r); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	var response = ChatCompletionResponse{}

	// Handle chat completions logic here
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
