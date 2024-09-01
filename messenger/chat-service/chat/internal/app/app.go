package app

import (
	grpcapp "Messenger-android/messenger/chat-service/chat/internal/app/grpc"
	"Messenger-android/messenger/internal/storage"
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

	chatService := chat.New(log, strg, strg, strg, strg, strg)

	grpcApp := grpcapp.New(log, chatService, grpcPort)

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
