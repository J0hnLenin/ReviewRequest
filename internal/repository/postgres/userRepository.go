package postgres

import (
	"context"
	"database/sql"

	"github.com/J0hnLenin/ReviewRequest/domain"
	"github.com/J0hnLenin/ReviewRequest/service"
)

func (r *PostgresRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, user_name, team_name, is_active 
		FROM users 
		WHERE id = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.TeamName,
		&user.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, service.ErrQueryExecution
	}

	return &user, nil
}

func (r *PostgresRepository) SaveUser(ctx context.Context, u *domain.User) error {
	return r.saveUser(ctx, r.db, u)
}

func (r *PostgresRepository) saveUser(ctx context.Context, execer interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}, u *domain.User) error {
	query := `
		INSERT INTO users (id, user_name, team_name, is_active) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (id) DO UPDATE SET 
			user_name = EXCLUDED.user_name,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active`

	_, err := execer.ExecContext(ctx, query, u.ID, u.Name, u.TeamName, u.IsActive)
	if err != nil {
		return service.ErrQueryExecution
	}

	return nil
}