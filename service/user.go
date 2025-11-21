package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) UserChangeActive(ctx context.Context, id string, newValue bool) (*domain.User, error) {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrNotFound
	}
	if user.IsActive == newValue{
		return user, nil
	}
	user.IsActive = newValue
	err = s.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UserGetReviews(ctx context.Context, id string) ([]*domain.PullRequest, error) {
	return s.prRepo.GetPRByAuthor(ctx, id)
}