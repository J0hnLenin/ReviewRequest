package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (h *Handler) PRCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	pr, err := h.service.PRCreate(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"pr": h.convertPRToResponse(pr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) PRMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	pr, err := h.service.PRMerge(r.Context(), req.PullRequestID); 
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"pr": h.convertPRToResponse(pr),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) PRReassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_reviewer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	pr, replacedBy, err := h.service.PRreassign(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"pr":          h.convertPRToResponse(pr),
		"replaced_by": replacedBy,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) convertPRToResponse(pr *domain.PullRequest) map[string]interface{} {
	status := "OPEN"
	if pr.Status == domain.Merged {
		status = "MERGED"
	}

	response := map[string]interface{}{
		"pull_request_id":   pr.ID,
		"pull_request_name": pr.Title,
		"author_id":         pr.AuthorID,
		"status":            status,
		"assigned_reviewers": pr.ReviewersID,
	}

	if pr.Status == domain.Merged && pr.MergedAt != nil {
		response["mergedAt"] = pr.MergedAt.Format(time.RFC3339)
	}
	
	return response
}

func (h *Handler) convertPRsToShortResponse(prs []*domain.PullRequest) []map[string]interface{} {
	result := make([]map[string]interface{}, len(prs))
	for i, pr := range prs {
		status := "OPEN"
		if pr.Status == domain.Merged {
			status = "MERGED"
		}

		result[i] = map[string]interface{}{
			"pull_request_id":   pr.ID,
			"pull_request_name": pr.Title,
			"author_id":         pr.AuthorID,
			"status":            status,
		}
	}
	return result
}