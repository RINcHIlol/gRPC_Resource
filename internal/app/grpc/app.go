package grpcapp

import (
	resourceGrpc "ResourceService/internal/grpc/resource"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, resourceService resourceGrpc.ResourceServer) *App {
	grpcServer := grpc.NewServer()

	resourceGrpc.Register(grpcServer, resourceService)

	return &App{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.log.Info("grpc server starting")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}
	if err := a.grpcServer.Serve(l); err != nil {
		return err
	}
	a.log.Info("grpc server started")
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op))

	a.grpcServer.GracefulStop()
}
