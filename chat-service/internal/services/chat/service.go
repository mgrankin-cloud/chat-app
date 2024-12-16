package chat

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	ssov4 "github.com/mgrankin-cloud/messenger/contract/gen/go/chat"
	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Chat struct {
	log         *slog.Logger
	chtService  ChtService
	chtProvider ChtProvider
	appProvider AppProvider
}

type ChtService interface {
	GetChat(ctx context.Context, chatID int64) (models.Chat, error)
}

type ChtProvider interface {
	CreateChat(ctx context.Context, chatName string, photo []byte, chatType ssov4.ChatType) (int64, error)
	UpdateChat(ctx context.Context, chatName string, photo []byte, chatID int64) (success bool, err error)
	DeleteChat(ctx context.Context, chatId int64) error
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(
	log *slog.Logger,
	chtService ChtService,
	cthProvider ChtProvider,
	appProvider AppProvider,
) *Chat {
	return &Chat{
		log:         log,
		chtService:  chtService,
		chtProvider: cthProvider,
		appProvider: appProvider,
	}
}

func (c *Chat) GetChat(ctx context.Context, chatID int64) (models.Chat, error) {
	const op = "Chat.GetChat"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
	)

	log.Info("attempting to get chat")

	var chat models.Chat
	chat, err := c.chtService.GetChat(ctx, chatID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			c.log.Warn("chat not found", err)
			return models.Chat{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		c.log.Error("failed to get chat", err)
		return models.Chat{} ,fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (c *Chat) CreateChat(ctx context.Context, chatName string, photo []byte, chatType ssov4.ChatType) (int64, error) {
	const op = "Chat.CreateChat"

	log := c.log.With(
		slog.String("op", op),
		slog.String("chat_name", chatName),
		slog.String("photo", string(photo)),
		slog.String("chat_type", string(chatType)),
	)

	log.Info("attempting to create chat")

	chatID, err := c.chtProvider.CreateChat(ctx, chatName, photo, chatType)
	if err != nil {
		log.Error("failed to create chat", slog.String("error", err.Error()))
		return chatID, fmt.Errorf("%s: %w", op, err)
	}

	return chatID, nil
}

func (c *Chat) UpdateChat(ctx context.Context, chatName string, photo []byte, chatID int64) (success bool, err error) {
	const op = "Chat.UpdateChat"

	log := c.log.With(
		slog.String("op", op),
		slog.String("chat_name", chatName),
		slog.String("photo", string(photo)),
	)

	log.Info("attempting to update chat")

	success, err = c.chtProvider.UpdateChat(ctx, chatName, photo, chatID)
	if err != nil {
		log.Error("failed to update chat", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return success, nil
}

func (c *Chat) DeleteChat(ctx context.Context, chatID int64) error {
	const op = "Chat.DeleteChat"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("chat_id", chatID),
	)

	log.Info("attempting to delete chat")

	err := c.chtProvider.DeleteChat(ctx, chatID)
	if err != nil {
		log.Error("failed to delete chat", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Chat) GetAppSettings(ctx context.Context, appID int64) (models.App, error) {
	const op = "Chat.GetAppSettings"

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
