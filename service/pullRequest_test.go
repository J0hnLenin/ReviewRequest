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
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

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

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)
	mockRepo.On("GetTeamByUser", mock.Anything, authorID).Return(team, nil)
	mockRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
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

	mockRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestPRCreate_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)
	mockRepo.On("GetTeamByUser", mock.Anything, authorID).Return(nil, nil)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, domain.ErrNotFound, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestPRCreate_PRAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"

	existingPR := &domain.PullRequest {
		ID: prID,
		Title: "Title",
		AuthorID: "user123",
		ReviewersID: make([]string, 0),
		Status: domain.Open,
		MergedAt: nil,
	}

	mockRepo.On("GetPRById", mock.Anything, prID).Return(existingPR, nil)
	mockRepo.On("GetTeamByUser", mock.Anything, authorID).Return(nil, nil)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, domain.ErrPRExists, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetTeamByUser")
}

func TestPRMerge_Success(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	existingPR := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    "user1",
		ReviewersID: []string{"user2", "user3"},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	mockRepo.On("GetPRById", mock.Anything, prID).Return(existingPR, nil)
	mockRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Merged && pr.MergedAt != nil
	})).Return(nil)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, domain.Merged, pr.Status)
	assert.NotNil(t, pr.MergedAt)

	mockRepo.AssertExpectations(t)
}

func TestPRMerge_NotFound(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, domain.ErrNotFound, err)
	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRMerge_AlreadyMerged(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

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

	mockRepo.On("GetPRById", mock.Anything, prID).Return(existingPR, nil)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, domain.Merged, pr.Status)
	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_Success(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "team"

	author := &domain.User{
		ID: "author",
		Name: "Andrey",
		TeamName: teamName,
		IsActive: true,
	}
	oldReviewer := &domain.User{
		ID: "oldReviewer",
		Name: "Ivan",
		TeamName: teamName,
		IsActive: true,
	}
	reassignReviewer := &domain.User{
		ID: "reassignReviewer",
		Name: "Petr",
		TeamName: teamName,
		IsActive: true,
	}
	unactiveMember := &domain.User{
		ID: "unactiveMember",
		Name: "Alice",
		TeamName: teamName,
		IsActive: false,
	}
	activeMember := &domain.User{
		ID: "activeMember",
		Name: "Ilona",
		TeamName: teamName,
		IsActive: true,
	}
	teamMembers := []*domain.User {
		author, oldReviewer, reassignReviewer, unactiveMember, activeMember,
	}
	team := &domain.Team{
		Name: teamName,
		Members: teamMembers,
	}

	prID := "pr-123"
	pr := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    author.ID,
		ReviewersID: []string{oldReviewer.ID, reassignReviewer.ID},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	expectReviewers := []string {
		oldReviewer.ID, activeMember.ID,
	}

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Open})).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Nil(t, pr.MergedAt)
	assert.Equal(t, domain.Open, pr.Status)
	assert.Equal(t, author.ID, pr.AuthorID)
	assert.Equal(t, newReviewerID, activeMember.ID)
	assert.ElementsMatch(t, pr.ReviewersID, expectReviewers)

	mockRepo.AssertExpectations(t)
}

