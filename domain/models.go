package domain

type User struct {
	ID uint64
	Name string
	Team *Team 
	IsActive bool
}

type Team struct {
	ID uint64
	Name string
}

type PullRequest struct {
	ID uint64
	Title string
	Author *User
	Merged bool 
}

type Reviewer struct {
	User *User
	PR *PullRequest
}
