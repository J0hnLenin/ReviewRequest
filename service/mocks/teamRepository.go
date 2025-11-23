package mocks

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (m *MockRepository) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockRepository) GetTeamByUser(ctx context.Context, userID string) (*domain.Team, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *MockRepository) SaveTeam(ctx context.Context, t *domain.Team) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockRepository) ChangeTeamActive(ctx context.Context, name string, active bool) (*domain.Team, error) {
    args := m.Called(ctx, name, active)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Team), args.Error(1)
}