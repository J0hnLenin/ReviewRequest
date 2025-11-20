package service

import (
	"context"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) PRCreate(ctx context.Context, prID string, title string, authorID string) (*domain.PullRequest, error) {
	_, err := s.prRepo.GetById(ctx, prID)
	if err != domain.ErrNotFound {
		return nil, domain.ErrPRExists
	}
	author, err := s.userRepo.GetById(ctx, authorID)
	if err != nil {
		return nil, err
	}
	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}
	pr := &domain.PullRequest{
		ID:     prID,
		Title:  title,
		Author: author,
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
	if !prContainsReviewer(pr, reviewer) {
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
