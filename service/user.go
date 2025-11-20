package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) UserChangeActive(ctx context.Context, id string, newValue bool) (*domain.User, error) {
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user.IsActive == newValue{
		return user, nil
	}
	user.IsActive = newValue
	err = s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UserGetReviews(ctx context.Context, id string) ([]*domain.PullRequest, error) {
	return s.prRepo.GetByAuthor(ctx, id)
}

func userEquals(u1, u2 *domain.User) bool {
	if u1 == nil || u2 == nil {
		return u1 == u2
	}
	return u1.ID == u2.ID
}