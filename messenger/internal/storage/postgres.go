package storage

import (
	"Messenger-android/messenger/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	db *sql.DB
}

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrChatNotFound    = errors.New("message not found")
	ErrMessageNotFound = errors.New("message not found")
	ErrMediaNotFound   = errors.New("media not found")
	ErrAppNotFound     = errors.New("app not found")
)

// Конструктор Storage для PostgreSQL
func New(connString string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, username string, passHash []byte, phone string, photo []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, username, pass_hash, phone, photo) VALUES($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.QueryRowContext(ctx, email, username, passHash, phone, photo).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateUser(ctx context.Context, userID int64, email, username, phone string, photo []byte) error {
	const op = "storage.postgres.UpdateUser"

	query := `
		UPDATE users 
		SET email = $1, username = $2, phone = $3, photo = $4 
		WHERE id = $5
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, email, username, phone, photo, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return nil
}

func (s *Storage) User(ctx context.Context, email, username, phone string, photo []byte) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
		SELECT id, email, username, pass_hash, phone 
		FROM users 
		WHERE 
			($1 = '' OR email = $1) AND 
			($2 = '' OR username = $2) AND 
			($3 = '' OR phone = $3) AND
			($4 = '' OR photo = $4) 
	`

	stmt, err := s.db.Prepare(query)

	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email, username, phone, photo)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.PassHash, &user.Phone, &user.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "storage.postgres.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	const op = "storage.postgres.GetUserByID"

	var user models.User
	err := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by ID: %w", ErrUserNotFound)
	}

	return user, nil
}

func (s *Storage) ChangePassword(ctx context.Context, userID int64, newPassword string) error {
	const op = "storage.postgres.ChangePassword"

	query := `
		UPDATE users 
		SET pass_hash = $1 
		WHERE id = $2
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, newPassword, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int64) error {
	const op = "storage.postgres.DeleteUser"

	query := `
		DELETE FROM users 
		WHERE id = $1
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return nil
}
