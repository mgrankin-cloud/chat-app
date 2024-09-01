package user

import (
	"Messenger-android/messenger/internal/domain/models"
	"Messenger-android/messenger/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Chat struct {
	log         *slog.Logger
	chtService  ChtService
	chtChanger  ChtChanger
	chtRemover  ChtRemover
	appProvider AppProvider
}

type ChtService interface {
	Chat(ctx context.Context, chatName string, photo []byte, chatType int) (models.Chat, error)
}

type ChtChanger interface {
	UpdateChat(ctx context.Context, chatName string, photo []byte) error
}

type ChtRemover interface {
	DeleteChat(ctx context.Context, chatId int64) error
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(
	log *slog.Logger,
	chtService ChtService,
	chtChanger ChtChanger,
	chtRemover ChtRemover,
	appProvider AppProvider,
) *Chat {
	return &Chat{
		log:         log,
		chtService:  chtService,
		chtChanger:  chtChanger,
		chtRemover:  chtRemover,
		appProvider: appProvider,
	}
}

func (c *Chat) Chat(ctx context.Context, chatName string, photo []byte, chatType int) (models.Chat, error) {
	const op = "Chat.GetChat"

	var chat models.Chat

	log := c.log.With(
		slog.String("op", op),
		slog.String("chat_name", chatName),
		slog.String("photo", string(photo)),
		slog.Int("chat_type", chatType),
	)

	log.Info("attempting to get message")

	chat, err := c.chtService.Chat(ctx, chatName, photo, chatType)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			c.log.Warn("message not found", err)
			return chat, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		c.log.Error("failed to get message", err)
		return chat, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (c *Chat) UpdateChat(ctx context.Context, chatName string, photo []byte) error {
	const op = "Chat.UpdateChat"

	log := c.log.With(
		slog.String("op", op),
		slog.String("chat_name", chatName),
		slog.String("photo", string(photo)),
	)

	log.Info("attempting to update message")

	err := c.chtChanger.UpdateChat(ctx, chatName, photo)
	if err != nil {
		log.Error("Failed to update message", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Chat) DeleteChat(ctx context.Context, chatID int64) error {
	const op = "Chat.DeleteChat"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
	)

	log.Info("attempting to delete message")

	err := c.chtRemover.DeleteChat(ctx, chatID)
	if err != nil {
		log.Error("Failed to delete message", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Chat) GetAppSettings(ctx context.Context, appID int64) (models.App, error) {
	const op = "User.GetAppSettings"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("app_id", appID),
	)

	log.Info("attempting to get app settings")

	app, err := c.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("Failed to get app settings", slog.String("error", err.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
