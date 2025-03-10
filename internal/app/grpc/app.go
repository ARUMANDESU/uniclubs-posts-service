package grpcapp

import (
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/grpc/event"
	postgrpc "github.com/arumandesu/uniclubs-posts-service/internal/grpc/post"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, eventServices eventgrpc.Services, postServices postgrpc.Services) *App {
	gRPCServer := grpc.NewServer()

	eventgrpc.Register(gRPCServer, eventServices)
	postgrpc.Register(gRPCServer, postServices)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.grpc.run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPCServer is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "app.grpc.stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC Server")
	a.gRPCServer.GracefulStop()
}
