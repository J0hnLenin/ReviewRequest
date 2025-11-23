package postgres

import (
	"context"
	"database/sql"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service"
)

func (r *PostgresRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	// 1. Получаем общее количество открытых и закрытых PR
	totalQuery := `
		SELECT 
			COUNT(*) FILTER (WHERE NOT is_merged) as open_prs,
			COUNT(*) FILTER (WHERE is_merged) as closed_prs
		FROM pull_requests`

	var openPRs, closedPRs int
	err := r.db.QueryRowContext(ctx, totalQuery).Scan(&openPRs, &closedPRs)
	if err != nil {
		return nil, service.ErrQueryExecution
	}

	// 2. Получаем ревьюера с наибольшим количеством открытых PR
	topOpenReviewerQuery := `
		SELECT u.id, u.user_name, COUNT(*) as pr_count
		FROM pull_requests pr
		CROSS JOIN UNNEST(pr.reviewers_id) AS reviewer_id
		JOIN users u ON u.id = reviewer_id
		WHERE NOT pr.is_merged
		GROUP BY u.id, u.user_name
		ORDER BY pr_count DESC
		LIMIT 1`

	var topOpenReviewer domain.UserStats
	err = r.db.QueryRowContext(ctx, topOpenReviewerQuery).Scan(
		&topOpenReviewer.UserID,
		&topOpenReviewer.UserName,
		&topOpenReviewer.Count,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, service.ErrQueryExecution
	}

	// 3. Получаем ревьюера с наибольшим количеством закрытых PR
	topClosedReviewerQuery := `
		SELECT u.id, u.user_name, COUNT(*) as pr_count
		FROM pull_requests pr
		CROSS JOIN UNNEST(pr.reviewers_id) AS reviewer_id
		JOIN users u ON u.id = reviewer_id
		WHERE pr.is_merged
		GROUP BY u.id, u.user_name
		ORDER BY pr_count DESC
		LIMIT 1`

	var topClosedReviewer domain.UserStats
	err = r.db.QueryRowContext(ctx, topClosedReviewerQuery).Scan(
		&topClosedReviewer.UserID,
		&topClosedReviewer.UserName,
		&topClosedReviewer.Count,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, service.ErrQueryExecution
	}

	// 4. Получаем автора с наибольшим количеством PR
	topAuthorQuery := `
		SELECT u.id, u.user_name, COUNT(*) as pr_count
		FROM pull_requests pr
		JOIN users u ON u.id = pr.author_id
		GROUP BY u.id, u.user_name
		ORDER BY pr_count DESC
		LIMIT 1`

	var topAuthor domain.UserStats
	err = r.db.QueryRowContext(ctx, topAuthorQuery).Scan(
		&topAuthor.UserID,
		&topAuthor.UserName,
		&topAuthor.Count,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, service.ErrQueryExecution
	}

	stats := &domain.Statistics{
		TotalOpenPRs:   openPRs,
		TotalClosedPRs: closedPRs,
	}

	// Добавляем топовых пользователей только если они есть
	if topOpenReviewer.UserID != "" {
		stats.TopOpenReviewer = &topOpenReviewer
	}
	if topClosedReviewer.UserID != "" {
		stats.TopClosedReviewer = &topClosedReviewer
	}
	if topAuthor.UserID != "" {
		stats.TopAuthor = &topAuthor
	}

	return stats, nil
}