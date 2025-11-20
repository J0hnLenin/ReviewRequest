package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) TeamSave(ctx context.Context, t *domain.Team) error {
	_, err := s.teamRepo.GetByName(ctx, t.Name)
	if err == domain.ErrNotFound {
		return s.teamRepo.Save(ctx, t)
	}
	return domain.ErrTeamExists
}

func (s *Service) TeamGetByName(ctx context.Context, n string) (*domain.Team, error) {
	team, err := s.teamRepo.GetByName(ctx, n)
	if err != nil {
		return nil, err
	}
	return team, nil
}