package service

import (
	"context"
	"math/rand"
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
	team, err := s.teamRepo.GetByName(ctx, n)
	if err != nil {
		return nil, err
	}
	return team, nil
}

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

func (s *Service) PRCreate(ctx context.Context, id string) error {
	pr, err := s.prRepo.GetById(ctx, id)
	if err != domain.ErrNotFound {
		return domain.ErrPRExists
	}
	team, err := s.teamRepo.GetByName(ctx, pr.Author.TeamName)
	if err != nil {
		return err
	}
	fillReviewers(pr, team)
	return s.prRepo.Save(ctx, pr)
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
	team, err := s.teamRepo.GetByName(ctx, pr.Author.TeamName)
	if err != nil {
		return err
	}
	newReviewer := newReviewer(team, pr)
	if newReviewer == nil {
		return domain.ErrNoCandidate
	}
	replaceReviewer(pr, reviewer, newReviewer)
	return s.prRepo.Save(ctx, pr)
}

func newReviewer(t *domain.Team, pr *domain.PullRequest) *domain.User{
	candidates := make([]*domain.User, len(pr.Reviewers))
	
	for _, member := range t.Members {
		if !slices.Contains(pr.Reviewers, member) && pr.Author != member {
			candidates = append(candidates, member)
		}
	}
	
	if len(candidates) == 0 {
		return nil
	}
	
	ind := rand.Intn(len(candidates))
	return candidates[ind]
}

func replaceReviewer(pr *domain.PullRequest, oldReviewer *domain.User, newReviewer *domain.User) {
	ind := slices.Index(pr.Reviewers, oldReviewer)
	pr.Reviewers[ind] = newReviewer
} 

func addReviewer(pr *domain.PullRequest, u *domain.User) {
	pr.Reviewers = append(pr.Reviewers, u)
}

func fillReviewers(pr *domain.PullRequest, t *domain.Team) {
	for range domain.MaxReviewers {
		reviewer := newReviewer(t, pr)
		if reviewer == nil {
			break
		}
		addReviewer(pr, reviewer)
	}
}