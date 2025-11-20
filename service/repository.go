package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

type TeamRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Team, error)
	Save(ctx context.Context, t *domain.Team) error
}

type UserRepository interface {
	GetById(ctx context.Context, id string) (*domain.User, error)
	Save(ctx context.Context, u *domain.User) error
}

type PullRequestRepository interface {
	GetByAuthor(ctx context.Context, id string) ([]*domain.PullRequest, error)
	GetById(ctx context.Context, id string) (*domain.PullRequest, error)
	Save(ctx context.Context, pr *domain.PullRequest) error
}