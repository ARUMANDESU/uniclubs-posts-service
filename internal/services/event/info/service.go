package info

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
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
	ListEvents(ctx context.Context, filters domain.Filters) ([]domain.Event, *domain.PaginationMetadata, error)
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
			return nil, eventservice.ErrEventNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventservice.ErrInvalidID
		default:
			log.Error("failed to get event", logger.Err(err))
			return nil, err
		}
	}

	if !event.IsOrganizer(userId) && event.Status == domain.EventStatusDraft {
		return nil, eventservice.ErrEventNotFound
	}

	return event, nil
}

func (s Service) ListEvents(ctx context.Context, filters domain.Filters) ([]domain.Event, *domain.PaginationMetadata, error) {
	const op = "services.event.management.listEvents"
	log := s.log.With(slog.String("op", op))

	events, pagination, err := s.eventProvider.ListEvents(ctx, filters)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEventNotFound):
			return nil, nil, eventservice.ErrEventNotFound
		default:
			log.Error("failed to list events", logger.Err(err))
			return nil, nil, err
		}
	}

	return events, pagination, nil
}
