package main

import (
	"log/slog"
	"os"
	"sso/internal/app"
	"sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//TODO: Инициализация объекста конфига
	cfg := config.MustLoad()

	// TODO: инициализация логгера
	log := setupLogger(cfg.Env)

	log.Info("starting application")
	log.Debug("Something")

	//TODO: Инициализировать приложение (app)
	application := app.New(log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GRPCSrv.MustRun()
	//TODO: Запустить gRPC-сервер приложения
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
