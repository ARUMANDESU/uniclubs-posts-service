package collaborator

import (
	"context"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
)

func (s Service) SendJoinRequestToClub(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.Event, error) {
	// todo: implement me
	panic("implement me")
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
	// todo: implement me
	panic("implement me")
}
