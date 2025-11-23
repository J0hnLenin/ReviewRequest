package mocks

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (m *MockRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Statistics), args.Error(1)
}