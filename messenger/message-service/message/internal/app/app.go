package app

import (
	"Messenger-android/messenger/internal/storage"
	grpcapp "Messenger-android/messenger/message-service/message/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    *storage.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	postgresConnString string,
) *App {
	strg, err := storage.New(postgresConnString)
	if err != nil {
		panic(err)
	}

	messageService := message.New(log, strg, strg, strg, strg, strg)

	grpcApp := grpcapp.New(log, messageService, grpcPort)

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
