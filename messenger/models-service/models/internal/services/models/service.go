package models

import (
	"Messenger-android/messenger/domain/models"
	"Messenger-android/messenger/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	User(ctx context.Context, email, username, phone string) (models.User, error)
}

type ChtService interface {
	Chat(ctx context.Context) (models.Chat, error)
}

type AppService interface {
	App(ctx context.Context) (models.App, error)
}

type MsgService interface {
	Message(ctx context.Context) (models.Message, error)
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

func (m *Models) User(ctx context.Context, email, username, phone string) (models.User, error) {
	const op = "User.GetUser"

	var user models.User

	log := m.log.With(
		slog.String("op", op),
		slog.Int64("user_id", user.ID),
		slog.String("email", user.Email),
		slog.String("username", user.Username),
		slog.String("password", string(user.PassHash)),
		slog.String("photo", string(user.Photo)),
	)

	log.Info("attempting to get user")

	user, err := m.usrService.User(ctx, email, username, phone)
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

func (m *Models) Chat(ctx context.Context, email, username, phone string) (models.Chat, error) {
	const op = "Chat.GetChat"

	var chat models.Chat

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get chat")

	chat, err := m.chtService.Chat(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			m.log.Warn("chat not found", err)
			return chat, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get user", err)
		return chat, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (m *Models) App(ctx context.Context, email, username, phone string) (models.App, error) {
	const op = "App.GetApp"

	var app models.App

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get app")

	app, err := m.appService.App(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			m.log.Warn("chat not found", err)
			return app, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get user", err)
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (m *Models) Message(ctx context.Context, email, username, phone string) (models.Message, error) {
	const op = "Message.GetMessage"

	var msg models.Message

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get message")

	msg, err := m.chtService.Chat(ctx)
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
