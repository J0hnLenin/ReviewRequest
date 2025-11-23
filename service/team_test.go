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
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)
	teamName := "testers"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
			{ID: "user2", Name: "User Two", TeamName: teamName, IsActive: true},
		},
	}

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, nil)
	mockRepo.On("SaveTeam", mock.Anything, team).Return(nil)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamSave_TeamExists(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)
	teamName := "existing-team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}

	existingTeam := &domain.Team{Name:teamName}

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(existingTeam, nil)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrTeamExists, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamGetByName_Success(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)
	teamName := "teamName"
	expectedTeam := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(expectedTeam, nil)

	// Act
	team, err := service.TeamGetByName(context.Background(), teamName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTeam, team)
	mockRepo.AssertExpectations(t)
}

func TestTeamGetByName_NotFound(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	mockRepo.On("GetTeamByName", mock.Anything, "non-existent").Return(nil, nil)

	// Act
	team, err := service.TeamGetByName(context.Background(), "non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, domain.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamSave_GetTeamByNameError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)
	teamName := "team team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}
	expectedError := ErrQueryExecution

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, expectedError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertNotCalled(t, "SaveTeam")
	mockRepo.AssertExpectations(t)
}

func TestTeamSave_SaveTeamError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)
	teamName := "test team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
			{ID: "user2", Name: "User Two", TeamName: teamName, IsActive: true},
		},
	}
	expectedError := ErrQueryExecution

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, nil)
	mockRepo.On("SaveTeam", mock.Anything, team).Return(expectedError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamSave_ConnectionError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)
	teamName := "test-team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}
	connectionError := ErrConnection

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, connectionError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, connectionError, err)
	mockRepo.AssertNotCalled(t, "SaveTeam")
	mockRepo.AssertExpectations(t)
}

func TestTeamGetByName_GetTeamByNameError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "test-team"
	expectedError := ErrQueryExecution

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, expectedError)

	// Act
	team, err := service.TeamGetByName(context.Background(), teamName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamGetByName_ConnectionError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "test-team"
	connectionError := ErrConnection

	mockRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, connectionError)

	// Act
	team, err := service.TeamGetByName(context.Background(), teamName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, connectionError, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamSave_EmptyTeamName(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	team := &domain.Team{
		Name:    "",
		Members: []*domain.User{},
	}

	mockRepo.On("GetTeamByName", mock.Anything, "").Return(nil, nil)
	mockRepo.On("SaveTeam", mock.Anything, mock.Anything).Return(nil, nil)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamGetByName_EmptyTeamName(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	mockRepo.On("GetTeamByName", mock.Anything, "").Return(nil, nil)

	// Act
	team, err := service.TeamGetByName(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, domain.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestTeamSave_PartialSaveError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)
	teamName := "large-team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
			{ID: "user2", Name: "User Two", TeamName: teamName, IsActive: true},
			{ID: "user3", Name: "User Three", TeamName: teamName, IsActive: true},
		},
	}
	expectedError := ErrQueryExecution

	mockRepo.On("GetTeamByName", mock.Anything, "large-team").Return(nil, nil)
	mockRepo.On("SaveTeam", mock.Anything, team).Return(expectedError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}