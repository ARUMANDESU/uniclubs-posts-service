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
	"time"
)

func (s Service) SendJoinRequestToClub(ctx context.Context, dto *dto.SendJoinRequestToClub) (*domain.Event, error) {
	const op = "collaborator.Service.sendJoinRequestToClub"
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

	if !event.IsOwner(dto.UserId) {
		return nil, eventService.ErrPermissionsDenied
	}

	if event.IsCollaborator(dto.Club.ID) {
		return nil, eventService.ErrClubAlreadyCollaborator
	}

	invite, err := s.clubInviteStorage.GetJoinRequestByClubId(ctx, dto.EventId, dto.Club.ID)
	if err != nil && !errors.Is(err, storage.ErrInviteNotFound) {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventService.ErrInvalidID
		default:
			log.Error("failed to get join request by club id", logger.Err(err))
			return nil, err
		}
	}

	if invite != nil {
		return nil, eventService.ErrInviteAlreadyExists
	}

	createCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	invite, err = s.clubInviteStorage.CreateJoinRequestToClub(createCtx, dto)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return nil, eventService.ErrInvalidID
		default:
			log.Error("failed to get join request by club id", logger.Err(err))
			return nil, err
		}
	}

	return event, nil
}

func (s Service) AcceptClubJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error) {
	// todo: implement me
	panic("implement me")
}

func (s Service) RejectClubJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error) {
	// todo: implement me
	panic("implement me")
}

func (s Service) KickClub(ctx context.Context, userId, targetId int64) error {
	// todo: implement me
	panic("implement me")
}

func (s Service) RevokeInviteClub(ctx context.Context, inviteId string, userId int64) error {
	const op = "collaborator.Service.revokeInviteClub"
	log := s.log.With(slog.String("op", op))

	invite, err := s.clubInviteStorage.GetJoinRequestsByClubInviteId(ctx, inviteId)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidID):
			return eventService.ErrInvalidID
		case errors.Is(err, storage.ErrInviteNotFound):
			return eventService.ErrInviteNotFound
		default:
			log.Error("failed to get join requests by club invite id", logger.Err(err))
			return err
		}
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

	if !event.IsOwner(userId) {
		return eventService.ErrPermissionsDenied
	}

	err = s.inviteDeleter.DeleteInvite(ctx, inviteId)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidID) {
			return eventService.ErrInvalidID
		}
		log.Error("failed to delete invite", logger.Err(err))
		return err
	}

	return nil
}
