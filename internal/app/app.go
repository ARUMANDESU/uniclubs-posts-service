package app

import (
	"context"
	amqpapp "github.com/arumandesu/uniclubs-posts-service/internal/app/amqp"
	grpcapp "github.com/arumandesu/uniclubs-posts-service/internal/app/grpc"
	"github.com/arumandesu/uniclubs-posts-service/internal/config"
	"github.com/arumandesu/uniclubs-posts-service/internal/grpc/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/rabbitmq"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/club"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event/collaborator"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event/info"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event/management"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/user"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage/mongodb"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
	"time"
)

type App struct {
	log     *slog.Logger
	GRPCSrv *grpcapp.App
	AMQPApp *amqpapp.App
	mongoDB *mongodb.Storage
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

	rmq, err := rabbitmq.New(cfg.Rabbitmq, log)
	if err != nil {
		l.Error("failed to connect to rabbitmq", logger.Err(err))
		panic(err)
	}
	l.Info("connected to rabbitmq")

	userService := userservice.New(log, mongoDB)
	clubService := clubservice.New(log, mongoDB)
	eventCollaboratorService := eventcollab.New(log, mongoDB, mongoDB, mongoDB, mongoDB)

	services := eventgrpc.NewServices(
		eventmanagement.New(log, mongoDB),
		eventCollaboratorService,
		eventCollaboratorService,
		eventinfo.New(log, mongoDB),
	)

	grpcApp := grpcapp.New(log, cfg.GRPC.Port, services)
	amqpApp := amqpapp.New(log, userService, clubService, rmq)
	return &App{
		log:     log,
		GRPCSrv: grpcApp,
		AMQPApp: amqpApp,
		mongoDB: mongoDB,
	}
}

func (a *App) Stop() {
	const op = "app.stop"
	log := a.log.With(slog.String("op", op))

	a.GRPCSrv.Stop()

	err := a.AMQPApp.Shutdown()
	if err != nil {
		log.Error("failed to shutdown amqp app", logger.Err(err))
	}

	mongoCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = a.mongoDB.Close(mongoCtx)
	if err != nil {
		log.Error("failed to close mongodb connection", logger.Err(err))
	}
}
