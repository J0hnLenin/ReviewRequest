package service

import (
	"testing"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/stretchr/testify/assert"
)

func TestPrContainsReviewer(t *testing.T) {
	pr := &domain.PullRequest{
		ReviewersID: []string{"user1", "user2"},
	}

	testCases := []struct {
		name     string
		userID   string
		expected bool
	}{
		{"Reviewer exists", "user1", true},
		{"Reviewer exists", "user2", true},
		{"Reviewer not exists", "user3", false},
		{"Empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := prContainsReviewer(pr, tc.userID)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidCandidate(t *testing.T) {
	pr := &domain.PullRequest{
		AuthorID:    "author1",
		ReviewersID: []string{"reviewer1"},
	}

	user := &domain.User{
		ID:       "candidate1",
		IsActive: true,
	}

	assert.True(t, validCandidate(pr, user))

	user.ID = "author1"
	assert.False(t, validCandidate(pr, user))

	user.ID = "reviewer1"
	assert.False(t, validCandidate(pr, user))

	user.ID = "candidate1"
	user.IsActive = false
	assert.False(t, validCandidate(pr, user))
}

func TestFillReviewers_ZeroReviewers(t *testing.T) {
	pr := &domain.PullRequest{
		AuthorID:    "author1",
		ReviewersID: []string{},
	}

	team := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "author1", Name: "Author", TeamName: "test-team", IsActive: true},
		},
	}

	fillReviewers(pr, team)

	assert.Len(t, pr.ReviewersID, 0)
	assert.NotContains(t, pr.ReviewersID, "author1")
}

func TestFillReviewers_OneReviewer(t *testing.T) {
	pr := &domain.PullRequest{
		AuthorID:    "author1",
		ReviewersID: []string{},
	}

	team := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "author1", Name: "Author", TeamName: "test-team", IsActive: true},
			{ID: "user2", Name: "User 2", TeamName: "test-team", IsActive: true},
			{ID: "user4", Name: "User 3", TeamName: "test-team", IsActive: false},
		},
	}

	fillReviewers(pr, team)

	assert.Len(t, pr.ReviewersID, 1)
	assert.Contains(t, pr.ReviewersID, "user2")
	assert.NotContains(t, pr.ReviewersID, "author1")
	assert.NotContains(t, pr.ReviewersID, "user3")
}

func TestFillReviewers_TwoReviewers(t *testing.T) {
	pr := &domain.PullRequest{
		AuthorID:    "author1",
		ReviewersID: []string{},
	}

	team := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "author1", Name: "Author", TeamName: "test-team", IsActive: true},
			{ID: "user2", Name: "User 2", TeamName: "test-team", IsActive: true},
			{ID: "user3", Name: "User 3", TeamName: "test-team", IsActive: true},
			{ID: "user4", Name: "User 4", TeamName: "test-team", IsActive: false},
		},
	}

	fillReviewers(pr, team)

	assert.Len(t, pr.ReviewersID, 2)
	assert.NotContains(t, pr.ReviewersID, "author1")
	assert.NotContains(t, pr.ReviewersID, "user4")
}

func TestReplaceReviewer_Sucsess(t *testing.T) {
	pr := &domain.PullRequest{
		AuthorID:    "author1",
		ReviewersID: []string{"old_reviewer", "need_to_replace_reviewer",},
	}
	err := replaceReviewer(pr, "need_to_replace_reviewer", "new_reviewer")
	assert.NoError(t, err)
	assert.ElementsMatch(t, pr.ReviewersID, []string{"old_reviewer", "new_reviewer",})
}	

func TestReplaceReviewer_NotAssigned(t *testing.T) {
	pr := &domain.PullRequest{
		AuthorID:    "author1",
		ReviewersID: []string{},
	}
	err := replaceReviewer(pr, "user1", "user2")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotAssigned, err)
}	