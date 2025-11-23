package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) TeamSave(ctx context.Context, t *domain.Team) error {
	team, err := s.repo.GetTeamByName(ctx, t.Name)
	if err != nil {
		return err
	}
	if team != nil {
		return domain.ErrTeamExists
	}
	return s.repo.SaveTeam(ctx, t)
}

func (s *Service) TeamGetByName(ctx context.Context, n string) (*domain.Team, error) {
	team, err := s.repo.GetTeamByName(ctx, n)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrNotFound
	}
	return team, nil
}

func (s *Service) TeamChangeActive(ctx context.Context, name string, newValue bool) (*domain.Team, error) {
	team, err := s.repo.ChangeTeamActive(ctx, name, newValue)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrNotFound
	}
	return team, nil
}