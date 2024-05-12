package app

import (
	"context"
	amqpapp "github.com/arumandesu/uniclubs-posts-service/internal/app/amqp"
	grpcapp "github.com/arumandesu/uniclubs-posts-service/internal/app/grpc"
	"github.com/arumandesu/uniclubs-posts-service/internal/config"
	"github.com/arumandesu/uniclubs-posts-service/internal/rabbitmq"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/club"
	eventManagement "github.com/arumandesu/uniclubs-posts-service/internal/services/event/management"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/user"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage/mongodb"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
	AMQPApp *amqpapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	const op = "app.new"
	l := log.With(slog.String("op", op))

	mongoDB, err := mongodb.New(context.Background(), cfg.MongoDB)
	if err != nil {
		l.Error("failed to connect to mongodb", logger.Err(err))
		panic(err)
	}
	l.Info("connected to mongodb")

	userService := user.New(l, mongoDB)
	clubService := club.New(l, mongoDB)

	rmq, err := rabbitmq.New(cfg.Rabbitmq, l)
	if err != nil {
		l.Error("failed to connect to rabbitmq", logger.Err(err))
		panic(err)
	}

	managementService := eventManagement.New(l, mongoDB)

	grpcApp := grpcapp.New(l, cfg.GRPC.Port, managementService)

	amqpApp := amqpapp.New(l, userService, clubService, rmq)
	return &App{
		GRPCSrv: grpcApp,
		AMQPApp: amqpApp,
	}
}
