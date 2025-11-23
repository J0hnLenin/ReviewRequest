package main

import (
	"log"
	"net/http"
	"os"

	"github.com/J0hnLenin/ReviewRequest/internal/api/handler"
	"github.com/J0hnLenin/ReviewRequest/internal/repository/postgres"
	"github.com/J0hnLenin/ReviewRequest/service"
)

func main() {
    connStr := os.Getenv("DATABASE_URL");
    if connStr == "" {
        panic("Connection string empty")
    }

    repo, err := postgres.NewPostgresRepository(connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer repo.Close()

    svc := service.NewService(repo)

    h := handler.NewHandler(svc)

    http.HandleFunc("/team/add", h.TeamAdd)
    http.HandleFunc("/team/get", h.TeamGet)
    http.HandleFunc("/users/setIsActive", h.UserSetIsActive)
    http.HandleFunc("/pullRequest/create", h.PRCreate)
    http.HandleFunc("/pullRequest/merge", h.PRMerge)
    http.HandleFunc("/pullRequest/reassign", h.PRReassign)
    http.HandleFunc("/users/getReview", h.UserGetReviews)
    http.HandleFunc("/health", h.HealthCheck)
    http.HandleFunc("/statistics", h.GetStatistics)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}