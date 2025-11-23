package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

type Repository interface {
	GetTeamByName(ctx context.Context, name string) (*domain.Team, error)
	GetTeamByUser(ctx context.Context, userID string) (*domain.Team, error)
	SaveTeam(ctx context.Context, t *domain.Team) error
	ChangeTeamActive(ctx context.Context, name string, active bool) (*domain.Team, error)

	GetUserById(ctx context.Context, id string) (*domain.User, error)
	SaveUser(ctx context.Context, u *domain.User) error

	GetPRByAuthor(ctx context.Context, id string) ([]*domain.PullRequest, error)
	GetPRById(ctx context.Context, id string) (*domain.PullRequest, error)
	GetPRAndTeam(ctx context.Context, id string) (*domain.PullRequest, *domain.Team, error)
	SavePR(ctx context.Context, pr *domain.PullRequest) error

	GetStatistics(ctx context.Context) (*domain.Statistics, error)
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}
