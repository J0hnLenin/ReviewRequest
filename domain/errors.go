package domain

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrTeamExists    = errors.New("team_name already exists")
	ErrPRExists 	 = errors.New("PR id already exists")
	ErrPRMerged		 = errors.New("cannot reassign on merged PR")
	ErrNotAssigned   = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate   = errors.New("no active replacement candidate in team")
)