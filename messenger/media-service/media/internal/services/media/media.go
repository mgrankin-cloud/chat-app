package media

import (
	"Messenger-android/messenger/internal/domain/models"

	errors2 "Messenger-android/messenger/internal/storage"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var redisClient *redis.Client

type Media struct {
	log         *slog.Logger
	mdUploader    MediaUploader
	mdDownloader MediaDownloader
	appProvider AppProvider
	tokenTTL    time.Duration
}

type MediaSaver interface {
	SaveMedia(
		ctx context.Context,

	) ( , err error)
}

type MediaUploader interface {
	UploadMedia(ctx context.Context, ) (models.Media, error)
}

type MediaDownloader interface {
	DownloadMedia(ctx context.Context, ) (models.Media, error)
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
	mediaUploader MediaUploader,
	mediaDownloader MediaDownloader,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Media {
	return &Media{
		mdUploader:    mediaUploader,
		mdDownloader:  mediaDownloader,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (m *Media) UploadMedia(
	ctx context.Context,
	) {

}

func (m *Media) DownloadMedia(
	ctx context.Context,
	) {

}
