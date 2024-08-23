package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresStorage struct {
	db *sql.DB
}

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Stop() error {
	return s.db.Close()
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, id int) (*User, error) {
	// Implement the method to get user by ID from PostgreSQL
	return nil, nil
}
