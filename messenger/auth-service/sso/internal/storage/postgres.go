package storage

import (
	"Messenger-android/messenger/auth-service/sso/internal/domain/models"
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
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
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

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, username string, passHash []byte, phone string) (int64, error) {
	const op = "storage.postgres.SaveUser"

	// Простенький запрос на добавление пользователя
	stmt, err := s.db.Prepare("INSERT INTO users(email, username, pass_hash, phone) VALUES($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем запрос, передав параметры
	var id int64
	err = stmt.QueryRowContext(ctx, email, username, passHash, phone).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		// Обработка ошибки уникального ограничения (constraint violation)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user	.
func (s *Storage) User(ctx context.Context, email, username, phone string) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
		SELECT id, email, username, pass_hash, phone 
		FROM users 
		WHERE 
			($1 = '' OR email = $1) AND 
			($2 = '' OR username = $2) AND 
			($3 = '' OR phone = $3)
	`

	stmt, err := s.db.Prepare(query)

	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email, username, phone)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.PassHash, &user.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// App returns app by id.
func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

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

/*func (s *Storage) GetUserByID(ctx context.Context, id int) (*User, error) {
	// Implement the method to get user by ID from PostgreSQL
	// Example:
	// var user User
	// err := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get user by ID: %w", err)
	// }
	// return &user, nil
	return nil, nil
}*/
