package app

import (
	"Messenger-android/messenger/internal/storage"
	grpcapp "Messenger-android/messenger/models-service/models/internal/app/grpc"
	"Messenger-android/messenger/models-service/models/internal/services/models"
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

	modelsService := models.New(log, strg, strg, strg, strg)

	grpcApp := grpcapp.New(log, modelsService, grpcPort)

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
