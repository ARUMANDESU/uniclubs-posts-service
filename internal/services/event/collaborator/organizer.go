package eventcollab

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

func (s Service) SendJoinRequestToUser(ctx context.Context, dto *dtos.SendJoinRequestToUser) (*domain.Event, error) {
	const op = "services.event.management.sendJoinRequestToUser"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, dto.EventId)
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

	// Check if the user is event owner or organizer
	if !event.IsOrganizer(dto.UserId) {
		return nil, eventservice.ErrPermissionsDenied
	}

	// Check if the target user is already an organizer
	if event.IsOrganizer(dto.Target.ID) {
		return nil, eventservice.ErrUserAlreadyOrganizer
	}

	userInvite, err := s.userInviteStorage.GetJoinRequestByUserId(ctx, dto.EventId, dto.Target.ID)
	if err != nil && !errors.Is(err, storage.ErrInviteNotFound) {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventservice.ErrInvalidID
		default:
			log.Error("failed to get join request by user id", logger.Err(err))
			return nil, err
		}
	}

	if userInvite != nil {
		return nil, eventservice.ErrInviteAlreadyExists
	}

	organizer := event.GetOrganizerById(dto.UserId)
	if organizer == nil {
		return nil, eventservice.ErrPermissionsDenied
	}

	if organizer.ClubId != dto.TargetClubId {
		return nil, eventservice.ErrUserIsFromAnotherClub
	}

	sendJoinRequestCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = s.userInviteStorage.CreateJoinRequestToUser(sendJoinRequestCtx, dto)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventservice.ErrInvalidID
		default:
			log.Error("failed to send join request", logger.Err(err))
			return nil, err
		}
	}

	// todo: send push notification to the target user

	return event, nil
}

func (s Service) AcceptUserJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error) {
	const op = "services.event.management.acceptUserJoinRequest"
	log := s.log.With(slog.String("op", op))

	invite, err := s.userInviteStorage.GetJoinRequestsByUserInviteId(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInviteNotFound):
			return domain.Event{}, eventservice.ErrInviteNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return domain.Event{}, eventservice.ErrInvalidID
		default:
			log.Error("failed to get join request by id", logger.Err(err))
			return domain.Event{}, err
		}
	}

	if !invite.IsInvited(userId) {
		return domain.Event{}, eventservice.ErrPermissionsDenied
	}

	event, err := s.eventStorage.GetEvent(ctx, invite.EventId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEventNotFound):
			return domain.Event{}, eventservice.ErrEventNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return domain.Event{}, eventservice.ErrInvalidID
		default:
			log.Error("failed to get event", logger.Err(err))
			return domain.Event{}, err
		}
	}

	if event.IsOrganizer(invite.User.ID) {
		return domain.Event{}, eventservice.ErrUserAlreadyOrganizer
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = s.inviteDeleter.DeleteInvite(deleteCtx, inviteId)
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		newOrganizer := domain.Organizer{
			User:    invite.User,
			ClubId:  invite.ClubId,
			ByWhoId: invite.ByWhoId,
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
				return domain.Event{}, eventservice.ErrInviteNotFound
			case errors.Is(err, storage.ErrInvalidID):
				return domain.Event{}, eventservice.ErrInvalidID
			case errors.Is(err, storage.ErrOptimisticLockingFailed):
				return domain.Event{}, eventservice.ErrEventUpdateConflict
			default:
				log.Error("failed to accept join request", logger.Err(err))
				return domain.Event{}, err
			}
		}
	}
	return *event, nil
}

func (s Service) RejectUserJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error) {
	const op = "services.event.collaborator.rejectUserJoinRequest"
	log := s.log.With(slog.String("op", op))

	invite, err := s.userInviteStorage.GetJoinRequestsByUserInviteId(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInviteNotFound):
			return domain.Event{}, eventservice.ErrInviteNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return domain.Event{}, eventservice.ErrInvalidID
		default:
			log.Error("failed to get join request by id", logger.Err(err))
			return domain.Event{}, err
		}
	}
	if !invite.IsInvited(userId) {
		return domain.Event{}, eventservice.ErrPermissionsDenied
	}

	deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.inviteDeleter.DeleteInvite(deleteCtx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInviteNotFound):
			return domain.Event{}, eventservice.ErrInviteNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return domain.Event{}, eventservice.ErrInvalidID
		default:
			log.Error("failed to delete join request", logger.Err(err))
			return domain.Event{}, err
		}
	}

	return domain.Event{}, nil
}

func (s Service) KickOrganizer(ctx context.Context, eventId string, userId, targetId int64) (*domain.Event, error) {
	const op = "services.event.collaborator.kickOrganizer"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, eventId)
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

	if !event.IsOrganizer(userId) {
		return nil, eventservice.ErrPermissionsDenied
	}

	target := event.GetOrganizerById(targetId)
	if target == nil {
		return nil, eventservice.ErrUserIsNotEventOrganizer
	}

	if !(target.IsByWho(userId) || event.IsOwner(userId)) {
		return nil, eventservice.ErrPermissionsDenied
	}

	err = event.RemoveOrganizer(targetId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrganizerNotFound):
			return nil, eventservice.ErrUserIsNotEventOrganizer
		case errors.Is(err, domain.ErrOrganizersEmpty):
			return nil, eventservice.ErrOrganizerNotFound
		default:
			log.Error("failed to remove organizer", logger.Err(err))
			return nil, err
		}
	}

	updateCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	event, err = s.eventStorage.UpdateEvent(updateCtx, event)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrOptimisticLockingFailed):
			return nil, eventservice.ErrEventUpdateConflict
		default:
			log.Error("failed to update event", logger.Err(err))
			return nil, err
		}
	}

	return event, nil
}

func (s Service) RevokeInviteOrganizer(ctx context.Context, inviteId string, userId int64) error {
	const op = "services.event.collaborator.revokeInviteOrganizer"
	log := s.log.With(slog.String("op", op))

	invite, err := s.userInviteStorage.GetJoinRequestsByUserInviteId(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInviteNotFound):
			return eventservice.ErrInviteNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return eventservice.ErrInvalidID
		default:
			log.Error("failed to get join request by id", logger.Err(err))
			return err
		}
	}

	event, err := s.eventStorage.GetEvent(ctx, invite.EventId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEventNotFound):
			return eventservice.ErrEventNotFound
		case errors.Is(err, storage.ErrInvalidID):
			return eventservice.ErrInvalidID
		default:
			log.Error("failed to get event", logger.Err(err))
			return err
		}
	}

	if !(invite.IsByWho(userId) || event.IsOwner(userId)) {
		return eventservice.ErrPermissionsDenied
	}

	deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = s.inviteDeleter.DeleteInvite(deleteCtx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return eventservice.ErrInvalidID
		default:
			log.Error("failed to delete join request", logger.Err(err))
			return err
		}
	}

	return nil
}
