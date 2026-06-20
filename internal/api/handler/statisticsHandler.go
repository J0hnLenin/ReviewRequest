package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	stats, err := h.service.GetStatistics(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		log.Printf("response encode error: %v", err)
	}
}
