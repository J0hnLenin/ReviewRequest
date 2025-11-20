package service

import (
	"math/rand"
	"slices"

	"github.com/J0hnLenin/ReviewRequest/domain"
)

func prContainsReviewer(pr *domain.PullRequest, u *domain.User) bool {
	return slices.ContainsFunc(pr.Reviewers, func(r *domain.User) bool { 
                return userEquals(r, u) 
            })
}

func validCandidate(pr *domain.PullRequest, u *domain.User) bool {
	return u.IsActive &&
		!prContainsReviewer(pr, u) &&
		!userEquals(pr.Author, u)
}

func newReviewer(t *domain.Team, pr *domain.PullRequest) *domain.User {
	candidates := make([]*domain.User, 0, len(pr.Reviewers))

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

func addReviewer(pr *domain.PullRequest, u *domain.User) {
	pr.Reviewers = append(pr.Reviewers, u)
}

func replaceReviewer(pr *domain.PullRequest, oldReviewer *domain.User, newReviewer *domain.User) {
	ind := slices.IndexFunc(pr.Reviewers, func(r *domain.User) bool {
		return userEquals(r, oldReviewer)
	})
	pr.Reviewers[ind] = newReviewer
}



func fillReviewers(pr *domain.PullRequest, t *domain.Team) {
	for len(pr.Reviewers) < domain.MaxReviewers {
		reviewer := newReviewer(t, pr)
		if reviewer == nil {
			break
		}
		addReviewer(pr, reviewer)
	}
}

