package main

import (
	"fmt"
	"github.com/mgrankin-cloud/messenger/internal/config/user"
	app "github.com/mgrankin-cloud/messenger/pkg/grpc/user"
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
	usrServerAddr string = "http://localhost:8088"
	c                    = cache.New(5*time.Minute, 10*time.Minute)
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

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

func usrServerHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Handling request from user server: %s %s", req.Method, req.URL.Path)
	if response, found := c.Get("serverUser"); found {
		fmt.Fprintf(w, response.(string))
	} else {
		response := "hello from server user"
		c.Set("serverUser", response, cache.DefaultExpiration)
		fmt.Fprintf(w, response)
	}
}

var ServerUser *http.Server = &http.Server{
	Addr:    fmt.Sprintf(":%s", strings.Split(usrServerAddr, ":")[2]),
	Handler: http.HandlerFunc(usrServerHandler),
}
