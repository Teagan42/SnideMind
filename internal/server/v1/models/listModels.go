package models

import (
	"encoding/json"
	"net/http"

	"github.com/teagan42/snidemind/internal/types"
)

func ListModelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var response = types.ModelListResponse{}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
