package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//TODO: Инициализация объекста конфига
	ctx := context.Background()
	cfg := config.MustLoad()

	// TODO: инициализация логгера
	log := setupLogger(cfg.Env)

	log.Info("starting application")

	//TODO: Инициализировать приложение (app)
	application := app.New(ctx, log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()
	//TODO: Запустить gRPC-сервер приложения

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sigNtf := <-stop

	log.Info("spotting application", slog.String("signal", sigNtf.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
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
