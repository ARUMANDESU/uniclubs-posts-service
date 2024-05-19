package management

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	eventService "github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"log/slog"
	"time"

	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
)

type Service struct {
	log          *slog.Logger
	eventStorage EventStorage
}

type EventStorage interface {
	CreateEvent(ctx context.Context, club *domain.Club, user *domain.User) (*domain.Event, error)
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
	DeleteEventById(ctx context.Context, eventId string) error
}

func New(log *slog.Logger, eventStorage EventStorage) Service {
	return Service{log: log, eventStorage: eventStorage}
}

func (s Service) CreateEvent(ctx context.Context, club domain.Club, user domain.User) (*domain.Event, error) {
	const op = "services.event.management.createEvent"
	log := s.log.With(slog.String("op", op))

	createCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	event, err := s.eventStorage.CreateEvent(createCtx, &club, &user)
	if err != nil {
		log.Error("failed to create event", logger.Err(err))
		return nil, err
	}

	return event, nil
}

func (s Service) UpdateEvent(ctx context.Context, dto *dto.UpdateEvent) (*domain.Event, error) {
	const op = "services.event.management.updateEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, dto.EventId)
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

	if event.OwnerId != dto.UserId {
		return nil, eventService.ErrUserIsNotEventOwner
	}

	for _, path := range dto.Paths {
		switch path {
		case "title":
			event.Title = dto.Title
		case "description":
			event.Description = dto.Description
		case "type":
			event.Type = dto.Type
		case "tags":
			event.Tags = dto.Tags
		case "max_participants":
			event.MaxParticipants = dto.MaxParticipants
		case "location_link":
			event.LocationLink = dto.LocationLink
		case "location_university":
			event.LocationUniversity = dto.LocationUniversity
		case "start_date":
			startDate, err := time.Parse(time.RFC3339, dto.StartDate)
			if err != nil {
				log.Error("failed to parse start date", logger.Err(err))
				return nil, err
			}
			event.StartDate = startDate
		case "end_date":
			endDate, err := time.Parse(time.RFC3339, dto.EndDate)
			if err != nil {
				log.Error("failed to parse end date", logger.Err(err))
				return nil, err
			}
			event.EndDate = endDate
		case "cover_images":
			event.CoverImages = dto.CoverImages
		case "attached_images":
			event.AttachedImages = dto.AttachedImages
		case "attached_files":
			event.AttachedFiles = dto.AttachedFiles
		}
	}

	updateCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	updatedEvent, err := s.eventStorage.UpdateEvent(updateCtx, event)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrOptimisticLockingFailed):
			return nil, eventService.ErrEventUpdateConflict
		default:
			log.Error("failed to update event", logger.Err(err))
			return nil, err
		}

	}

	return updatedEvent, nil
}

func (s Service) DeleteEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.deleteEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, eventId)
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
	if event.OwnerId != userId {
		return nil, eventService.ErrUserIsNotEventOwner
	}

	deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.eventStorage.DeleteEventById(deleteCtx, eventId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEventNotFound):
			return nil, eventService.ErrEventNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventService.ErrInvalidID
		default:
			log.Error("failed to delete event", logger.Err(err))
			return nil, err
		}
	}

	return event, nil
}
