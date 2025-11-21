package domain

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
}

type PRStatus bool

const (
	Merged PRStatus = true
	Open   PRStatus = false
)

const (
	MaxReviewers = 2
)