func TestPRreassign_ButPRAlreadyMerged(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "team"

	author := &domain.User{
		ID: "author",
		Name: "Andrey",
		TeamName: teamName,
		IsActive: true,
	}
	oldReviewer := &domain.User{
		ID: "oldReviewer",
		Name: "Ivan",
		TeamName: teamName,
		IsActive: true,
	}
	reassignReviewer := &domain.User{
		ID: "reassignReviewer",
		Name: "Petr",
		TeamName: teamName,
		IsActive: true,
	}
	unactiveMember := &domain.User{
		ID: "unactiveMember",
		Name: "Alice",
		TeamName: teamName,
		IsActive: false,
	}
	activeMember := &domain.User{
		ID: "activeMember",
		Name: "Ilona",
		TeamName: teamName,
		IsActive: true,
	}
	teamMembers := []*domain.User {
		author, oldReviewer, reassignReviewer, unactiveMember, activeMember,
	}
	team := &domain.Team{
		Name: teamName,
		Members: teamMembers,
	}

	prID := "pr-123"
	pr := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    author.ID,
		ReviewersID: []string{oldReviewer.ID, reassignReviewer.ID},
		Status:      domain.Merged,
		MergedAt:    nil,
	}

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Merged})).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrPRMerged, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButNoCandidate(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "team"

	author := &domain.User{
		ID: "author",
		Name: "Andrey",
		TeamName: teamName,
		IsActive: true,
	}
	oldReviewer := &domain.User{
		ID: "oldReviewer",
		Name: "Ivan",
		TeamName: teamName,
		IsActive: true,
	}
	reassignReviewer := &domain.User{
		ID: "reassignReviewer",
		Name: "Petr",
		TeamName: teamName,
		IsActive: true,
	}
	unactiveMember := &domain.User{
		ID: "activeMember",
		Name: "Alice",
		TeamName: teamName,
		IsActive: false,
	}
	teamMembers := []*domain.User {
		author, oldReviewer, reassignReviewer, unactiveMember,
	}
	team := &domain.Team{
		Name: teamName,
		Members: teamMembers,
	}

	prID := "pr-123"
	pr := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    author.ID,
		ReviewersID: []string{oldReviewer.ID, reassignReviewer.ID},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Open})).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNoCandidate, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButReviewerNotAssigned(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "team"

	author := &domain.User{
		ID: "author",
		Name: "Andrey",
		TeamName: teamName,
		IsActive: true,
	}
	oldReviewer := &domain.User{
		ID: "oldReviewer",
		Name: "Ivan",
		TeamName: teamName,
		IsActive: true,
	}
	reassignReviewer := &domain.User{
		ID: "reassignReviewer",
		Name: "Petr",
		TeamName: teamName,
		IsActive: true,
	}
	unactiveMember := &domain.User{
		ID: "unactiveMember",
		Name: "Alice",
		TeamName: teamName,
		IsActive: false,
	}
	activeMember := &domain.User{
		ID: "activeMember",
		Name: "Ilona",
		TeamName: teamName,
		IsActive: true,
	}
	teamMembers := []*domain.User {
		author, oldReviewer, reassignReviewer, unactiveMember, activeMember,
	}
	team := &domain.Team{
		Name: teamName,
		Members: teamMembers,
	}

	prID := "pr-123"
	pr := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    author.ID,
		ReviewersID: []string{oldReviewer.ID},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Open})).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNotAssigned, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButPRNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "team"

	reassignReviewer := &domain.User{
		ID: "reassignReviewer",
		Name: "Petr",
		TeamName: teamName,
		IsActive: true,
	}


	prID := "pr-123"

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(nil, nil, nil)
	mockRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockRepo.On("SavePR", mock.Anything, mock.Anything).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNotFound, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButReviewerNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	teamName := "team"

	author := &domain.User{
		ID: "author",
		Name: "Andrey",
		TeamName: teamName,
		IsActive: true,
	}
	oldReviewer := &domain.User{
		ID: "oldReviewer",
		Name: "Ivan",
		TeamName: teamName,
		IsActive: true,
	}
	unactiveMember := &domain.User{
		ID: "unactiveMember",
		Name: "Alice",
		TeamName: teamName,
		IsActive: false,
	}
	activeMember := &domain.User{
		ID: "activeMember",
		Name: "Ilona",
		TeamName: teamName,
		IsActive: true,
	}
	teamMembers := []*domain.User {
		author, oldReviewer, unactiveMember, activeMember,
	}
	team := &domain.Team{
		Name: teamName,
		Members: teamMembers,
	}

	prID := "pr-123"
	pr := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    author.ID,
		ReviewersID: []string{oldReviewer.ID},
		Status:      domain.Open,
		MergedAt:    nil,
	}
	
	reassignReviewerID := "r-321"
	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockRepo.On("GetUserById", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("SavePR", mock.Anything, mock.Anything).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNotFound, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

// Новые тесты для покрытия ошибок базы данных

func TestPRCreate_GetPRByIdError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"
	expectedError := ErrQueryExecution

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, expectedError)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertNotCalled(t, "GetTeamByUser")
	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRCreate_GetTeamByUserError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"
	expectedError := ErrQueryExecution

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)
	mockRepo.On("GetTeamByUser", mock.Anything, authorID).Return(nil, expectedError)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRCreate_SavePRError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"
	expectedError := ErrQueryExecution

	team := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "user1", Name: "Author", TeamName: "test-team", IsActive: true},
			{ID: "user2", Name: "Reviewer 1", TeamName: "test-team", IsActive: true},
		},
	}

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)
	mockRepo.On("GetTeamByUser", mock.Anything, authorID).Return(team, nil)
	mockRepo.On("SavePR", mock.Anything, mock.AnythingOfType("*domain.PullRequest")).Return(expectedError)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestPRMerge_GetPRByIdError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}
	service := NewService(mockRepo)

	prID := "pr-123"
	expectedError := ErrQueryExecution

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, expectedError)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRMerge_SavePRError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	expectedError := ErrQueryExecution

	existingPR := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    "user1",
		ReviewersID: []string{"user2", "user3"},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	mockRepo.On("GetPRById", mock.Anything, prID).Return(existingPR, nil)
	mockRepo.On("SavePR", mock.Anything, mock.AnythingOfType("*domain.PullRequest")).Return(expectedError)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestPRreassign_GetPRAndTeamError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	reviewerID := "reviewer1"
	expectedError := ErrQueryExecution

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(nil, nil, expectedError)

	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reviewerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertNotCalled(t, "GetUserById")
	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_GetUserByIdError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	reviewerID := "reviewer1"
	expectedError := ErrQueryExecution

	team := &domain.Team{
		Name: "test-team",
		Members: []*domain.User{
			{ID: "user1", Name: "Author", TeamName: "test-team", IsActive: true},
			{ID: "reviewer1", Name: "Reviewer", TeamName: "test-team", IsActive: true},
		},
	}

	pr := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    "user1",
		ReviewersID: []string{reviewerID},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockRepo.On("GetUserById", mock.Anything, reviewerID).Return(nil, expectedError)

	// Act
	resultPR, newReviewerID, err := service.PRreassign(context.Background(), prID, reviewerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultPR)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_SavePRError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	teamName := "test-team"
	oldReviewer := &domain.User{
		ID:       "reviewer1",
		Name:     "Reviewer 1",
		TeamName: teamName,
		IsActive: true,
	}
	
	expectedError := ErrQueryExecution

	team := &domain.Team{
		Name: teamName,
		Members: []*domain.User{
			{ID: "user1", Name: "Author", TeamName: teamName, IsActive: true},
			oldReviewer,
			{ID: "reviewer2", Name: "Reviewer 2", TeamName: teamName, IsActive: true},
		},
	}

	pr := &domain.PullRequest{
		ID:          prID,
		Title:       "Test PR",
		AuthorID:    "user1",
		ReviewersID: []string{oldReviewer.ID},
		Status:      domain.Open,
		MergedAt:    nil,
	}

	

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockRepo.On("GetUserById", mock.Anything, oldReviewer.ID).Return(oldReviewer, nil)
	mockRepo.On("SavePR", mock.Anything, mock.AnythingOfType("*domain.PullRequest")).Return(expectedError)

	// Act
	resultPR, newReviewerID, err := service.PRreassign(context.Background(), prID, oldReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultPR)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestPRCreate_ConnectionError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	title := "Test PR"
	authorID := "user1"
	connectionError := ErrConnection

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, connectionError)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, connectionError, err)
}

func TestPRMerge_ConnectionError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	connectionError := ErrConnection

	mockRepo.On("GetPRById", mock.Anything, prID).Return(nil, connectionError)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, connectionError, err)
}

func TestPRreassign_ConnectionError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.MockRepository{}

	service := NewService(mockRepo)

	prID := "pr-123"
	reviewerID := "reviewer1"
	connectionError := ErrConnection

	mockRepo.On("GetPRAndTeam", mock.Anything, prID).Return(nil, nil, connectionError)

	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reviewerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, connectionError, err)
}