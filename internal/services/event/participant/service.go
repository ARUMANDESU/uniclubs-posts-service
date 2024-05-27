package eventparticipant

import (
	"context"
	"errors"
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	eventservice "github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"time"
)

type Service struct {
	log *slog.Logger
	Storage
}
type EventStorage interface {
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
}

type ParticipantStorage interface {
	GetEventParticipant(ctx context.Context, eventId string, userId int64) (*domain.Participant, error)
	AddEventParticipant(ctx context.Context, participant *domain.Participant) error
	DeleteEventParticipant(ctx context.Context, eventId string, userId int64) error
}

type BanStorage interface {
	GetBanRecord(ctx context.Context, eventId string, userId int64) (*domain.BanRecord, error)
	BanParticipant(ctx context.Context, dto *dtos.BanParticipant) error
	UnBanParticipant(ctx context.Context, eventId string, userId int64) error
}

type UserProvider interface {
	GetUserById(ctx context.Context, id int64) (*domain.User, error)
}

type ClubProvider interface {
	IsClubMember(ctx context.Context, userId, clubId int64) (bool, error)
	IsBanned(ctx context.Context, userId, clubId int64) (bool, error)
}

func New(log *slog.Logger, storage Storage) Service {
	return Service{
		log:     log,
		Storage: storage,
	}
}

func (s Service) ParticipateEvent(ctx context.Context, eventId string, userId int64) (*eventv1.EventObject, error) {
	const op = "service.event.participant.participateEvent"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, eventId)
	if err != nil {
		return nil, s.handleError("failed to get event", log, err)
	}

	if event.Status != domain.EventStatusInProgress {
		return nil, fmt.Errorf("%w: can't participate in event that is not in progress", eventservice.ErrInvalidEventStatus)
	}

	if event.MaxParticipants != 0 && event.ParticipantsCount >= event.MaxParticipants {
		return nil, fmt.Errorf("can't participate: %w", eventservice.ErrEventIsFull)
	}

	record, err := s.banStorage.GetBanRecord(ctx, eventId, userId)
	if err != nil && !errors.Is(err, storage.ErrBanRecordNotFound) {
		return nil, s.handleError("failed to get ban record", log, err)
	}
	if record != nil {
		return nil, fmt.Errorf("can't participate: %w", eventservice.ErrUserIsBanned)
	}

	banned, err := s.clubProvider.IsBanned(ctx, userId, event.ClubId)
	if err != nil && !errors.Is(err, storage.ErrBanRecordNotFound) {
		return nil, s.handleError("failed to check if user is banned ", log, err)
	}
	if banned {
		return nil, fmt.Errorf("can't participate: %w", eventservice.ErrUserIsBanned)
	}

	participant, err := s.participantStorage.GetEventParticipant(ctx, eventId, userId)
	if err != nil && !errors.Is(err, storage.ErrParticipantNotFound) {
		return nil, s.handleError("failed to get participant", log, err)
	}

	if participant != nil {
		return nil, fmt.Errorf("can't participate: %w", eventservice.ErrAlreadyParticipating)
	}

	if event.Type == domain.EventTypeIntraClub {
		isMemberOfCollabClubs, err := s.IsMemberOfCollabClubs(ctx, event, userId)
		if err != nil {
			return nil, fmt.Errorf("can't participate: %w, %w", err, eventservice.ErrUserIsFromAnotherClub)
		}
		if !isMemberOfCollabClubs {
			return nil, fmt.Errorf("can't participate: %w", eventservice.ErrUserIsFromAnotherClub)
		}
	}

	user, err := s.userProvider.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	participant = &domain.Participant{
		ID:       primitive.NewObjectID().Hex(),
		EventId:  eventId,
		User:     *user,
		JoinedAt: time.Now(),
	}

	err = s.participantStorage.AddEventParticipant(ctx, participant)
	if err != nil {
		return nil, s.handleError("failed to add participant", log, err)
	}

	event.ParticipantsCount++

	event, err = s.eventStorage.UpdateEvent(ctx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return event.ToProto(), nil
}

func (s Service) CancelParticipation(ctx context.Context, eventId string, userId int64) (*eventv1.EventObject, error) {
	const op = "service.event.participant.cancelParticipation"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, eventId)
	if err != nil {
		return nil, s.handleError("failed to get event", log, err)
	}

	_, err = s.participantStorage.GetEventParticipant(ctx, eventId, userId)
	if err != nil {
		return nil, s.handleError("failed to get participant", log, err)
	}

	err = s.participantStorage.DeleteEventParticipant(ctx, eventId, userId)
	if err != nil {
		return nil, s.handleError("failed to delete participant", log, err)
	}

	if event.ParticipantsCount > 0 {
		event.ParticipantsCount--
	} else {
		event.ParticipantsCount = 0
	}

	event, err = s.eventStorage.UpdateEvent(ctx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return event.ToProto(), nil
}

