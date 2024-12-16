package app

import (
	"github.com/mgrankin-cloud/messenger/internal/services/user"

	"log/slog"
	"time"

	grpcapp "github.com/mgrankin-cloud/messenger/pkg/grpc/grpcapp/user"
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

	userService := user.New(log, strg, strg)

	grpcApp := grpcapp.New(log, userService, grpcPort)

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
