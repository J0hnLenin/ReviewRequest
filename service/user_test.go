package service

import (
	"context"
	"testing"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserChangeActive_SuccessActivate(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: false,
	}

	mockRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
	mockRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == userID && user.IsActive == true
	})).Return(nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.IsActive)
	assert.Equal(t, userID, updatedUser.ID)

	mockRepo.AssertExpectations(t)
}

func TestUserChangeActive_SuccessDeactivate(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: true,
	}

	mockRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
	mockRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == userID && user.IsActive == false
	})).Return(nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.False(t, updatedUser.IsActive)
	assert.Equal(t, userID, updatedUser.ID)

	mockRepo.AssertExpectations(t)
}

func TestUserChangeActive_NoChangeSameValue(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: true,
	}

	mockRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.IsActive)
	assert.Equal(t, userID, updatedUser.ID)

	mockRepo.AssertNotCalled(t, "SaveUser")
	mockRepo.AssertExpectations(t)
}

func TestUserChangeActive_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "non-existent-user"

	mockRepo.On("GetUserById", mock.Anything, userID).Return(nil, nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, domain.ErrNotFound, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "SaveUser")
}

func TestUserChangeActive_GetUserError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user123"
	expectedError := ErrQueryExecution

	mockRepo.On("GetUserById", mock.Anything, userID).Return(nil, expectedError)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "SaveUser")
}

func TestUserChangeActive_SaveUserError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: false,
	}
	expectedError := ErrQueryExecution

	mockRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
	mockRepo.On("SaveUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(expectedError)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestUserGetReviews_Success(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "author123"
	expectedPRs := []*domain.PullRequest{
		{
			ID:          "pr1",
			Title:       "First PR",
			AuthorID:    userID,
			ReviewersID: []string{"reviewer1", "reviewer2"},
			Status:      domain.Open,
		},
		{
			ID:          "pr2",
			Title:       "Second PR",
			AuthorID:    userID,
			ReviewersID: []string{"reviewer3"},
			Status:      domain.Merged,
		},
	}

	mockRepo.On("GetPRByAuthor", mock.Anything, userID).Return(expectedPRs, nil)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Len(t, prs, 2)
	assert.Equal(t, expectedPRs, prs)

	mockRepo.AssertExpectations(t)
}

func TestUserGetReviews_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user-with-no-prs"

	mockRepo.On("GetPRByAuthor", mock.Anything, userID).Return([]*domain.PullRequest{}, nil)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Empty(t, prs)

	mockRepo.AssertExpectations(t)
}

func TestUserGetReviews_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user123"
	expectedError := ErrQueryExecution

	mockRepo.On("GetPRByAuthor", mock.Anything, userID).Return(nil, expectedError)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, prs)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestUserGetReviews_NilResult(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	userID := "user123"

	mockRepo.On("GetPRByAuthor", mock.Anything, userID).Return(nil, nil)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, prs)

	mockRepo.AssertExpectations(t)
}