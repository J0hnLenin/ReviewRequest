package service

import (
	"math/rand"
	"slices"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func prContainsReviewer(pr *domain.PullRequest, userID string) bool {
	return slices.Contains(pr.ReviewersID,userID)
}

func validCandidate(pr *domain.PullRequest, u *domain.User) bool {
	return u.IsActive &&
		!prContainsReviewer(pr, u.ID) &&
		pr.AuthorID != u.ID
}

func newReviewer(t *domain.Team, pr *domain.PullRequest) *domain.User {
	candidates := make([]*domain.User, 0, len(t.Members))

	for _, member := range t.Members {
		if validCandidate(pr, member) {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	ind := rand.Intn(len(candidates))
	return candidates[ind]
}

func addReviewer(pr *domain.PullRequest, userID string) {
	pr.ReviewersID = append(pr.ReviewersID, userID)
}

func replaceReviewer(pr *domain.PullRequest, oldReviewerID string, newReviewerID string) error {
	ind := slices.Index(pr.ReviewersID, oldReviewerID)
	if ind == -1 {
		return domain.ErrNotAssigned
	}
	pr.ReviewersID[ind] = newReviewerID
	return nil
}

func fillReviewers(pr *domain.PullRequest, t *domain.Team) {
	for len(pr.ReviewersID) < domain.MaxReviewers {
		reviewer := newReviewer(t, pr)
		if reviewer == nil {
			break
		}
		addReviewer(pr, reviewer.ID)
	}
}

