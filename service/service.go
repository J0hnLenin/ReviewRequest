package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

type TeamRepository interface {
	GetTeamByName(ctx context.Context, name string) (*domain.Team, error)
	SaveTeam(ctx context.Context, t *domain.Team) error
}

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	SaveUser(ctx context.Context, u *domain.User) error
}

type PullRequestRepository interface {
	GetPRByAuthor(ctx context.Context, id string) ([]*domain.PullRequest, error)
	GetPRById(ctx context.Context, id string) (*domain.PullRequest, error)
	GetPRAndTeam(ctx context.Context, id string) (*domain.PullRequest, *domain.Team, error)
	SavePR(ctx context.Context, pr *domain.PullRequest) error
}

type Service struct {
	teamRepo TeamRepository
	userRepo UserRepository
	prRepo   PullRequestRepository
}

func NewService(t TeamRepository, u UserRepository, pr PullRequestRepository) *Service {
	return &Service{
		teamRepo: t,
		userRepo: u,
		prRepo:   pr,
	}
}