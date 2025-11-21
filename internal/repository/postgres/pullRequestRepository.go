package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service"
	"github.com/lib/pq"
)

func (r *PostgresRepository) GetPRByAuthor(ctx context.Context, authorID string) ([]*domain.PullRequest, error) {
	query := `
		SELECT id, title, author_id, reviewers_id, is_merged, merged_at 
		FROM pull_requests 
		WHERE author_id = $1`

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, service.ErrQueryExecution
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		pr, err := r.scanPullRequest(rows)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *PostgresRepository) GetPRById(ctx context.Context, id string) (*domain.PullRequest, error) {
	query := `
		SELECT id, title, author_id, reviewers_id, is_merged, merged_at 
		FROM pull_requests 
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	pr, err := r.scanPullRequest(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return pr, err
}

func (r *PostgresRepository) GetPRAndTeam(ctx context.Context, id string) (*domain.PullRequest, *domain.Team, error) {
	query := `
		SELECT 
			pr.id,
			pr.title,
			pr.author_id,
			pr.reviewers_id,
			pr.is_merged,
			pr.merged_at,
			t.team_name,
			COALESCE(array_agg(um.id ORDER BY um.id) FILTER (WHERE um.id IS NOT NULL), '{}') as member_ids,
			COALESCE(array_agg(um.user_name ORDER BY um.id) FILTER (WHERE um.id IS NOT NULL), '{}') as member_names,
			COALESCE(array_agg(um.is_active ORDER BY um.id) FILTER (WHERE um.id IS NOT NULL), '{}') as member_active,
			u.user_name as author_name,
			u.team_name as author_team_name,
			u.is_active as author_active
			
		FROM pull_requests pr
		LEFT JOIN users u ON pr.author_id = u.id
		LEFT JOIN teams t ON u.team_name = t.team_name
		LEFT JOIN users um ON t.team_name = um.team_name
		
		WHERE pr.id = $1
		GROUP BY 
			pr.id, pr.title, pr.author_id, pr.reviewers_id, pr.is_merged, pr.merged_at,
			t.team_name,
			u.user_name, u.team_name, u.is_active`

	var (
		prID, prTitle, authorID string
		reviewers               []string
		isMerged                bool
		mergedAt                *time.Time

		teamName string
		
		memberIDs, memberNames []string
		memberActive           []bool
		
		authorName, authorTeamName string
		authorActive               bool
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&prID,
		&prTitle,
		&authorID,
		pq.Array(&reviewers),
		&isMerged,
		&mergedAt,
		
		&teamName,
		
		pq.Array(&memberIDs),
		pq.Array(&memberNames),
		pq.Array(&memberActive),
		
		&authorName,
		&authorTeamName,
		&authorActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, service.ErrQueryExecution
	}

	pr := &domain.PullRequest{
		ID:          prID,
		Title:       prTitle,
		AuthorID:    authorID,
		ReviewersID: reviewers,
		Status:      domain.PRStatus(isMerged),
		MergedAt:    mergedAt,
	}

	team := &domain.Team{
		Name: teamName,
	}

	team.Members = make([]*domain.User, len(memberIDs))
	for i := range memberIDs {
		team.Members[i] = &domain.User{
			ID:       memberIDs[i],
			Name:     memberNames[i],
			TeamName: teamName,
			IsActive: memberActive[i],
		}
	}

	return pr, team, nil
}

func (r *PostgresRepository) SavePR(ctx context.Context, pr *domain.PullRequest) error {
	query := `
		INSERT INTO pull_requests (id, title, author_id, reviewers_id, is_merged, merged_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET 
			title = EXCLUDED.title,
			author_id = EXCLUDED.author_id,
			reviewers_id = EXCLUDED.reviewers_id,
			is_merged = EXCLUDED.is_merged,
			merged_at = EXCLUDED.merged_at`

	_, err := r.db.ExecContext(ctx, query, 
		pr.ID, 
		pr.Title, 
		pr.AuthorID, 
		pq.Array(pr.ReviewersID), 
		bool(pr.Status),
		pr.MergedAt,
	)
	if err != nil {
		return service.ErrQueryExecution
	}

	return nil
}

func (r *PostgresRepository) scanPullRequest(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	var reviewers []string
	var isMerged bool
	var mergedAt *time.Time

	err := scanner.Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		pq.Array(&reviewers),
		&isMerged,
		&mergedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, service.ErrQueryExecution
	}

	pr.ReviewersID = reviewers
	pr.Status = domain.PRStatus(isMerged)
	pr.MergedAt = mergedAt
	return &pr, nil
}