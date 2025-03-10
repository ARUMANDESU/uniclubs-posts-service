package userservice

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

var ErrUserNotExist = errors.New("user does not exist")

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
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = s.usrStorage.UpdateUser(ctx, &user)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotExists):
			return ErrUserNotExist
		default:
			log.Error("failed to update user", logger.Err(err))
			return err
		}
	}

	return nil
}
