package management

import (
	"context"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
	"time"
)

type Service struct {
	log          *slog.Logger
	eventStorage EventStorage
}

type EventStorage interface {
	CreateEvent(ctx context.Context, dto *dto.CreateEventDTO) (*domain.Event, error)
	GetEvent(ctx context.Context, id int64) (*domain.Event, error)
}

func New(log *slog.Logger, eventStorage EventStorage) Service {
	return Service{
		log:          log,
		eventStorage: eventStorage,
	}
}

func (s Service) CreateEvent(ctx context.Context, dto *dto.CreateEventDTO) (*domain.Event, error) {
	const op = "services.event.management.CreateEvent"
	log := s.log.With(slog.String("op", op))

	createCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	event, err := s.eventStorage.CreateEvent(createCtx, dto)
	if err != nil {
		log.Error("failed to create event", logger.Err(err))
		return nil, err
	}

	return event, nil
}

func (s Service) GetEvent(ctx context.Context, dto *dto.GetEventDTO) (*domain.Event, error) {
	//TODO implement me
	panic("implement me")
}
