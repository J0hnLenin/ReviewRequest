package service

import (
	"context"
	"testing"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserChangeActive_Success_Activate(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: false,
	}

	mockUserRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
	mockUserRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == userID && user.IsActive == true
	})).Return(nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.IsActive)
	assert.Equal(t, userID, updatedUser.ID)

	mockUserRepo.AssertExpectations(t)
}

func TestUserChangeActive_Success_Deactivate(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: true,
	}

	mockUserRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
	mockUserRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == userID && user.IsActive == false
	})).Return(nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.False(t, updatedUser.IsActive)
	assert.Equal(t, userID, updatedUser.ID)

	mockUserRepo.AssertExpectations(t)
}

func TestUserChangeActive_NoChange_SameValue(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: true,
	}

	mockUserRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.IsActive)
	assert.Equal(t, userID, updatedUser.ID)

	mockUserRepo.AssertNotCalled(t, "SaveUser")
	mockUserRepo.AssertExpectations(t)
}

func TestUserChangeActive_UserNotFound(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "non-existent-user"

	mockUserRepo.On("GetUserById", mock.Anything, userID).Return(nil, nil)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, domain.ErrNotFound, err)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "SaveUser")
}

func TestUserChangeActive_GetUserError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user123"
	expectedError := ErrQueryExecution

	mockUserRepo.On("GetUserById", mock.Anything, userID).Return(nil, expectedError)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, expectedError, err)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "SaveUser")
}

func TestUserChangeActive_SaveUserError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user123"
	currentUser := &domain.User{
		ID:       userID,
		Name:     "Test User",
		TeamName: "team1",
		IsActive: false,
	}
	expectedError := ErrQueryExecution

	mockUserRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
	mockUserRepo.On("SaveUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(expectedError)

	// Act
	updatedUser, err := service.UserChangeActive(context.Background(), userID, true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, expectedError, err)

	mockUserRepo.AssertExpectations(t)
}

func TestUserGetReviews_Success(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

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

	mockPRRepo.On("GetPRByAuthor", mock.Anything, userID).Return(expectedPRs, nil)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Len(t, prs, 2)
	assert.Equal(t, expectedPRs, prs)

	mockPRRepo.AssertExpectations(t)
}

func TestUserGetReviews_EmptyList(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user-with-no-prs"

	mockPRRepo.On("GetPRByAuthor", mock.Anything, userID).Return([]*domain.PullRequest{}, nil)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Empty(t, prs)

	mockPRRepo.AssertExpectations(t)
}

func TestUserGetReviews_RepositoryError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user123"
	expectedError := ErrQueryExecution

	mockPRRepo.On("GetPRByAuthor", mock.Anything, userID).Return(nil, expectedError)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, prs)
	assert.Equal(t, expectedError, err)

	mockPRRepo.AssertExpectations(t)
}

func TestUserGetReviews_NilResult(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	userID := "user123"

	mockPRRepo.On("GetPRByAuthor", mock.Anything, userID).Return(nil, nil)

	// Act
	prs, err := service.UserGetReviews(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, prs)

	mockPRRepo.AssertExpectations(t)
}

// Table-driven tests для UserChangeActive
func TestUserChangeActive_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		initialActive  bool
		newActive      bool
		shouldSave     bool
		expectedError  error
		mockSetup      func(*mocks.MockRepository, string, bool, bool)
	}{
		{
			name:          "activate inactive user",
			initialActive: false,
			newActive:     true,
			shouldSave:    true,
			expectedError: nil,
			mockSetup: func(mockUserRepo *mocks.MockRepository, userID string, initialActive, newActive bool) {
				currentUser := &domain.User{
					ID:       userID,
					Name:     "Test User",
					TeamName: "team1",
					IsActive: initialActive,
				}
				mockUserRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
				mockUserRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.IsActive == newActive
				})).Return(nil)
			},
		},
		{
			name:          "deactivate active user",
			initialActive: true,
			newActive:     false,
			shouldSave:    true,
			expectedError: nil,
			mockSetup: func(mockUserRepo *mocks.MockRepository, userID string, initialActive, newActive bool) {
				currentUser := &domain.User{
					ID:       userID,
					Name:     "Test User",
					TeamName: "team1",
					IsActive: initialActive,
				}
				mockUserRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
				mockUserRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.IsActive == newActive
				})).Return(nil)
			},
		},
		{
			name:          "no change when same value",
			initialActive: true,
			newActive:     true,
			shouldSave:    false,
			expectedError: nil,
			mockSetup: func(mockUserRepo *mocks.MockRepository, userID string, initialActive, newActive bool) {
				currentUser := &domain.User{
					ID:       userID,
					Name:     "Test User",
					TeamName: "team1",
					IsActive: initialActive,
				}
				mockUserRepo.On("GetUserById", mock.Anything, userID).Return(currentUser, nil)
			},
		},
		{
			name:          "user not found",
			initialActive: true,
			newActive:     false,
			shouldSave:    false,
			expectedError: domain.ErrNotFound,
			mockSetup: func(mockUserRepo *mocks.MockRepository, userID string, initialActive, newActive bool) {
				mockUserRepo.On("GetUserById", mock.Anything, userID).Return(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockTeamRepo := &mocks.MockRepository{}
			mockUserRepo := &mocks.MockRepository{}
			mockPRRepo := &mocks.MockRepository{}

			service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)
			userID := "test-user"

			tt.mockSetup(mockUserRepo, userID, tt.initialActive, tt.newActive)

			// Act
			updatedUser, err := service.UserChangeActive(context.Background(), userID, tt.newActive)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, updatedUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updatedUser)
				assert.Equal(t, tt.newActive, updatedUser.IsActive)
			}

			if tt.shouldSave {
				mockUserRepo.AssertCalled(t, "SaveUser", mock.Anything, mock.AnythingOfType("*domain.User"))
			} else {
				mockUserRepo.AssertNotCalled(t, "SaveUser")
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserChangeActive_EdgeCases(t *testing.T) {
	t.Run("empty user ID", func(t *testing.T) {
		// Arrange
		mockTeamRepo := &mocks.MockRepository{}
		mockUserRepo := &mocks.MockRepository{}
		mockPRRepo := &mocks.MockRepository{}

		service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

		mockUserRepo.On("GetUserById", mock.Anything, "").Return(nil, nil)

		// Act
		updatedUser, err := service.UserChangeActive(context.Background(), "", true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.Equal(t, domain.ErrNotFound, err)

		mockUserRepo.AssertExpectations(t)
	})
}