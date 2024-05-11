package club

import (
	"context"
	"encoding/json"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

type Service struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	SaveClub(ctx context.Context, club *domain.Club) error
	UpdateClub(ctx context.Context, club *domain.Club) error
	GetClubByID(ctx context.Context, clubID int64) (club *domain.Club, err error)
	DeleteClub(ctx context.Context, club *domain.Club) error
}

func New(log *slog.Logger, storage Storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s Service) HandleCreateClub(msg amqp091.Delivery) error {
	const op = "services.club.handleCreateClub"
	log := s.log.With(slog.String("op", op))
	var input *domain.Club

	err := json.Unmarshal(msg.Body, &input)
	if err != nil {
		log.Error("failed to unmarshal message", logger.Err(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.storage.SaveClub(ctx, input)
	if err != nil {
		log.Error("failed to save club", logger.Err(err))
		return err
	}

	return nil
}

func (s Service) HandleUpdateClub(msg amqp091.Delivery) error {
	const op = "services.club.handleUpdateClub"
	log := s.log.With(slog.String("op", op))

	var input struct {
		ID      int64   `json:"id"`
		Name    *string `json:"name" bson:"name"`
		LogoURL *string `json:"logo_url" bson:"logo_url"`
	}

	err := json.Unmarshal(msg.Body, &input)
	if err != nil {
		log.Error("failed to unmarshal message", logger.Err(err))
		return err
	}

	getClubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	club, err := s.storage.GetClubByID(getClubCtx, input.ID)
	if err != nil {
		log.Error("failed to get club", logger.Err(err))
		return err
	}

	if input.Name != nil {
		club.Name = *input.Name
	}
	if input.LogoURL != nil {
		club.LogoURL = *input.LogoURL
	}

	updateClubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.storage.UpdateClub(updateClubCtx, club)
	if err != nil {
		log.Error("failed to update club", logger.Err(err))
		return err
	}

	return nil
}

func (s Service) HandleDeleteClub(msg amqp091.Delivery) error {
	const op = "services.club.handleDeleteClub"
	log := s.log.With(slog.String("op", op))

	var clubID int64

	err := json.Unmarshal(msg.Body, &clubID)
	if err != nil {
		log.Error("failed to unmarshal message", logger.Err(err))
		return err
	}

	getClubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	club, err := s.storage.GetClubByID(getClubCtx, clubID)
	if err != nil {
		log.Error("failed to get club", logger.Err(err))
		return err
	}

	deleteClubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.storage.DeleteClub(deleteClubCtx, club)
	if err != nil {
		log.Error("failed to delete club", logger.Err(err))
		return err
	}

	return nil
}
