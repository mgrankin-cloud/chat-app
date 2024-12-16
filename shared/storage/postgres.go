package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	ssov4 "github.com/mgrankin-cloud/messenger/contract/gen/go/chat"
	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Storage struct {
	db *sql.DB
}

// SendNotification implements notification.NotificationService.
func (s *Storage) SendNotification(ctx context.Context, userID int64, message string) error {
	panic("unimplemented")
}

func (s *Storage) SubscribeNotification(ctx context.Context, userID int64, message string) error {
	panic("unimplemented")
}

var (
	ErrUserExists      = errors.New("user already exists")
	ErrChatExists      = errors.New("chat already exists")
	ErrFileExists      = errors.New("file already exists")
	ErrMessageExists   = errors.New("message already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrChatNotFound    = errors.New("message not found")
	ErrMessageNotFound = errors.New("message not found")
	ErrFileNotFound    = errors.New("file not found")
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

func (s *Storage) SaveUser(ctx context.Context, email string, username string, passHash []byte, phone string) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, username, pass_hash, phone) VALUES($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.QueryRowContext(ctx, email, username, passHash, phone).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email, username, phone string, photo []byte, active bool) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
		SELECT id, email, username, pass_hash, phone 
		FROM users 
		WHERE 
			($1 = '' OR email = $1) AND 
			($2 = '' OR username = $2) AND 
			($3 = '' OR phone = $3) AND
			($4 = '' OR photo = $4) AND
			($5 = '' OR active = $5)
	`

	stmt, err := s.db.Prepare(query)

	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email, username, phone, photo, active)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.PassHash, &user.Phone, &user.Photo, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	const op = "storage.postgres.GetUserByID"

	var user models.User
	err := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by ID: %w", op, ErrUserNotFound)
	}

	return user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, userID int64, email, username, phone string, photo []byte, active bool) (bool, error) {
	const op = "storage.postgres.UpdateUser"

	query := `
		UPDATE users 
		SET email = $1, username = $2, phone = $3, photo = $4, active = $5
		WHERE id = $6
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, email, username, phone, photo, active, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return true, nil
}

func (s *Storage) UpdateUserActiveStatus(ctx context.Context, userID int64, active bool) (bool, error) {
	const op = "storage.postgres.UpdateUserActiveStatus"

	query := `UPDATE users SET active = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, active, userID)
	if err != nil {
		return false, fmt.Errorf("failed to update user active status: %w", op, err)
	}
	return true, nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.DeleteUser"

	query := `
		DELETE FROM users 
		WHERE id = $1
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return true, nil
}

func (s *Storage) ChangePassword(ctx context.Context, userID int64, newPassword string) (bool, error) {
	const op = "storage.postgres.ChangePassword"

	query := `
		UPDATE users 
		SET pass_hash = $1 
		WHERE id = $2
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, newPassword, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return true, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) GetChat(ctx context.Context, chatID int64) (models.Chat, error) {
	const op = "storage.postgres.Chat"

	query := `
		SELECT id, name, chat_type, photo
		FROM chats
		WHERE id = $1
	`

	stmt, err := s.db.Prepare(query)

	var chat models.Chat
	if err != nil {
		return chat, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, chatID)

	err = row.Scan(&chat.ID, &chat.Name, &chat.Photo, &chat.ChatType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return chat, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return chat, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (s *Storage) CreateChat(ctx context.Context, name string, photo []byte, chatType ssov4.ChatType) (int64, error) {
	const op = "storage.postgres.SaveChat"
	var id int64

	stmt, err := s.db.Prepare("INSERT INTO chats(name, photo, chat_type) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRowContext(ctx, name, photo, chatType).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrChatExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateChat(ctx context.Context, name string, photo []byte, chatID int64) (bool, error) {
	const op = "storage.postgres.UpdateChat"

	query := `
		UPDATE chats 
		SET name = $1, photo = $2
		WHERE id = $3
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, name, photo, chatID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s: %w", op, ErrChatNotFound)
	}

	return true, nil
}

func (s *Storage) DeleteChat(ctx context.Context, chatID int64) error {
	const op = "storage.postgres.DeleteChat"

	query := `
		DELETE FROM chats
		WHERE id = $1
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrChatNotFound)
	}

	return nil
}

func (s *Storage) Message(ctx context.Context, msgID int64) (models.Message, error) {
	const op = "storage.postgres.Message"

	query := `
		SELECT id, content, created_at, created_by, reply_to_id, received_by, received_at, status, is_read
		FROM messages
		WHERE id = $1`

	stmt, err := s.db.Prepare(query)

	if err != nil {
		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, msgID)

	var message models.Message
	err = row.Scan(&message.ID, &message.Content, &message.CreatedAt, &message.CreatedBy, &message.ReplyToID, &message.ReceivedBy, &message.ReceivedAt, &message.Status, &message.IsRead)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Message{}, fmt.Errorf("%s: %w", op, ErrMessageNotFound)
		}

		return models.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

func (s *Storage) CreateMessage(ctx context.Context, content string, createdAt *timestamppb.Timestamp, createdBy, replyToID, receivedBy int64, receivedAt *timestamppb.Timestamp, status string, isRead bool) (int64, error) {
	const op = "storage.postgres.CreateMessage"

	stmt, err := s.db.Prepare("INSERT INTO messages(content, created_by, reply_to_id, received_by, created_at, received_at, status, is_read) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.QueryRowContext(ctx, content, createdBy, replyToID, receivedBy, createdAt, receivedAt, status, isRead).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrMessageExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateMessage(ctx context.Context, messageID int64, content string, status string, isRead bool) (bool, error) {
	const op = "storage.postgres.UpdateMessage"

	query := `
		UPDATE messages
		SET content = $1, status = $2, isRead = $3
		WHERE id = $4
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, content, status, isRead, messageID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s: %w", op, ErrMessageNotFound)
	}

	return true, nil
}

func (s *Storage) UpdateMessageStatus(messageID int64, status string) error {
	const op = "storage.postgres.UpdateMessagestatus"

	query := `UPDATE messages SET status = $1 WHERE id = $2`
	_, err := s.db.Exec(query, status, messageID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteMessage(ctx context.Context, messageID int64) (bool, error) {
	const op = "storage.postgres.DeleteMessage"

	query := `
		DELETE FROM messages
		WHERE id = $1
	`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, messageID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s: %w", op, ErrMessageNotFound)
	}

	return true, nil
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

func (s *Storage) GetAppByID(ctx context.Context, id int64) (models.App, error) {
	const op = "storage.postgres.GetAppByID"

	var app models.App
	err := s.db.QueryRowContext(ctx, "SELECT * FROM apps WHERE id = $1", id).Scan(&app)
	if err != nil {
		return models.App{}, fmt.Errorf("failed to get app by ID: %w", op, ErrAppNotFound)
	}

	return app, nil
}

/** Реализация метода UploadMedia     ДЛЯ S3!!!
func (s *Storage) UploadMedia(ctx context.Context, file io.Reader, fileName string) (string, error) {
	result, err := s.s3Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return result.Location, nil
}

// Реализация метода DownloadMedia
func (s *Storage) DownloadMedia(ctx context.Context, fileID string) (io.ReadCloser, error) {
	result, err := s.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return result.Body, nil
}**/