package media

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var redisClient *redis.Client

type Media struct {
	log          *slog.Logger
	mdUploader   MediaUploader
	mdDownloader MediaDownloader
	tokenTTL     time.Duration
}

type MediaUploader interface {
	UploadMedia(ctx context.Context, data []byte, fileName string, mimeType string) (models.Media, error)
}

type MediaDownloader interface {
	DownloadMedia(ctx context.Context, fileID int64) (models.Media, error)
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
	tokenTTL time.Duration,
) *Media {
	return &Media{
		mdUploader:   mediaUploader,
		mdDownloader: mediaDownloader,
		log:          log,
		tokenTTL:     tokenTTL,
	}
}

func (m *Media) UploadMedia(
	ctx context.Context,
	data []byte,
	fileName string,
	mimeType string,
) (bool, int64, error) {
	const op = "Media.UploadMedia"

	log := m.log.With(
		slog.String("op", op),
		slog.String("fileName", fileName),
		slog.String("mimeType", mimeType),
	)

	log.Info("Uploading media content")

	media, err := m.mdUploader.UploadMedia(ctx, data, fileName, mimeType)
	if err != nil {
		log.Error("failed to upload media", slog.String("error", err.Error()))
		return false, 0, fmt.Errorf("%s: %w", op, err)
	}

	mediaJSON, err := json.Marshal(media)
	if err != nil {
		log.Error("failed to marshal media metadata", slog.String("error", err.Error()))
	} else {
		err = redisClient.Set(ctx, fmt.Sprintf("media:%d", media.ID), mediaJSON, m.tokenTTL).Err()
		if err != nil {
			log.Error("failed to cache media metadata", slog.String("error", err.Error()))
		}
	}

	return true, media.ID, nil
}

func (m *Media) DownloadMedia(
	ctx context.Context,
	fileID int64,
) ([]byte, string, string, error) {
	const op = "Media.DownloadMedia"

	log := m.log.With(
		slog.String("op", op),
		slog.Int64("file_id", fileID),
	)

	log.Info("Downloading media content")

	cacheKey := fmt.Sprintf("media:%d", fileID)
	mediaJSON, err := redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var cachedMedia models.Media
		err = json.Unmarshal(mediaJSON, &cachedMedia)
		if err == nil {
			return cachedMedia.Data, cachedMedia.FileName, cachedMedia.MimeType, nil
		}
	}

	media, err := m.mdDownloader.DownloadMedia(ctx, fileID)
	if err != nil {
		log.Error("failed to download media", slog.String("error", err.Error()))
		return nil, "", "", fmt.Errorf("%s: %w", op, err)
	}

	mediaJSON, err = json.Marshal(media)
	if err != nil {
		log.Error("failed to marshal media metadata", slog.String("error", err.Error()))
	} else {
		err = redisClient.Set(ctx, cacheKey, mediaJSON, m.tokenTTL).Err()
		if err != nil {
			log.Error("failed to cache media metadata", slog.String("error", err.Error()))
		}
	}

	return media.Data, media.FileName, media.MimeType, nil
}