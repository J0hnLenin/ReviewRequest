package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) PRCreate(ctx context.Context, prID string, title string, authorID string) (*domain.PullRequest, error) {
	pr, team, err := s.prRepo.GetPRAndTeam(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr != nil {
		return nil, domain.ErrPRExists
	}
	if team == nil {
		return nil, domain.ErrNotFound
	}
	
	pr = &domain.PullRequest{
		ID:     prID,
		Title:  title,
		AuthorID: authorID,
		Status: domain.Open,
	}
	fillReviewers(pr, team)
	err = s.prRepo.Save(ctx, pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) PRMerge(ctx context.Context, id string) error {
	pr, err := s.prRepo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if pr == nil {
		return domain.ErrNotFound
	}
	if pr.Status == domain.Merged {
		return nil
	}
	pr.Status = domain.Merged
	return s.prRepo.Save(ctx, pr)
}

func (s *Service) PRreassign(ctx context.Context, prID string, reviewerID string) error {
	pr, team, err := s.prRepo.GetPRAndTeam(ctx, prID)
	if err != nil {
		return err
	}
	if pr == nil || team == nil {
		return domain.ErrNotFound
	}
	if pr.Status == domain.Merged {
		return domain.ErrPRMerged
	}
	if !prContainsReviewer(pr, reviewerID) {
		return domain.ErrNotAssigned
	}
	newReviewer := newReviewer(team, pr)
	if newReviewer == nil {
		return domain.ErrNoCandidate
	}
	err = replaceReviewer(pr, reviewerID, newReviewer.ID)
	if err != nil {
		return err
	}
	return s.prRepo.Save(ctx, pr)
}
