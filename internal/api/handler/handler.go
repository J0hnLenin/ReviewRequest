package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/J0hnLenin/ReviewRequest/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
	if err != nil {
		log.Printf("response encode error: %v", err)
	}
}
