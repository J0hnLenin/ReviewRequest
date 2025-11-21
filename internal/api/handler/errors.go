package handler

import (
	"encoding/json"
	"net/http"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *Handler) writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	errorResp := ErrorResponse{}
	errorResp.Error.Code = code
	errorResp.Error.Message = message
	
	json.NewEncoder(w).Encode(errorResp)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrNotFound:
		h.writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case domain.ErrTeamExists:
		h.writeError(w, http.StatusBadRequest, "TEAM_EXISTS", err.Error())
	case domain.ErrPRExists:
		h.writeError(w, http.StatusConflict, "PR_EXISTS", err.Error())
	case domain.ErrPRMerged:
		h.writeError(w, http.StatusConflict, "PR_MERGED", err.Error())
	case domain.ErrNotAssigned:
		h.writeError(w, http.StatusConflict, "NOT_ASSIGNED", err.Error())
	case domain.ErrNoCandidate:
		h.writeError(w, http.StatusConflict, "NO_CANDIDATE", err.Error())
	default:
		h.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}