package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) UserSetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	user, err := h.service.UserChangeActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"user_id":   user.ID,
			"username":  user.Name,
			"team_name": user.TeamName,
			"is_active": user.IsActive,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UserGetReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.writeError(w, http.StatusBadRequest, "MISSING_PARAM", "user_id is required")
		return
	}

	prs, err := h.service.UserGetReviews(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"user_id": userID,
		"pull_requests": h.convertPRsToShortResponse(prs),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}