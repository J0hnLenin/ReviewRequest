package postgres

import (
	"context"
	"database/sql"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service"
	"github.com/lib/pq"
)

func (r *PostgresRepository) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	query := `
		SELECT t.team_name, 
		       COALESCE(array_agg(u.id) FILTER (WHERE u.id IS NOT NULL), '{}') as member_ids,
		       COALESCE(array_agg(u.user_name) FILTER (WHERE u.id IS NOT NULL), '{}') as member_names,
		       COALESCE(array_agg(u.is_active) FILTER (WHERE u.id IS NOT NULL), '{}') as member_active
		FROM teams t
		LEFT JOIN users u ON t.team_name = u.team_name
		WHERE t.team_name = $1
		GROUP BY t.team_name`

	var team domain.Team
	var memberIDs, memberNames []string
	var memberActive []bool

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&team.Name,
		pq.Array(&memberIDs),
		pq.Array(&memberNames),
		pq.Array(&memberActive),
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, service.ErrQueryExecution
	}

	team.Members = make([]*domain.User, len(memberIDs))
	for i := range memberIDs {
		team.Members[i] = &domain.User{
			ID:       memberIDs[i],
			Name:     memberNames[i],
			TeamName: team.Name,
			IsActive: memberActive[i],
		}
	}

	return &team, nil
}

func (r *PostgresRepository) SaveTeam(ctx context.Context, t *domain.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return service.ErrQueryExecution
	}
	defer tx.Rollback()

	query := `INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING`
	_, err = tx.ExecContext(ctx, query, t.Name)
	if err != nil {
		return service.ErrQueryExecution
	}

	for _, user := range t.Members {
		if err := r.saveUser(ctx, tx, user); err != nil {
			return service.ErrQueryExecution
		}
	}

	return tx.Commit()
}