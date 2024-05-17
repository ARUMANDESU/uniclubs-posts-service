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

type Service struct {
	log                    *slog.Logger
	eventProvider          EventProvider
	organizerInviteStorage OrganizerInviteStorage
}

type EventProvider interface {
	GetEvent(ctx context.Context, eventId string) (*domain.Event, error)
}

type OrganizerInviteStorage interface {
	SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.UserInvite, error)
	GetJoinRequests(ctx context.Context, eventId string) ([]domain.UserInvite, error)
	GetJoinRequestByUserId(ctx context.Context, userId int64) (*domain.UserInvite, error)
	GetJoinRequestsById(ctx context.Context, requestId string) (*domain.UserInvite, error)
}

func New(log *slog.Logger, eventProvider EventProvider, organizerInviteStorage OrganizerInviteStorage) Service {
	return Service{
		log:                    log,
		eventProvider:          eventProvider,
		organizerInviteStorage: organizerInviteStorage,
	}
}

func (s Service) SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.Event, error) {
	const op = "services.event.management.sendJoinRequestToUser"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventProvider.GetEvent(ctx, dto.EventId)
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
	if event.IsOrganizer(dto.TargetId) {
		return nil, eventService.ErrUserAlreadyOrganizer
	}

	userInvite, err := s.organizerInviteStorage.GetJoinRequestByUserId(ctx, dto.TargetId)
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
	if event.User.ID == dto.UserId {
		if event.Club.ID != dto.TargetClubId {
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
