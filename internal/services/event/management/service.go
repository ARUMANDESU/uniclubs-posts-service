package management

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"log/slog"
	"time"

	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
)

var (
	ErrClubNotExists = errors.New("club not found")
	ErrEventNotFound = errors.New("event not found")
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
