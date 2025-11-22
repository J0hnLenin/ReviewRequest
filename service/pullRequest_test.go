package service

import (
	"context"
	"testing"
	"time"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPRCreate_Success(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"

	team := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "user1", Name: "Author", TeamName: "test-team", IsActive: true},
			{ID: "user2", Name: "Reviewer 1", TeamName: "test-team", IsActive: true},
			{ID: "user3", Name: "Reviewer 2", TeamName: "test-team", IsActive: true},
		},
	}

	mockPRRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)
	mockTeamRepo.On("GetTeamByUser", mock.Anything, authorID).Return(team, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.ID == prID &&
			pr.Title == title &&
			pr.AuthorID == authorID &&
			len(pr.ReviewersID) == 2 &&
			pr.Status == domain.Open
	})).Return(nil)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, prID, pr.ID)
	assert.Equal(t, title, pr.Title)
	assert.Equal(t, authorID, pr.AuthorID)
	assert.Len(t, pr.ReviewersID, 2)
	assert.Equal(t, domain.Open, pr.Status)

	mockPRRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestPRCreate_TeamNotFound(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"

	mockPRRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)
	mockTeamRepo.On("GetTeamByUser", mock.Anything, authorID).Return(nil, nil)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, domain.ErrNotFound, err)

	mockPRRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestPRMerge_Success(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	prID := "pr-123"
	existingPR := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    "user1",
		ReviewersID: []string{"user2", "user3"},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	mockPRRepo.On("GetPRById", mock.Anything, prID).Return(existingPR, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Merged && pr.MergedAt != nil
	})).Return(nil)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, domain.Merged, pr.Status)
	assert.NotNil(t, pr.MergedAt)

	mockPRRepo.AssertExpectations(t)
}

func TestPRMerge_AlreadyMerged(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	prID := "pr-123"
	mergedTime := time.Now()
	existingPR := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    "user1",
		ReviewersID: []string{"user2", "user3"},
		Status:      domain.Merged,
		MergedAt:    &mergedTime,
	}

	mockPRRepo.On("GetPRById", mock.Anything, prID).Return(existingPR, nil)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, domain.Merged, pr.Status)
	mockPRRepo.AssertNotCalled(t, "SavePR")
}