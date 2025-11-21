package handler

import (
	"encoding/json"
	"net/http"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (h *Handler) TeamAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var req struct {
		TeamName string                   `json:"team_name"`
		Members  []map[string]interface{} `json:"members"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	team := &domain.Team{
		Name:    req.TeamName,
		Members: make([]*domain.User, len(req.Members)),
	}

	for i, member := range req.Members {
		userID, _ := member["user_id"].(string)
		username, _ := member["username"].(string)
		isActive, _ := member["is_active"].(bool)

		team.Members[i] = &domain.User{
			ID:       userID,
			Name:     username,
			TeamName: req.TeamName,
			IsActive: isActive,
		}
	}

	if err := h.service.TeamSave(r.Context(), team); err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"team": map[string]interface{}{
			"team_name": team.Name,
			"members":   h.convertMembersToResponse(team.Members),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) TeamGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		h.writeError(w, http.StatusBadRequest, "MISSING_PARAM", "team_name is required")
		return
	}

	team, err := h.service.TeamGetByName(r.Context(), teamName)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"team_name": team.Name,
		"members":   h.convertMembersToResponse(team.Members),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) convertMembersToResponse(members []*domain.User) []map[string]interface{} {
	result := make([]map[string]interface{}, len(members))
	for i, member := range members {
		result[i] = map[string]interface{}{
			"user_id":   member.ID,
			"username":  member.Name,
			"is_active": member.IsActive,
		}
	}
	return result
}