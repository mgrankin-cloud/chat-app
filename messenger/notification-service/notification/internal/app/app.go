package app

import (
	notification "Messenger-android/messenger/notification-service/notification/internal/services/notification"
	"log/slog"
	"time"

	"Messenger-android/messenger/internal/storage"
	grpcapp "Messenger-android/messenger/notification-service/notification/internal/app/grpc"
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

	notificationService := notification.New(log, strg, strg, strg, tokenTTL)

	grpcApp := grpcapp.New(log, notificationService, grpcPort)

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
