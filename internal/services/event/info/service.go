package info

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	eventService "github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
)

type Service struct {
	log           *slog.Logger
	eventProvider EventProvider
}

type EventProvider interface {
	GetEvent(ctx context.Context, eventId string) (*domain.Event, error)
}

func New(log *slog.Logger, eventProvider EventProvider) Service {
	return Service{
		log:           log,
		eventProvider: eventProvider,
	}
}

func (s Service) GetEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.getEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventProvider.GetEvent(ctx, eventId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEventNotFound):
			return nil, eventService.ErrEventNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventService.ErrInvalidID
		default:
			log.Error("failed to get event", logger.Err(err))
			return nil, err
		}
	}

	if event.User.ID != userId && event.Status == domain.EventStatusDraft {
		return nil, eventService.ErrEventNotFound
	}

	return event, nil
}
