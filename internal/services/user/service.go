package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotExist = errors.New("user does not exist")
)

type Service struct {
	log        *slog.Logger
	usrStorage Storage
}

type Storage interface {
	UpdateUser(ctx context.Context, user *domain.User) error
}

func New(log *slog.Logger, storage Storage) *Service {
	return &Service{
		log:        log,
		usrStorage: storage,
	}
}

func (s Service) HandleUpdateUser(msg amqp091.Delivery) error {
	const op = "services.user.handleUpdateUser"

	log := s.log.With(slog.String("op", op))

	var user domain.User

	err := json.Unmarshal(msg.Body, &user)
	if err != nil {
		log.Error("failed to unmarshal message", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = s.usrStorage.UpdateUser(ctx, &user)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotExists):
			log.Error("user not found", logger.Err(err))
			return fmt.Errorf("%s: %w", op, ErrUserNotExist)
		default:
			log.Error("failed to update user", logger.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
