package message

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Message struct {
	log         *slog.Logger
	msgService  MsgService
	msgProvider MsgProvider
}

type MsgService interface {
	Message(ctx context.Context, msgID int64) (models.Message, error)
}

type MsgProvider interface {
	CreateMessage(ctx context.Context, content string, createdAt *timestamppb.Timestamp, createdBy int64, replyToID int64, receivedBy int64, receivedAt *timestamppb.Timestamp, status string, isRead bool) (int64, error)
	UpdateMessage(ctx context.Context, messageID int64, content string, status string, isRead bool) (bool, error)
	DeleteMessage(ctx context.Context, messageID int64) (bool, error)
}


func New(
	log *slog.Logger,
	msgService MsgService,
	provider MsgProvider,
) *Message {
	return &Message{
		log:         log,
		msgService:  msgService,
		msgProvider: provider,
	}
}

func (m *Message) GetMessage(ctx context.Context, msgID int64) (models.Message, error) {
	const op = "Message.GetMessage"

	var message models.Message

	log := m.log.With(
		slog.String("op", op),
		slog.String("content", message.Content),
		slog.String("created_at", message.CreatedAt.AsTime().Format(time.RFC3339)),
		slog.Int64("created_by", message.CreatedBy),
		slog.Int64("reply_to_id", message.ReplyToID),
		slog.Int64("received_by", message.ReceivedBy),
		slog.String("received_at", message.ReceivedAt.AsTime().Format(time.RFC3339)),
	)

	log.Info("attempting to get message")

	message, err := m.msgService.Message(ctx, msgID)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotFound) {
			m.log.Warn("message not found",  slog.String("error message not found", err.Error()))
			return message, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		m.log.Error("failed to get message",  slog.String("error getting message", err.Error()))
		return message, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

func (m *Message) CreateMessage(ctx context.Context, content string, createdAt *timestamppb.Timestamp, createdBy int64, replyToID int64, receivedBy int64, receivedAt *timestamppb.Timestamp, status string, isRead bool) (int64, error) {
	const op = "Message.CreateMessage"

	log := m.log.With(
		slog.String("op", op),
		slog.String("content", content),
		slog.String("createdAt", createdAt.String()),
		slog.Int64("createdBy", createdBy),
		slog.Int64("replyToID", replyToID),
		slog.Int64("receivedBy", receivedBy),
		slog.String("receivedAt", receivedAt.String()),
		slog.String("status", status),
		slog.Bool("is_read", isRead),
	)

	log.Info("attempting to create message")


	messageID, err := m.msgProvider.CreateMessage(ctx, content, createdAt, createdBy, replyToID, receivedBy, receivedAt, status, isRead)
	if err != nil {
		log.Error("failed to create message", slog.String("error creating message", err.Error()))
		return messageID, fmt.Errorf("%s: %w", op, err)
	}

	return messageID, nil
}

func (m *Message) UpdateMessage(ctx context.Context, messageID int64, content string, status string, isRead bool) (bool, error) {
	const op = "Message.UpdateMessage"

	log := m.log.With(
		slog.String("op", op),
		slog.String("content", content),
		slog.String("status", status),
		slog.Bool("is_read", isRead),
	)

	log.Info("attempting to update message")

	success, err := m.msgProvider.UpdateMessage(ctx, messageID, content, status, isRead)
	if err != nil {
		log.Error("Failed to update message", slog.String("error updating message", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return success, nil
}

func (m *Message) DeleteMessage(ctx context.Context, messageID int64) (bool, error) {
	const op = "Message.DeleteMessage"

	log := m.log.With(
		slog.String("op", op),
		slog.Int64("message_id", messageID),
	)

	log.Info("attempting to delete message")

	success, err := m.msgProvider.DeleteMessage(ctx, messageID)
	if err != nil {
		log.Error("Failed to delete message", slog.String("error deleting message", err.Error()))
		return success, fmt.Errorf("%s: %w", op, err)
	}

	return success, nil
}