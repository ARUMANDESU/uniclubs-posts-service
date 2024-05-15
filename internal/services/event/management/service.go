package management

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"log/slog"
	"time"

	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
)

var (
	ErrClubNotExists       = errors.New("club not found")
	ErrEventNotFound       = errors.New("event not found")
	ErrEventUpdateConflict = errors.New("event update conflict")
	ErrUserIsNotEventOwner = errors.New("permissions denied: user is not event owner")
)

type Service struct {
	log          *slog.Logger
	eventStorage EventStorage
	clubProvider ClubProvider
	userProvider UserProvider
}

type EventStorage interface {
	CreateEvent(ctx context.Context, clubId, userId int64) (*domain.Event, error)
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
}

type ClubProvider interface {
	GetClubByID(ctx context.Context, id int64) (*domain.Club, error)
}

type UserProvider interface {
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

func New(log *slog.Logger, eventStorage EventStorage, clubProvider ClubProvider, userProvider UserProvider) Service {
	return Service{
		log:          log,
		eventStorage: eventStorage,
		clubProvider: clubProvider,
		userProvider: userProvider,
	}
}

func (s Service) CreateEvent(ctx context.Context, clubId, userId int64) (*domain.Event, error) {
	const op = "services.event.management.createEvent"
	log := s.log.With(slog.String("op", op))

	_, err := s.clubProvider.GetClubByID(ctx, clubId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrClubNotExists):
			return nil, ErrClubNotExists
		default:
			log.Error("failed to get club", logger.Err(err))
			return nil, err
		}
	}

	createCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	event, err := s.eventStorage.CreateEvent(createCtx, clubId, userId)
	if err != nil {
		log.Error("failed to create event", logger.Err(err))
		return nil, err
	}

	return event, nil
}

func (s Service) GetEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.getEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, eventId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEventNotFound):
			return nil, ErrEventNotFound
		default:
			log.Error("failed to get event", logger.Err(err))
			return nil, err
		}
	}

	if event.User.ID != userId && event.Status == domain.EventStatusDraft {
		return nil, ErrEventNotFound
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
			return nil, ErrEventNotFound
		default:
			log.Error("failed to get event", logger.Err(err))
			return nil, err
		}
	}

	if event.User.ID != dto.UserId {
		return nil, ErrUserIsNotEventOwner
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
			return nil, ErrEventUpdateConflict
		default:
			log.Error("failed to update event", logger.Err(err))
			return nil, err
		}

	}

	return updatedEvent, nil
}
