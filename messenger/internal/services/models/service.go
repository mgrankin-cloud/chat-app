package models

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Models struct {
	log        *slog.Logger
	usrService UsrService
	chtService ChtService
	appService AppService
	msgService MsgService
}

type UsrService interface {
	User(ctx context.Context, email, username, phone string, photo []byte) (models.User, error)
}

type ChtService interface {
	Chat(ctx context.Context, name string, photo []byte, chatType int) (models.Chat, error)
}

type AppService interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

type MsgService interface {
	Message(ctx context.Context, content string, createdAt time.Time, createdBy int64, replyToID int64, receivedBy int64, receivedAt time.Time) (models.Message, error)
}

func New(
	log *slog.Logger,
	usrService UsrService,
	chtService ChtService,
	appService AppService,
	msgService MsgService,
) *Models {
	return &Models{
		log:        log,
		usrService: usrService,
		chtService: chtService,
		appService: appService,
		msgService: msgService,
	}
}

func (m *Models) User(ctx context.Context, email, username, phone string, photo []byte) (models.User, error) {
	const op = "User.GetUser"

	var user models.User

	log := m.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("username", username),
		slog.String("phone", phone),
		slog.String("photo", string(photo)),
	)

	log.Info("attempting to get user")

	user, err := m.usrService.User(ctx, email, username, phone, photo)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			m.log.Warn("user not found", err)
			return user, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get user", err)
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (m *Models) Chat(ctx context.Context, name string, photo []byte, chatType int) (models.Chat, error) {
	const op = "Chat.GetChat"

	var chat models.Chat

	log := m.log.With(
		slog.String("op", op),
		slog.String("chat_name", name),
		slog.String("chat_photo", string(photo)),
		slog.String("chat_type", strconv.Itoa(chatType)),
	)

	log.Info("attempting to get message")

	chat, err := m.chtService.Chat(ctx, name, photo, chatType)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			m.log.Warn("message not found", err)
			return chat, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get user", err)
		return chat, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (m *Models) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "App.GetApp"

	var app models.App

	log := m.log.With(
		slog.String("op", op),
		slog.Int64("app_id", appID),
	)

	log.Info("attempting to get app")

	app, err := m.appService.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			m.log.Warn("message not found", err)
			return app, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get user", err)
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (m *Models) Message(ctx context.Context, content string, createdAt time.Time, createdBy int64, replyToID int64, receivedBy int64, receivedAt time.Time) (models.Message, error) {
	const op = "Message.GetMessage"

	var msg models.Message

	log := m.log.With(
		slog.String("op", op),
		slog.String("content", content),
		slog.String("created_at", createdAt.Format(time.RFC3339)),
		slog.Int64("created_by", createdBy),
		slog.Int64("reply_to_id", replyToID),
		slog.Int64("received_by", receivedBy),
		slog.String("received_at", receivedAt.Format(time.RFC3339)),
	)

	log.Info("attempting to get message")

	msg, err := m.msgService.Message(ctx, content, createdAt, createdBy, replyToID, receivedBy, receivedAt)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotFound) {
			m.log.Warn("message not found", err)
			return msg, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get message", err)
		return msg, fmt.Errorf("%s: %w", op, err)
	}

	return msg, nil
}
