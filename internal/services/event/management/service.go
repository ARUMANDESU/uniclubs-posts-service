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

	if dto.Paths["is_hidden_for_non_members"] && !(event.Type == domain.EventTypeIntraClub || dto.Type == domain.EventTypeIntraClub) {
		return nil, fmt.Errorf("%w: university scope event cannot be hidden from non member users", eventservice.ErrEventInvalidFields)
	}

	hasUnchangeableFields := dto.HasUnchangeableFields()

	switch event.Status {
	case domain.EventStatusFinished, domain.EventStatusCanceled, domain.EventStatusArchived:
		return nil, eventservice.ErrEventIsNotEditable
	case domain.EventStatusApproved:
		if hasUnchangeableFields {
			event.ChangeStatus(domain.EventStatusPending)
		}
	case domain.EventStatusInProgress, domain.EventStatusPending:
		if hasUnchangeableFields {
			return nil, fmt.Errorf("%w for status %s", eventservice.ErrContainsUnchangeable, event.Status)
		}
	case domain.EventStatusRejected, domain.EventStatusDraft:
		// Allow updates without additional checks
	default:
		return nil, fmt.Errorf("%w: %s", eventservice.ErrUnknownStatus, event.Status)
	}

	updateFunctions := map[string]func(){
		"title":                     func() { event.Title = dto.Title },
		"description":               func() { event.Description = dto.Description },
		"type":                      func() { event.Type = dto.Type },
		"tags":                      func() { event.Tags = dto.Tags },
		"max_participants":          func() { event.MaxParticipants = dto.MaxParticipants },
		"location_link":             func() { event.LocationLink = dto.LocationLink },
		"location_university":       func() { event.LocationUniversity = dto.LocationUniversity },
		"start_date":                func() { event.StartDate = dto.StartDate },
		"end_date":                  func() { event.EndDate = dto.EndDate },
		"cover_images":              func() { event.CoverImages = dto.CoverImages },
		"attached_images":           func() { event.AttachedImages = dto.AttachedImages },
		"attached_files":            func() { event.AttachedFiles = dto.AttachedFiles },
		"is_hidden_for_non_members": func() { event.IsHiddenForNonMembers = dto.IsHiddenForNonMembers },
	}

	for path, exists := range dto.Paths {
		if !exists {
			continue
		}
		if updateFunc, ok := updateFunctions[path]; ok {
			updateFunc()
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

func (s Service) UnpublishEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error) {
	const op = "services.event.management.unpublishEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.FetchEventAndCheckOwner(ctx, eventId, userId)
	if err != nil {
		return nil, err
	}

	err = event.Unpublish()
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
