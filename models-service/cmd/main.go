package main

import (
	"fmt"
	"github.com/mgrankin-cloud/messenger/internal/config/models"
	app "github.com/mgrankin-cloud/messenger/pkg/grpc/models"
	"github.com/patrickmn/go-cache"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var (
	modelsServerAddr string = "http://localhost:8086"
	c                       = cache.New(5*time.Minute, 10*time.Minute)
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath)

	go func() {
		application.GRPCServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Waiting for SIGINT (pkill -2) or SIGTERM
	<-stop

	// initiate graceful shutdown
	err := application.GRPCServer.Stop()
	if err != nil {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		logger.Error("Failed to stop gRPC server", "error", err)
	}

	log.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func modelsServerHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Handling request from models server: %s %s", req.Method, req.URL.Path)
	if response, found := c.Get("serverModels"); found {
		fmt.Fprintf(w, response.(string))
	} else {
		response := "hello from server models"
		c.Set("serverModels", response, cache.DefaultExpiration)
		fmt.Fprintf(w, response)
	}
}

var ServerModels *http.Server = &http.Server{
	Addr:    fmt.Sprintf(":%s", strings.Split(modelsServerAddr, ":")[2]),
	Handler: http.HandlerFunc(modelsServerHandler),
}