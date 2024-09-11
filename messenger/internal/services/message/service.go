package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Message struct {
	log         *slog.Logger
	msgService  MsgService
	msgChanger  MsgChanger
	msgProvider MsgProvider
	msgRemover  MsgRemover
	appProvider AppProvider
}

type MsgService interface {
	Message(ctx context.Context, content string, createdAt time.Time, createdBy int64, replyToID int64, receivedBy int64, receivedAt time.Time) (models.Message, error)
}

type MsgChanger interface {
	UpdateMessage(ctx context.Context, content string) error
}

type MsgProvider interface {
	CreateMessage(ctx context.Context, content string, createdAt time.Time, createdBy int64, replyToID int64, receivedBy int64, receivedAt time.Time) (models.Message, error)
}

type MsgRemover interface {
	DeleteMessage(ctx context.Context, messageID int64) error
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(
	log *slog.Logger,
	msgService MsgService,
	msgChanger MsgChanger,
	provider MsgProvider,
	msgRemover MsgRemover,
	appProvider AppProvider,
) *Message {
	return &Message{
		log:         log,
		msgService:  msgService,
		msgChanger:  msgChanger,
		msgProvider: provider,
		msgRemover:  msgRemover,
		appProvider: appProvider,
	}
}

func (m *Message) Message(ctx context.Context, content string, createdAt time.Time, createdBy int64, replyToID int64, receivedBy int64, receivedAt time.Time) (models.Message, error) {
	const op = "Message.GetMessage"

	var message models.Message

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get message")

	message, err := m.msgService.Message(ctx, content, createdAt, createdBy, replyToID, receivedBy, receivedAt)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotFound) {
			m.log.Warn("message not found", err)
			return message, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get message", err)
		return message, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

func (m *Message) CreateMessage(ctx context.Context, content string, createdAt time.Time, createdBy int64, replyToID int64, receivedBy int64, receivedAt time.Time) (models.Message, error) {
	const op = "Message.CreateMessage"

	log := m.log.With(
		slog.String("op", op),
		slog.String("content", content),
		slog.String("createdAt", createdAt.String()),
		slog.Int64("createdBy", createdBy),
		slog.Int64("replyToID", replyToID),
		slog.Int64("receivedBy", receivedBy),
		slog.Int64("receivedAt", receivedAt.Unix()),
	)

	log.Info("attempting to create message")

	var message models.Message

	message, err := m.msgProvider.CreateMessage(ctx, content, createdAt, createdBy, replyToID, receivedBy, receivedAt)
	if err != nil {
		log.Error("failed to create message", err)
		return message, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

func (m *Message) UpdateMessage(ctx context.Context, content string) error {
	const op = "Message.UpdateMessage"

	log := m.log.With(
		slog.String("op", op),
		slog.String("content", content),
	)

	log.Info("attempting to update message")

	err := m.msgChanger.UpdateMessage(ctx, content)
	if err != nil {
		log.Error("Failed to update message", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *Message) DeleteMessage(ctx context.Context, messageID int64) error {
	const op = "Message.DeleteMessage"

	log := m.log.With(
		slog.String("op", op),
		slog.Int64("message_id", messageID),
	)

	log.Info("attempting to delete message")

	err := m.msgRemover.DeleteMessage(ctx, messageID)
	if err != nil {
		log.Error("Failed to delete message", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *Message) GetAppSettings(ctx context.Context, appID int64) (models.App, error) {
	const op = "User.GetAppSettings"

	log := m.log.With(
		slog.String("op", op),
		slog.Int64("app_id", appID),
	)

	log.Info("attempting to get app settings")

	app, err := m.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("Failed to get app settings", slog.String("error", err.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
