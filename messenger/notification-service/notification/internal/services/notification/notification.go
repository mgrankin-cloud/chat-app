package media

import (
	"Messenger-android/messenger/internal/domain/models"

	"errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log/slog"
	"time"
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
	SendNotification(ctx context.Context) error
}

type NotificationSubscriber interface {
	SubscribeNotification(ctx context.Context) error
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
