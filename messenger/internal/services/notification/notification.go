package media

import (
	"fmt"

	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/internal/logger/s1"

	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var redisClient *redis.Client

type Notification struct {
	log           *slog.Logger
	ntfSender     NotificationSender
	ntfSubscriber NotificationSubscriber
	appProvider   AppProvider
	tokenTTL      time.Duration
}

type NotificationSender interface {
	SendNotification(ctx context.Context, userID int64, message string) error
}

type NotificationSubscriber interface {
	SubscribeNotification(ctx context.Context, userID int64, message string) error
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func New(
	log *slog.Logger,
	ntfSender NotificationSender,
	ntfSubscriber NotificationSubscriber,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Notification {
	return &Notification{
		ntfSender:     ntfSender,
		ntfSubscriber: ntfSubscriber,
		log:           log,
		appProvider:   appProvider,
		tokenTTL:      tokenTTL,
	}
}

func (n *Notification) SendNotification(ctx context.Context, userID int64, message string) error {
	const op = "Notification.SendNotification"

	log := n.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
		slog.String("message", message),
	)

	log.Info("Sending notification")

	err := n.ntfSender.SendNotification(ctx, userID, message)
	if err != nil {
		n.log.Error("failed to send notification", s1.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (n *Notification) SubscribeNotification(ctx context.Context, userID int64, message string) error {
	const op = "Notification.SubscribeNotification"

	log := n.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
		slog.String("message", message),
	)

	log.Info("Subscribing notification")

	err := n.ntfSubscriber.SubscribeNotification(ctx, userID, message)
	if err != nil {
		n.log.Error("failed to subscribe notification", s1.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
