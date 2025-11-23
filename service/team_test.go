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
	teamName := "testers"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
			{ID: "user2", Name: "User Two", TeamName: teamName, IsActive: true},
		},
	}

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, nil)
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
	teamName := "existing-team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}

	existingTeam := &domain.Team{Name:teamName}

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(existingTeam, nil)

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
	teamName := "teamName"
	expectedTeam := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(expectedTeam, nil)

	// Act
	team, err := service.TeamGetByName(context.Background(), teamName)

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

func TestTeamSave_GetTeamByNameError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)
	teamName := "team team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}
	expectedError := ErrQueryExecution

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, expectedError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockTeamRepo.AssertNotCalled(t, "SaveTeam")
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamSave_SaveTeamError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)
	teamName := "test team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
			{ID: "user2", Name: "User Two", TeamName: teamName, IsActive: true},
		},
	}
	expectedError := ErrQueryExecution

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, nil)
	mockTeamRepo.On("SaveTeam", mock.Anything, team).Return(expectedError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamSave_ConnectionError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)
	teamName := "test-team"
	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "User One", TeamName: teamName, IsActive: true},
		},
	}
	connectionError := ErrConnection

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, connectionError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, connectionError, err)
	mockTeamRepo.AssertNotCalled(t, "SaveTeam")
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamGetByName_GetTeamByNameError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	teamName := "test-team"
	expectedError := ErrQueryExecution

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, expectedError)

	// Act
	team, err := service.TeamGetByName(context.Background(), teamName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, expectedError, err)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamGetByName_ConnectionError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	teamName := "test-team"
	connectionError := ErrConnection

	mockTeamRepo.On("GetTeamByName", mock.Anything, teamName).Return(nil, connectionError)

	// Act
	team, err := service.TeamGetByName(context.Background(), teamName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, connectionError, err)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamSave_EmptyTeamName(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	team := &domain.Team{
		Name:    "",
		Members: []*domain.User{},
	}

	mockTeamRepo.On("GetTeamByName", mock.Anything, "").Return(nil, nil)
	mockTeamRepo.On("SaveTeam", mock.Anything, mock.Anything).Return(nil, nil)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.NoError(t, err)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamGetByName_EmptyTeamName(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	mockTeamRepo.On("GetTeamByName", mock.Anything, "").Return(nil, nil)

	// Act
	team, err := service.TeamGetByName(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, domain.ErrNotFound, err)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamSave_PartialSaveError(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)
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

	mockTeamRepo.On("GetTeamByName", mock.Anything, "large-team").Return(nil, nil)
	mockTeamRepo.On("SaveTeam", mock.Anything, team).Return(expectedError)

	// Act
	err := service.TeamSave(context.Background(), team)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockTeamRepo.AssertExpectations(t)
}