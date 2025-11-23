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

func TestPRCreate_UserNotFound(t *testing.T) {
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

func TestPRCreate_PRAlreadyExists(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

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

	mockPRRepo.On("GetPRById", mock.Anything, prID).Return(existingPR, nil)
	mockTeamRepo.On("GetTeamByUser", mock.Anything, authorID).Return(nil, nil)

	// Act
	pr, err := service.PRCreate(context.Background(), prID, title, authorID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, domain.ErrPRExists, err)

	mockPRRepo.AssertExpectations(t)
	mockTeamRepo.AssertNotCalled(t, "GetTeamByUser")
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

func TestPRMerge_NotFound(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	prID := "pr-123"
	mockPRRepo.On("GetPRById", mock.Anything, prID).Return(nil, nil)

	// Act
	pr, err := service.PRMerge(context.Background(), prID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, domain.ErrNotFound, err)
	mockPRRepo.AssertNotCalled(t, "SavePR")
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

func TestPRreassign_Success(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

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

	mockPRRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockUserRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
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

	mockPRRepo.AssertExpectations(t)
}

func TestPRreassign_ButPRAlreadyMerged(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

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

	mockPRRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockUserRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Merged})).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrPRMerged, err)

	mockPRRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButNoCandidate(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

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

	mockPRRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockUserRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Open})).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNoCandidate, err)

	mockPRRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButReviewerNotAssigned(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

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

	mockPRRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockUserRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.MatchedBy(func(pr *domain.PullRequest) bool {
		return pr.Status == domain.Open})).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNotAssigned, err)

	mockPRRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButPRNotFound(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

	teamName := "team"

	reassignReviewer := &domain.User{
		ID: "reassignReviewer",
		Name: "Petr",
		TeamName: teamName,
		IsActive: true,
	}


	prID := "pr-123"

	mockPRRepo.On("GetPRAndTeam", mock.Anything, prID).Return(nil, nil, nil)
	mockUserRepo.On("GetUserById", mock.Anything, reassignReviewer.ID).Return(reassignReviewer, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.Anything).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewer.ID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNotFound, err)

	mockPRRepo.AssertNotCalled(t, "SavePR")
}

func TestPRreassign_ButReviewerNotFound(t *testing.T) {
	// Arrange
	mockTeamRepo := &mocks.MockRepository{}
	mockUserRepo := &mocks.MockRepository{}
	mockPRRepo := &mocks.MockRepository{}

	service := NewService(mockTeamRepo, mockUserRepo, mockPRRepo)

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
	mockPRRepo.On("GetPRAndTeam", mock.Anything, prID).Return(pr, team, nil)
	mockUserRepo.On("GetUserById", mock.Anything, mock.Anything).Return(nil, nil)
	mockPRRepo.On("SavePR", mock.Anything, mock.Anything).Return(nil)
	
	// Act
	pr, newReviewerID, err := service.PRreassign(context.Background(), prID, reassignReviewerID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.Equal(t, "", newReviewerID)
	assert.Equal(t, domain.ErrNotFound, err)

	mockPRRepo.AssertNotCalled(t, "SavePR")
}