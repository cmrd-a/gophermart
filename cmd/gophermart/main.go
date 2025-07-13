package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cmrd-a/gophermart/migrations"

	"github.com/cmrd-a/gophermart/internal/api"
	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/cmrd-a/gophermart/internal/repository"
	"github.com/cmrd-a/gophermart/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	config.InitConfig()
	migrations.Migrate()
	repo, err := repository.NewRepository()
	if err != nil {
		slog.Error("Failed to create repository", "error", err.Error())
		os.Exit(1)
	}
	svc := service.NewService(context.TODO(), *repo)
	r := api.SetupRouter(svc)
	err = r.Run(config.Config.RunAddress)
	if err != nil {
		slog.Error("Failed to start server", "error", err.Error())
		os.Exit(1)
	}
}
