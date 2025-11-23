package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	return s.repo.GetStatistics(ctx)
}