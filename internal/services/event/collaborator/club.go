package collaborator

import (
	"context"
	"errors"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

func (s Service) SendJoinRequestToClub(ctx context.Context, dto *dto.SendJoinRequestToClub) (*domain.Event, error) {
	const op = "services.event.collaborator.sendJoinRequestToClub"
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

	if !event.IsOwner(dto.UserId) {
		return nil, eventservice.ErrPermissionsDenied
	}

	if event.IsCollaborator(dto.Club.ID) {
		return nil, eventservice.ErrClubAlreadyCollaborator
	}

	invite, err := s.clubInviteStorage.GetJoinRequestByClubId(ctx, dto.EventId, dto.Club.ID)
	if err != nil && !errors.Is(err, storage.ErrInviteNotFound) {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventservice.ErrInvalidID
		default:
			log.Error("failed to get join request by club id", logger.Err(err))
			return nil, err
		}
	}

	if invite != nil {
		return nil, eventservice.ErrInviteAlreadyExists
	}

	createCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	invite, err = s.clubInviteStorage.CreateJoinRequestToClub(createCtx, dto)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventservice.ErrInvalidID
		default:
			log.Error("failed to get join request by club id", logger.Err(err))
			return nil, err
		}
	}

	return event, nil
}

func (s Service) AcceptClubJoinRequest(ctx context.Context, dto *dto.AcceptJoinRequestClub) (domain.Event, error) {
	const op = "services.event.collaborator.acceptClubJoinRequest"
	log := s.log.With(slog.String("op", op))

	invite, err := s.clubInviteStorage.GetJoinRequestsByClubInviteId(ctx, dto.InviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return domain.Event{}, eventservice.ErrInvalidID
		case errors.Is(err, storage.ErrInviteNotFound):
			return domain.Event{}, eventservice.ErrInviteNotFound
		default:
			log.Error("failed to get join requests by club invite id", logger.Err(err))
			return domain.Event{}, err
		}
	}

	if invite.Club.ID != dto.ClubId {
		return domain.Event{}, fmt.Errorf("%w got %d", eventservice.ErrClubMismatch, dto.ClubId)
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

	if event.IsCollaborator(dto.ClubId) {
		return domain.Event{}, eventservice.ErrClubAlreadyCollaborator
	}

	if event.IsOrganizer(dto.User.ID) {
		return domain.Event{}, eventservice.ErrUserAlreadyOrganizer
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		deleteCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = s.inviteDeleter.DeleteInvite(deleteCtx, dto.InviteId)
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()

		event.AddCollaborator(invite.Club)
		event.AddOrganizer(dto.User.ToOrganizer(dto.ClubId, event.OwnerId))

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

func (s Service) RejectClubJoinRequest(ctx context.Context, inviteId string, clubId int64) (domain.Event, error) {
	const op = "services.event.collaborator.rejectClubJoinRequest"
	log := s.log.With(slog.String("op", op))

	invite, err := s.clubInviteStorage.GetJoinRequestsByClubInviteId(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return domain.Event{}, eventservice.ErrInvalidID
		case errors.Is(err, storage.ErrInviteNotFound):
			return domain.Event{}, eventservice.ErrInviteNotFound
		default:
			log.Error("failed to get join requests by club invite id", logger.Err(err))
			return domain.Event{}, err
		}
	}

	if invite.Club.ID != clubId {
		return domain.Event{}, fmt.Errorf("%w got %d", eventservice.ErrClubMismatch, clubId)
	}

	err = s.inviteDeleter.DeleteInvite(ctx, inviteId)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidID) {
			return domain.Event{}, eventservice.ErrInvalidID
		}
		log.Error("failed to delete invite", logger.Err(err))
		return domain.Event{}, err
	}

	return domain.Event{}, nil
}

func (s Service) KickClub(ctx context.Context, eventId string, userId, clubId int64) (*domain.Event, error) {
	const op = "services.event.collaborator.kickClub"
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

	if !event.IsOwner(userId) {
		return nil, eventservice.ErrPermissionsDenied
	}

	if !event.IsCollaborator(clubId) {
		return nil, eventservice.ErrCollaboratorNotFound
	}

	club := event.GetCollaboratorById(clubId)
	if club == nil {
		return nil, eventservice.ErrCollaboratorNotFound
	}

	err = event.RemoveCollaborator(clubId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrClubIsEventOwner):
			return nil, eventservice.ErrClubIsEventOwner
		case errors.Is(err, domain.ErrCollaboratorsEmpty), errors.Is(err, domain.ErrCollaboratorNotFound):
			return nil, eventservice.ErrCollaboratorNotFound
		default:
			log.Error("failed to remove collaborator", logger.Err(err))
			return nil, err
		}
	}

	//delete organizers from that club
	err = event.RemoveOrganizersByClubId(clubId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrganizerNotFound):
			return nil, eventservice.ErrOrganizerNotFound
		case errors.Is(err, domain.ErrUserIsEventOwner):
			return nil, eventservice.ErrUserIsEventOwner
		default:
			log.Error(fmt.Sprintf("failed to remove organizers with club id: %d", clubId), logger.Err(err))
			return nil, err
		}
	}

	updateCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

func (s Service) RevokeInviteClub(ctx context.Context, inviteId string, userId int64) error {
	const op = "services.event.collaborator.revokeInviteClub"
	log := s.log.With(slog.String("op", op))

	invite, err := s.clubInviteStorage.GetJoinRequestsByClubInviteId(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return eventservice.ErrInvalidID
		case errors.Is(err, storage.ErrInviteNotFound):
			return eventservice.ErrInviteNotFound
		default:
			log.Error("failed to get join requests by club invite id", logger.Err(err))
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

	if !event.IsOwner(userId) {
		return eventservice.ErrPermissionsDenied
	}

	err = s.inviteDeleter.DeleteInvite(ctx, inviteId)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidID) {
			return eventservice.ErrInvalidID
		}
		log.Error("failed to delete invite", logger.Err(err))
		return err
	}

	return nil
}
