package app

import (
	"github.com/mgrankin-cloud/messenger/message-service/internal/services/message"

	"log/slog"

	grpcapp "github.com/mgrankin-cloud/messenger/message-service/pkg/grpc/grpcapp/message"
	"github.com/mgrankin-cloud/messenger/shared/storage"
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

	msgService := message.New(log, strg, strg)

	grpcApp := grpcapp.New(log, msgService, grpcPort)

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
