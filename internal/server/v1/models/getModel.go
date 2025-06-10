package models

import (
	"encoding/json"
	"net/http"

	"github.com/teagan42/snidemind/internal/models"
	"github.com/teagan42/snidemind/internal/server/middleware"
)

func GetModelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if modelID, err := middleware.GetRouteParams[string](r); err != nil {
		http.Error(w, `{"error":"invalid model ID"}`, http.StatusBadRequest)
		return
	} else if modelID == "" {
		http.Error(w, `{"error":"model ID is required"}`, http.StatusBadRequest)
		return
	}

	var response = models.ModelResponse{}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
