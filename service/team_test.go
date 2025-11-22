package service

import (
	"context"
	"testing"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTeamSave_Success(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	team := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: "test-team", IsActive: true},
			{ID: "user2", Name: "User Two", TeamName: "test-team", IsActive: true},
		},
	}

	mockTeamRepo.On("GetTeamByName", mock.Anything, "test-team").Return(nil, nil)
	mockTeamRepo.On("SaveTeam", mock.Anything, team).Return(nil)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.NoError(t, err)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamSave_TeamExists(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	team := &domain.Team{
		Name: "existing-team",
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: "existing-team", IsActive: true},
		},
	}

	existingTeam := &domain.Team{Name: "existing-team"}

	mockTeamRepo.On("GetTeamByName", mock.Anything, "existing-team").Return(existingTeam, nil)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrTeamExists, err)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamGetByName_Success(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	expectedTeam := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: "test-team", IsActive: true},
		},
	}

	mockTeamRepo.On("GetTeamByName", mock.Anything, "test-team").Return(expectedTeam, nil)

	// Act
	team, err := service.TeamGetByName(context.Background(), "test-team")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTeam, team)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamGetByName_NotFound(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	mockTeamRepo.On("GetTeamByName", mock.Anything, "non-existent").Return(nil, nil)

	// Act
	team, err := service.TeamGetByName(context.Background(), "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, domain.ErrNotFound, err)
	mockTeamRepo.AssertExpectations(t)
}