package club

import (
	"context"
	"encoding/json"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

type Service struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	UpdateClub(ctx context.Context, club *domain.Club) error
}

func New(log *slog.Logger, storage Storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s Service) HandleUpdateClub(msg amqp091.Delivery) error {
	const op = "services.club.handleUpdateClub"
	log := s.log.With(slog.String("op", op))

	var club domain.Club
	err := json.Unmarshal(msg.Body, &club)
	if err != nil {
		log.Error("failed to unmarshal message", logger.Err(err))
		return err
	}

	updateClubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.storage.UpdateClub(updateClubCtx, &club)
	if err != nil {
		log.Error("failed to update club", logger.Err(err))
		return err
	}

	return nil
}
