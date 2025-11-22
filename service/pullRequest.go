package service

import (
	"context"
	"time"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func (s *Service) PRCreate(ctx context.Context, prID string, title string, authorID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetPRById(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr != nil {
		return nil, domain.ErrPRExists
	}

	team, err := s.teamRepo.GetTeamByUser(ctx, authorID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrNotFound
	}
	
	pr = &domain.PullRequest{
		ID:       prID,
		Title:    title,
		AuthorID: authorID,
		Status:   domain.Open,
		ReviewersID: make([]string, 0, 2),
		MergedAt: nil,
	}
	fillReviewers(pr, team)
	err = s.prRepo.SavePR(ctx, pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) PRMerge(ctx context.Context, id string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetPRById(ctx, id)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrNotFound
	}
	if pr.Status == domain.Merged {
		return pr, nil
	}
	
	now := time.Now()
	pr.Status = domain.Merged
	pr.MergedAt = &now
	
	err = s.prRepo.SavePR(ctx, pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) PRreassign(ctx context.Context, prID string, reviewerID string) (*domain.PullRequest, string, error) {
	pr, team, err := s.prRepo.GetPRAndTeam(ctx, prID)
	if err != nil {
		return nil, "", err
	}
	if pr == nil || team == nil {
		return nil, "", domain.ErrNotFound
	}
	if pr.Status == domain.Merged {
		return nil, "", domain.ErrPRMerged
	}
	if !prContainsReviewer(pr, reviewerID) {
		return nil, "", domain.ErrNotAssigned
	}
	newReviewer := newReviewer(team, pr)
	if newReviewer == nil {
		return nil, "", domain.ErrNoCandidate
	}
	err = replaceReviewer(pr, reviewerID, newReviewer.ID)
	if err != nil {
		return nil, "", err
	}
	err = s.prRepo.SavePR(ctx, pr)
	if err != nil {
		return nil, "", err
	}
	return pr, newReviewer.ID, err
}