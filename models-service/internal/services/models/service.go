package models

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	GetUserByID(ctx context.Context, userID int64) (models.User, error)
}

type ChtService interface {
	GetChat(ctx context.Context, chatID int64) (models.Chat, error)
}

type AppService interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

type MsgService interface {
	Message(ctx context.Context, msgID int64) (models.Message, error)
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

func (m *Models) GetUser(ctx context.Context, userID int64) (models.User, error) {
	const op = "User.GetUser"

	var user models.User

	log := m.log.With(
		slog.String("op", op),
		slog.String("email", user.Email),
		slog.String("username", user.Username),
		slog.String("phone", user.Phone),
		slog.Bool("active", user.Active),
	)

	log.Info("attempting to get user")

	user, err := m.usrService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			m.log.Warn("user not found", slog.String("error user not found", err.Error()))
			return user, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get user", slog.String("error getting user", err.Error()))
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (m *Models) GetChat(ctx context.Context, chatID int64) (models.Chat, error) {
	const op = "Chat.GetChat"

	var chat models.Chat

	log := m.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
	)

	log.Info("attempting to get chat")

	chat, err := m.chtService.GetChat(ctx, chatID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			m.log.Warn("chat not found", slog.String("error", err.Error()))
			return chat, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get chat", slog.String("error", err.Error()))
		return chat, fmt.Errorf("%s: %w", op, err)
	}

	log = log.With(
		slog.String("chat_name", chat.Name),
		slog.String("chat_photo", string(chat.Photo)),
	)
	log.Info("chat retrieved successfully")

	return chat, nil
}

func (m *Models) GetApp(ctx context.Context, appID int64) (models.App, error) {
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
			m.log.Warn("app not found", slog.String("error app not found", err.Error()))
			return app, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get app", slog.String("error getting app", err.Error()))
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (m *Models) GetMessage(ctx context.Context, messageID int64) (models.Message, error) {
	const op = "Message.GetMessage"

	var msg models.Message

	log := m.log.With(
		slog.String("op", op),
		slog.String("content", msg.Content),
		slog.String("created_at", msg.CreatedAt.AsTime().Format(time.RFC3339)),
		slog.Int64("created_by", msg.CreatedBy),
		slog.Int64("reply_to_id", msg.ReplyToID),
		slog.Int64("received_by", msg.ReceivedBy),
		slog.String("received_at", msg.ReceivedAt.AsTime().Format(time.RFC3339)),
	)

	log.Info("attempting to get message")

	msg, err := m.msgService.Message(ctx, messageID)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotFound) {
			m.log.Warn("message not found", slog.String("error message not found", err.Error()))
			return msg, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get message", slog.String("error getting message", err.Error()))
		return msg, fmt.Errorf("%s: %w", op, err)
	}

	return msg, nil
}
