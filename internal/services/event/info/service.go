package eventinfo

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
)

type Service struct {
	log *slog.Logger
	Storage
}

type EventProvider interface {
	GetEvent(ctx context.Context, eventId string) (*domain.Event, error)
	ListEvents(ctx context.Context, filters domain.Filters) ([]domain.Event, *domain.PaginationMetadata, error)
}

type ParticipantProvider interface {
	GetEventParticipant(ctx context.Context, eventId string, userId int64) (*domain.Participant, error)
}

type BanProvider interface {
	GetBanRecord(ctx context.Context, eventId string, userId int64) (*domain.BanRecord, error)
}

type ClubProvider interface {
	IsBanned(ctx context.Context, userId int64, clubId int64) (bool, error)
}

type InviteProvider interface {
	GetUserInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.UserInvite, error)
	GetClubInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.Invite, error)
}

func New(log *slog.Logger, storage Storage) Service {
	return Service{
		log:     log,
		Storage: storage,
	}
}

func (s Service) GetEvent(ctx context.Context, eventId string, userId int64) (*dtos.GetEvent, error) {
	const op = "services.event.management.getEvent"
	log := s.log.With(slog.String("op", op))

	participantStatus := domain.ParticipantStatusUnknown
	userStatus := domain.UserStatusUnknown

	event, err := s.eventProvider.GetEvent(ctx, eventId)
	if err != nil {
		return nil, s.handleError("failed to get event", log, err)
	}

	if !event.IsOrganizer(userId) && event.Status == domain.EventStatusDraft {
		return nil, eventservice.ErrEventNotFound
	}
	if event.IsOrganizer(userId) {
		userStatus = domain.UserStatusOrganizer
	}
	if event.IsOwner(userId) {
		userStatus = domain.UserStatusOwner
	}

	if userId != 0 {
		participantStatus, err = s.getParticipantStatus(ctx, event, userId)
		if err != nil {
			return nil, s.handleError("failed to get ban status", log, err)
		}
	}

	return &dtos.GetEvent{
		Event:             *event,
		UserStatus:        userStatus,
		ParticipantStatus: participantStatus,
	}, nil
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

func (s Service) GetUserInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.UserInvite, error) {
	const op = "services.event.management.getUserInvites"
	log := s.log.With(slog.String("op", op))

	invites, err := s.inviteProvider.GetUserInvites(ctx, dto)
	if err != nil {
		return nil, s.handleError("failed to get user invites", log, err)
	}

	return invites, nil
}

func (s Service) GetClubInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.Invite, error) {
	const op = "services.event.management.getClubInvites"
	log := s.log.With(slog.String("op", op))

	invites, err := s.inviteProvider.GetClubInvites(ctx, dto)
	if err != nil {
		return nil, s.handleError("failed to get club invites", log, err)
	}

	return invites, nil
}

// handleError handles common errors
func (s Service) handleError(msg string, log *slog.Logger, err error) error {
	switch {
	case errors.Is(err, storage.ErrEventNotFound):
		return eventservice.ErrEventNotFound
	case errors.Is(err, storage.ErrInvalidID):
		return eventservice.ErrInvalidID
	case errors.Is(err, storage.ErrParticipantNotFound):
		return eventservice.ErrParticipantNotFound
	case errors.Is(err, storage.ErrBanRecordNotFound):
		return eventservice.ErrBanRecordNotFound
	case errors.Is(err, storage.ErrInviteNotFound):
		return eventservice.ErrInviteNotFound
	default:
		log.Error(msg, logger.Err(err))
		return err
	}
}

func (s Service) getParticipantStatus(ctx context.Context, event *domain.Event, userId int64) (domain.ParticipantStatus, error) {
	const op = "services.event.management.getUserBanStatus"
	log := s.log.With(slog.String("op", op))

	participant, err := s.participantProvider.GetEventParticipant(ctx, event.ID, userId)
	if err != nil && !errors.Is(err, storage.ErrParticipantNotFound) {
		return domain.ParticipantStatusUnknown, s.handleError("failed to get participant", log, err)
	}

	if participant != nil {
		return domain.ParticipantStatusJoined, nil
	}

	banRecord, err := s.banProvider.GetBanRecord(ctx, event.ID, userId)
	if err != nil && !errors.Is(err, storage.ErrBanRecordNotFound) {
		return domain.ParticipantStatusUnknown, s.handleError("failed to get ban record", log, err)
	}
	if banRecord != nil {
		return domain.ParticipantStatusBanned, nil
	}

	banned, err := s.clubProvider.IsBanned(ctx, event.ClubId, userId)
	if err != nil && !errors.Is(err, storage.ErrBanRecordNotFound) {
		return domain.ParticipantStatusUnknown, err
	}
	if banned {
		return domain.ParticipantStatusBanned, nil
	}

	return domain.ParticipantStatusUnknown, nil
}

type Storage struct {
	eventProvider       EventProvider
	participantProvider ParticipantProvider
	banProvider         BanProvider
	clubProvider        ClubProvider
	inviteProvider      InviteProvider
}

func NewStorage(
	eventProvider EventProvider,
	participantProvider ParticipantProvider,
	banProvider BanProvider,
	clubProvider ClubProvider,
	inviteProvider InviteProvider,
) Storage {
	return Storage{
		eventProvider:       eventProvider,
		participantProvider: participantProvider,
		banProvider:         banProvider,
		clubProvider:        clubProvider,
		inviteProvider:      inviteProvider,
	}
}
