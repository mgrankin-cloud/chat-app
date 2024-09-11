package app

import (
	"github.com/mgrankin-cloud/messenger/internal/services/auth"
	chat "github.com/mgrankin-cloud/messenger/internal/services/chat"
	"github.com/mgrankin-cloud/messenger/internal/services/media"
	message "github.com/mgrankin-cloud/messenger/internal/services/message"
	"github.com/mgrankin-cloud/messenger/internal/services/models"
	notification "github.com/mgrankin-cloud/messenger/internal/services/notification"
	"github.com/mgrankin-cloud/messenger/internal/services/user"
	"log/slog"
	"time"

	grpcapp "github.com/mgrankin-cloud/messenger/pkg/grpc/grpcapp"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    *storage.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	postgresConnString string,
	tokenTTL time.Duration,
) *App {
	strg, err := storage.New(postgresConnString)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, strg, strg, strg, tokenTTL)

	userService := user.New(log, strg, strg, strg, tokenTTL)

	messageService := message.New(log, strg, strg, strg, tokenTTL)

	chatService := chat.New(log, strg, strg, strg, tokenTTL)

	modelsService := models.New(log, strg, strg, strg, tokenTTL)

	mediaService := media.New(log, strg, strg, strg, tokenTTL)

	notificationService := notification.New(log, strg, strg, strg, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    strg,
	}
}

func (a *App) Stop() error {
	if err := a.GRPCServer.Stop(); err != nil {
		return err
	}

	if err := a.Storage.Stop(); err != nil {
		return err
	}

	return nil
}
