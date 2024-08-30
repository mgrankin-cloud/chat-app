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
	Chat(ctx context.Context, chatName string, photo []byte) (models.Chat, error)
}

type ChtChanger interface {
	UpdateChat(ctx context.Context, chatID int64, chatName string, photo []byte) error
}

type ChtRemover interface {
	DeleteChat(ctx context.Context, chatId int64) error
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
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

func (c *Chat) Chat(ctx context.Context, chatID int64, chatName string, photo []byte) (models.Chat, error) {
	const op = "Chat.GetChat"

	var chat models.Chat

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chat.ID),
		slog.String("chat_name", chat.Name),
		slog.String("photo", string(chat.Photo)),
	)

	log.Info("attempting to get chat")

	chat, err := c.chtService.Chat(ctx, chatID, chatName, photo)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			c.log.Warn("chat not found", err)
			return chat, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		c.log.Error("failed to get chat", err)
		return chat, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (c *Chat) UpdateChat(ctx context.Context, chatID int64, chatName string, photo []byte) error {
	const op = "Chat.UpdateChat"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
		slog.String("chat_name", chatName),
		slog.String("photo", string(photo)),
	)

	log.Info("attempting to update chat")

	err := c.chtChanger.UpdateChat(ctx, chatID, chatName, photo)
	if err != nil {
		log.Error("Failed to update chat", slog.String("error", err.Error()))
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

	log.Info("attempting to delete chat")

	err := c.chtRemover.DeleteChat(ctx, chatID)
	if err != nil {
		log.Error("Failed to delete chat", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Chat) GetAppSettings(ctx context.Context, appID int) (models.App, error) {
	const op = "User.GetAppSettings"

	log := c.log.With(
		slog.String("op", op),
		slog.Int("app_id", appID),
	)

	log.Info("attempting to get app settings")

	app, err := c.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("Failed to get app settings", slog.String("error", err.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
