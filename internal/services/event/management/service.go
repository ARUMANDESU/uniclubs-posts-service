package eventmanagement

import (
	"context"
	"errors"
	"fmt"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"log/slog"
	"time"

	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
)

type Service struct {
	log          *slog.Logger
	eventStorage EventStorage
}

//go:generate mockery --name EventStorage
type EventStorage interface {
	CreateEvent(ctx context.Context, club domain.Club, user domain.User) (*domain.Event, error)
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
	event, err := s.eventStorage.CreateEvent(createCtx, club, user)
	if err != nil {
		log.Error("failed to create event", logger.Err(err))
		return nil, err
	}

	return event, nil
}

func (s Service) UpdateEvent(ctx context.Context, dto *dtos.UpdateEvent) (*domain.Event, error) {
	const op = "services.event.management.updateEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.FetchEventAndCheckOwner(ctx, dto.EventId, dto.UserId)
	if err != nil {
		return nil, err
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
			startDate, err := time.Parse(domain.TimeLayout, dto.StartDate)
			if err != nil {
				log.Error("failed to parse start date", logger.Err(err))
				return nil, err
			}
			event.StartDate = startDate
		case "end_date":
			endDate, err := time.Parse(domain.TimeLayout, dto.EndDate)
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
		return nil, s.handleError("failed to update event", log, err)
	}

	return updatedEvent, nil
}

func (s Service) DeleteEvent(ctx context.Context, dto *dtos.DeleteEvent) (*domain.Event, error) {
	const op = "services.event.management.deleteEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.getEvent(ctx, dto.EventId)
	if err != nil {
		return nil, err
	}

	if !(event.IsOwner(dto.UserId) || dto.IsAdmin) {
		return nil, eventservice.ErrPermissionsDenied
	}

	deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = s.eventStorage.DeleteEventById(deleteCtx, dto.EventId)
	if err != nil {
		return nil, s.handleError("failed to delete event", log, err)
	}

	return event, nil
}

func (s Service) PublishEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.publishEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.FetchEventAndCheckOwner(ctx, eventId, userId)
	if err != nil {
		return nil, err
	}

	err = validate.PublishEvent(event)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", eventservice.ErrEventInvalidFields, err)
	}

	err = event.Publish()
	if err != nil {
		return nil, s.handleError("failed to check if event can be published", log, err)
	}

	updateCtx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	updatedEvent, err := s.eventStorage.UpdateEvent(updateCtx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return updatedEvent, nil
}

func (s Service) SendToReview(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.sendToReview"
	log := s.log.With(slog.String("op", op))

	event, err := s.FetchEventAndCheckOwner(ctx, eventId, userId)
	if err != nil {
		return nil, err
	}

	err = validate.SendToReview(event)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", eventservice.ErrEventInvalidFields, err)
	}

	err = event.SendToReview()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", eventservice.ErrInvalidEventStatus, err)
	}

	updateCtx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	updatedEvent, err := s.eventStorage.UpdateEvent(updateCtx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return updatedEvent, nil
}

func (s Service) RevokeReview(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.revokeReview"
	log := s.log.With(slog.String("op", op))

	event, err := s.FetchEventAndCheckOwner(ctx, eventId, userId)
	if err != nil {
		return nil, err
	}

	err = event.RevokeReview()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", eventservice.ErrInvalidEventStatus, err)
	}

	updateCtx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	updatedEvent, err := s.eventStorage.UpdateEvent(updateCtx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return updatedEvent, nil
}

func (s Service) ApproveEvent(ctx context.Context, eventId string, user domain.User) (*domain.Event, error) {
	const op = "services.event.management.approveEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.getEvent(ctx, eventId)
	if err != nil {
		return nil, err
	}

	err = event.Approve(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", eventservice.ErrInvalidEventStatus, err)
	}

	updateCtx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	updatedEvent, err := s.eventStorage.UpdateEvent(updateCtx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return updatedEvent, nil
}

func (s Service) RejectEvent(ctx context.Context, dto *dtos.RejectEvent) (*domain.Event, error) {
	const op = "services.event.management.rejectEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.getEvent(ctx, dto.EventId)
	if err != nil {
		return nil, err
	}

	err = event.Reject(dto.User, dto.Reason)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", eventservice.ErrInvalidEventStatus, err)
	}

	updateCtx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	updatedEvent, err := s.eventStorage.UpdateEvent(updateCtx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return updatedEvent, nil
}

// FetchEventAndCheckOwner fetches an event and checks if the user is the owner
func (s Service) FetchEventAndCheckOwner(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.fetchEventAndCheckOwner"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, eventId)
	if err != nil {
		return nil, s.handleError("failed to get event", log, err)
	}

	if !event.IsOwner(userId) {
		return nil, eventservice.ErrUserIsNotEventOwner
	}

	return event, nil
}

func (s Service) getEvent(ctx context.Context, eventId string) (*domain.Event, error) {
	const op = "services.event.management.getEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, eventId)
	if err != nil {
		return nil, s.handleError("failed to get event", log, err)
	}

	return event, nil
}

// handleError handles common errors
func (s Service) handleError(msg string, log *slog.Logger, err error) error {
	switch {
	case errors.Is(err, storage.ErrEventNotFound):
		return eventservice.ErrEventNotFound
	case errors.Is(err, storage.ErrOptimisticLockingFailed):
		return eventservice.ErrEventUpdateConflict
	case errors.Is(err, storage.ErrInvalidID):
		return eventservice.ErrInvalidID
	case errors.Is(err, domain.ErrEventIsNotApproved):
		return eventservice.ErrEventIsNotApproved
	default:
		log.Error(msg, logger.Err(err))
		return err
	}
}
