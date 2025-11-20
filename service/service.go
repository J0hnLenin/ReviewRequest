package service

import (
	"context"
	"slices"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

type Service struct {
	teamRepo TeamRepository
	userRepo UserRepository
	prRepo   PullRequestRepository
}

func NewService(t TeamRepository, u UserRepository, pr PullRequestRepository) *Service {
	return &Service{
		teamRepo: t,
		userRepo: u,
		prRepo:   pr,
	}
}

func (s *Service) TeamSave(ctx context.Context, t *domain.Team) error {
	_, err := s.teamRepo.GetByName(ctx, t.Name)
	if err == domain.ErrNotFound {
		return s.teamRepo.Save(ctx, t)
	}
	return domain.ErrTeamExists
}

func (s *Service) TeamGetByName(ctx context.Context, n string) (*domain.Team, error) {
	t, err := s.teamRepo.GetByName(ctx, n)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) UserChangeActive(ctx context.Context, id string, newValue bool) (*domain.User, error) {
	u, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.IsActive == newValue{
		return u, nil
	}
	u.IsActive = newValue
	err = s.userRepo.Save(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) UserGetReviews(ctx context.Context, id string) ([]*domain.PullRequest, error) {
	return s.prRepo.GetByAuthor(ctx, id)
}

func (s *Service) PRCreate(ctx context.Context, pr *domain.PullRequest) error {
	_, err := s.prRepo.GetById(ctx, pr.ID)
	if err == domain.ErrNotFound {
		return s.prRepo.Save(ctx, pr)
	}
	return domain.ErrPRExists
}

func (s *Service) PRMerge(ctx context.Context, id string) error {
	pr, err := s.prRepo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if pr.Status == domain.Merged {
		return nil
	}
	pr.Status = domain.Merged
	return s.prRepo.Save(ctx, pr)
}

func (s *Service) PRreassign(ctx context.Context, prID string, reviewerID string) error {
	pr, err := s.prRepo.GetById(ctx, prID)
	if err != nil {
		return err
	}
	if pr.Status == domain.Merged {
		return domain.ErrPRMerged
	}
	reviewer, err := s.userRepo.GetById(ctx, reviewerID)
	if err != nil {
		return err
	}
	if !slices.Contains(pr.Reviewers, reviewer) {
		return domain.ErrNotAssigned
	}
	t, err := s.teamRepo.GetByName(ctx, pr.Author.TeamName)
	if err != nil {
		return err
	}
	newReviewer := getRandomTeamMember(t, pr)
	if newReviewer == nil {
		return domain.ErrNoCandidate
	}
	replaceReviewer(pr, reviewer, newReviewer)
	return s.prRepo.Save(ctx, pr)
}

func getRandomTeamMember(t *domain.Team, pr *domain.PullRequest) *domain.User{
	for _, u := range t.Members {
		if !slices.Contains(pr.Reviewers, u) && pr.Author!=u {
			return u
		}
	}
	return nil
}

func replaceReviewer(pr *domain.PullRequest, oldR *domain.User, newR *domain.User) {
	ind := slices.Index(pr.Reviewers, oldR)
	pr.Reviewers[ind] = newR
}