package app

import (
	grpcapp "ResourceService/internal/app/grpc"
	"ResourceService/internal/services/resourceService"
	"ResourceService/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}

	resourceService := resourceService.New(log, storage, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, resourceService)
	return &App{
		GRPCServer: grpcApp,
	}
}
