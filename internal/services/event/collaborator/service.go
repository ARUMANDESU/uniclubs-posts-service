package eventcollab

import (
	"context"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"log/slog"
)

type Service struct {
	log               *slog.Logger
	eventStorage      EventStorage
	userInviteStorage OrganizerInviteStorage
	clubInviteStorage ClubInviteStorage
	inviteDeleter     InviteDeleter
}

type EventStorage interface {
	GetEvent(ctx context.Context, eventId string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
}

type OrganizerInviteStorage interface {
	CreateJoinRequestToUser(ctx context.Context, dto *dtos.SendJoinRequestToUser) (*domain.UserInvite, error)
	GetUserJoinRequests(ctx context.Context, eventId string) ([]domain.UserInvite, error)
	GetJoinRequestsByUserInviteId(ctx context.Context, inviteId string) (*domain.UserInvite, error)
	GetJoinRequestByUserId(ctx context.Context, eventId string, userId int64) (*domain.UserInvite, error)
}

type ClubInviteStorage interface {
	CreateJoinRequestToClub(ctx context.Context, dto *dtos.SendJoinRequestToClub) (*domain.Invite, error)
	GetClubJoinRequests(ctx context.Context, eventId string) ([]domain.Invite, error)
	GetJoinRequestsByClubInviteId(ctx context.Context, inviteId string) (*domain.Invite, error)
	GetJoinRequestByClubId(ctx context.Context, eventId string, clubId int64) (*domain.Invite, error)
}

type InviteDeleter interface {
	DeleteInvite(ctx context.Context, inviteId string) error
}

func New(
	log *slog.Logger,
	eventProvider EventStorage,
	organizerInviteStorage OrganizerInviteStorage,
	clubInviteStorage ClubInviteStorage,
	inviteDeleter InviteDeleter,
) Service {
	return Service{
		log:               log,
		eventStorage:      eventProvider,
		userInviteStorage: organizerInviteStorage,
		clubInviteStorage: clubInviteStorage,
		inviteDeleter:     inviteDeleter,
	}
}
