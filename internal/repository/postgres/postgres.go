package postgres

import (
	"database/sql"

	"github.com/J0hnLenin/ReviewRequest/service"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, service.ErrConnection
	}

	if err := db.Ping(); err != nil {
		return nil, service.ErrConnection
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}



