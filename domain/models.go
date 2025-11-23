package domain

import "time"

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

type Team struct {
	Name    string
	Members []*User
}

type PullRequest struct {
	ID          string
	Title       string
	AuthorID    string
	ReviewersID []string
	Status      PRStatus
	MergedAt    *time.Time
}

type PRStatus bool

const (
	Merged PRStatus = true
	Open   PRStatus = false
)

const (
	MaxReviewers = 2
)

type Statistics struct {
	TotalOpenPRs     int          `json:"total_open_prs"`
	TotalClosedPRs   int          `json:"total_closed_prs"`
	TopOpenReviewer  *UserStats   `json:"top_open_reviewer,omitempty"`
	TopClosedReviewer *UserStats  `json:"top_closed_reviewer,omitempty"`
	TopAuthor        *UserStats   `json:"top_author,omitempty"`
}

type UserStats struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Count    int    `json:"count"`
}