func (s Service) KickParticipant(ctx context.Context, dto *dtos.KickParticipant) error {
	const op = "service.event.participant.kickParticipant"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, dto.EventId)
	if err != nil {
		return s.handleError("failed to get event", log, err)
	}

	if !event.IsOrganizer(dto.UserId) {
		return eventservice.ErrPermissionsDenied
	}

	if event.IsOrganizer(dto.ParticipantId) {
		return fmt.Errorf("%w: can't kick organizer from event", eventservice.ErrUserIsNotEventOrganizer)
	}

	if event.Status == domain.EventStatusFinished || event.Status == domain.EventStatusCanceled || event.Status == domain.EventStatusArchived {
		return fmt.Errorf("%w: can't kick participant from event, which status is %s", eventservice.ErrInvalidEventStatus, event.Status)
	}

	_, err = s.participantStorage.GetEventParticipant(ctx, dto.EventId, dto.ParticipantId)
	if err != nil {
		return s.handleError("failed to get participant", log, err)
	}

	err = s.participantStorage.DeleteEventParticipant(ctx, dto.EventId, dto.ParticipantId)
	if err != nil {
		return s.handleError("failed to delete participant", log, err)
	}

	if event.ParticipantsCount > 0 {
		event.ParticipantsCount--
	} else {
		event.ParticipantsCount = 0
	}

	event, err = s.eventStorage.UpdateEvent(ctx, event)
	if err != nil {
		return s.handleError("failed to update event", log, err)
	}

	return nil
}

func (s Service) BanParticipant(ctx context.Context, dto *dtos.BanParticipant) (*eventv1.EventObject, error) {
	const op = "service.event.participant.banParticipant"
	log := s.log.With(slog.String("op", op))

	event, err := s.eventStorage.GetEvent(ctx, dto.EventId)
	if err != nil {
		return nil, s.handleError("failed to get event", log, err)
	}

	if !event.IsOrganizer(dto.UserId) {
		return nil, eventservice.ErrPermissionsDenied
	}

	if event.IsOrganizer(dto.ParticipantId) {
		return nil, fmt.Errorf("%w: can't ban organizer from event", eventservice.ErrUserIsNotEventOrganizer)
	}

	if event.Status == domain.EventStatusFinished || event.Status == domain.EventStatusCanceled || event.Status == domain.EventStatusArchived {
		return nil, fmt.Errorf("%w: can't ban participant from event, which status is %s", eventservice.ErrInvalidEventStatus, event.Status)
	}

	record, err := s.banStorage.GetBanRecord(ctx, dto.EventId, dto.ParticipantId)
	if err != nil && !errors.Is(err, storage.ErrBanRecordNotFound) {
		return nil, err
	}
	if record != nil {
		return nil, eventservice.ErrUserAlreadyBanned
	}

	banned, err := s.clubProvider.IsBanned(ctx, dto.ParticipantId, event.ClubId)
	if err != nil {
		return nil, s.handleError("failed to check if user is banned ", log, err)
	}
	if banned {
		return nil, eventservice.ErrUserAlreadyBanned
	}

	err = s.banStorage.BanParticipant(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to ban participant: %w", err)
	}

	participant, err := s.participantStorage.GetEventParticipant(ctx, dto.EventId, dto.ParticipantId)
	if err != nil && !errors.Is(err, storage.ErrParticipantNotFound) {
		return nil, s.handleError("failed to get participant", log, err)
	}
	if participant == nil {
		return event.ToProto(), nil
	}

	err = s.participantStorage.DeleteEventParticipant(ctx, dto.EventId, dto.ParticipantId)
	if err != nil {
		return nil, s.handleError("failed to delete participant", log, err)
	}

	if event.ParticipantsCount > 0 {
		event.ParticipantsCount--
	} else {
		event.ParticipantsCount = 0
	}

	event, err = s.eventStorage.UpdateEvent(ctx, event)
	if err != nil {
		return nil, s.handleError("failed to update event", log, err)
	}

	return event.ToProto(), nil
}

func (s Service) IsMemberOfCollabClubs(ctx context.Context, event *domain.Event, userId int64) (bool, error) {
	const op = "service.event.participant.checkIsMemberOfCollaboratorClub"
	log := s.log.With(slog.String("op", op))

	var err []error

	for _, club := range event.CollaboratorClubs {
		isMember, e := s.clubProvider.IsClubMember(ctx, userId, club.ID)
		if e != nil {
			log.Warn("failed to check if user is a member of the club", logger.Err(e))
			err = append(err, e)
			continue
		}

		if isMember {
			return true, nil
		}
	}

	if len(err) > 0 {
		return false, err[0]
	}

	return false, nil
}

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
	case errors.Is(err, storage.ErrParticipantNotFound):
		return eventservice.ErrParticipantNotFound
	case errors.Is(err, storage.ErrBanRecordNotFound):
		return eventservice.ErrBanRecordNotFound
	default:
		log.Error(msg, logger.Err(err))
		return err
	}
}

type Storage struct {
	eventStorage       EventStorage
	userProvider       UserProvider
	clubProvider       ClubProvider
	participantStorage ParticipantStorage
	banStorage         BanStorage
}

func NewStorage(
	eventStorage EventStorage,
	userProvider UserProvider,
	clubProvider ClubProvider,
	participantStorage ParticipantStorage,
	banStorage BanStorage,
) Storage {
	return Storage{
		eventStorage:       eventStorage,
		userProvider:       userProvider,
		clubProvider:       clubProvider,
		participantStorage: participantStorage,
		banStorage:         banStorage,
	}
}
