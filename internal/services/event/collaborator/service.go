package collaborator

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	eventService "github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

type Service struct {
	log                    *slog.Logger
	eventStorage           EventStorage
	organizerInviteStorage OrganizerInviteStorage
}

type EventStorage interface {
	GetEvent(ctx context.Context, eventId string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
}

type OrganizerInviteStorage interface {
	SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.UserInvite, error)
	GetJoinRequests(ctx context.Context, eventId string) ([]domain.UserInvite, error)
	GetJoinRequestByUserId(ctx context.Context, userId int64) (*domain.UserInvite, error)
	GetJoinRequestsById(ctx context.Context, requestId string) (*domain.UserInvite, error)
	DeleteJoinRequest(ctx context.Context, requestId string) error
}

func New(log *slog.Logger, eventProvider EventStorage, organizerInviteStorage OrganizerInviteStorage) Service {
	return Service{
		log:                    log,
		eventStorage:           eventProvider,
		organizerInviteStorage: organizerInviteStorage,
	}
}

func (s Service) SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.Event, error) {
	const op = "services.event.management.sendJoinRequestToUser"
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

	// Check if the user is event owner or organizer
	if !event.IsOrganizer(dto.UserId) {
		return nil, eventService.ErrPermissionsDenied
	}

	// Check if the target user is already an organizer
	if event.IsOrganizer(dto.Target.ID) {
		return nil, eventService.ErrUserAlreadyOrganizer
	}

	userInvite, err := s.organizerInviteStorage.GetJoinRequestByUserId(ctx, dto.Target.ID)
	if err != nil && !errors.Is(err, storage.ErrInviteNotFound) {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventService.ErrInvalidID
		default:
			log.Error("failed to get join request by user id", logger.Err(err))
			return nil, err
		}
	}

	if userInvite != nil {
		return nil, eventService.ErrInviteAlreadyExists
	}

	// if the target user is the event owner then check if the target club is the same as the event club
	// if the target user is an organizer then check if the target club is the same as the organizer club
	if event.OwnerId == dto.UserId {
		if event.OwnerId != dto.TargetClubId {
			return nil, eventService.ErrUserIsFromAnotherClub
		}
	} else {
		organizer := event.GetOrganizerById(dto.UserId)
		if organizer == nil {
			return nil, eventService.ErrPermissionsDenied
		}

		if organizer.ClubId != dto.TargetClubId {
			return nil, eventService.ErrUserIsFromAnotherClub
		}
	}

	sendJoinRequestCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = s.organizerInviteStorage.SendJoinRequestToUser(sendJoinRequestCtx, dto)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventService.ErrInvalidID
		default:
			log.Error("failed to send join request", logger.Err(err))
			return nil, err
		}
	}

	// todo: send push notification to the target user

	return event, nil
}

func (s Service) AcceptUserJoinRequest(ctx context.Context, inviteId string, userId int64) error {
	const op = "services.event.management.acceptUserJoinRequest"
	log := s.log.With(slog.String("op", op))

	invite, err := s.organizerInviteStorage.GetJoinRequestsById(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInviteNotFound):
			return eventService.ErrInviteNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return eventService.ErrInvalidID
		default:
			log.Error("failed to get join request by id", logger.Err(err))
			return err
		}
	}

	if !invite.IsInvited(userId) {
		return eventService.ErrPermissionsDenied
	}

	event, err := s.eventStorage.GetEvent(ctx, invite.EventId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEventNotFound):
			return eventService.ErrEventNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return eventService.ErrInvalidID
		default:
			log.Error("failed to get event", logger.Err(err))
			return err
		}
	}

	if event.IsOrganizer(invite.User.ID) {
		return eventService.ErrUserAlreadyOrganizer
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = s.organizerInviteStorage.DeleteJoinRequest(deleteCtx, inviteId)
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		newOrganizer := domain.Organizer{
			User:   invite.User,
			ClubId: invite.ClubId,
		}
		event.AddOrganizer(newOrganizer)

		updateCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_, err = s.eventStorage.UpdateEvent(updateCtx, event)
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()

	close(errCh)

	// Check if there were any errors
	for err := range errCh {
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrInviteNotFound):
				return eventService.ErrInviteNotFound
			case errors.Is(err, storage.ErrInvalidID):
				return eventService.ErrInvalidID
			case errors.Is(err, storage.ErrOptimisticLockingFailed):
				return eventService.ErrEventUpdateConflict
			default:
				log.Error("failed to accept join request", logger.Err(err))
				return err
			}
		}
	}
	return nil
}

func (s Service) RejectUserJoinRequest(ctx context.Context, inviteId string, userId int64) error {
	const op = "services.event.management.rejectUserJoinRequest"
	log := s.log.With(slog.String("op", op))

	invite, err := s.organizerInviteStorage.GetJoinRequestsById(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInviteNotFound):
			return eventService.ErrInviteNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return eventService.ErrInvalidID
		default:
			log.Error("failed to get join request by id", logger.Err(err))
			return err
		}
	}
	if !invite.IsInvited(userId) {
		return eventService.ErrPermissionsDenied
	}

	deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.organizerInviteStorage.DeleteJoinRequest(deleteCtx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInviteNotFound):
			return eventService.ErrInviteNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return eventService.ErrInvalidID
		default:
			log.Error("failed to delete join request", logger.Err(err))
			return err
		}
	}

	return nil
}
