package chat

import (
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
	llm, ok := middleware.GetLLMFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"LLM not found in context"}`, http.StatusInternalServerError)
		return
	}
	chatCompletion, err := middleware.GetValidatedBody[types.ChatCompletionRequest](r)
	if err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	err = llm.CallCompletion(chatCompletion, w)
	if err != nil {
		http.Error(w, `{"error":"failed to call LLM"}`, http.StatusInternalServerError)
		return
	}

}
