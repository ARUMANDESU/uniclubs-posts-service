package main

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/app"
	"github.com/arumandesu/uniclubs-posts-service/internal/config"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Error("error loading .env file")
	}

	cfg := config.MustLoad()
	log := logger.Setup(cfg.Env)

	log.Info("starting application",
		slog.String("env", cfg.Env),
		slog.Int("port", cfg.GRPC.Port),
	)

	application := app.New(log, cfg)

	go application.GRPCSrv.MustRun()
	application.AMQPApp.SetupMessageConsumers()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	defer log.Info("application stopped", slog.String("signal", sign.String()))
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.Stop()
}
