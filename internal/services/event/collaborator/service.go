package collaborator

import (
	"context"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"log/slog"
)

type Service struct {
	log                    *slog.Logger
	eventStorage           EventStorage
	organizerInviteStorage OrganizerInviteStorage
}

type EventStorage interface {
	GetEvent(ctx context.Context, eventId string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
}

type OrganizerInviteStorage interface {
	SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.UserInvite, error)
	GetJoinRequests(ctx context.Context, eventId string) ([]domain.UserInvite, error)
	GetJoinRequestByUserId(ctx context.Context, userId int64) (*domain.UserInvite, error)
	GetJoinRequestsById(ctx context.Context, requestId string) (*domain.UserInvite, error)
	DeleteJoinRequest(ctx context.Context, requestId string) error
}

func New(log *slog.Logger, eventProvider EventStorage, organizerInviteStorage OrganizerInviteStorage) Service {
	return Service{
		log:                    log,
		eventStorage:           eventProvider,
		organizerInviteStorage: organizerInviteStorage,
	}
}
