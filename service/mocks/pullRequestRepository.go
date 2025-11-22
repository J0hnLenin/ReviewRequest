package mocks

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (m *MockRepository) GetPRByAuthor(ctx context.Context, authorID string) ([]*domain.PullRequest, error) {
	args := m.Called(ctx, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PullRequest), args.Error(1)
}

func (m *MockRepository) GetPRById(ctx context.Context, id string) (*domain.PullRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockRepository) GetPRAndTeam(ctx context.Context, id string) (*domain.PullRequest, *domain.Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*domain.PullRequest), nil, args.Error(2)
	}
	return args.Get(0).(*domain.PullRequest), args.Get(1).(*domain.Team), args.Error(2)
}

func (m *MockRepository) SavePR(ctx context.Context, pr *domain.PullRequest) error {
	args := m.Called(ctx, pr)
	return args.Error(0)
}