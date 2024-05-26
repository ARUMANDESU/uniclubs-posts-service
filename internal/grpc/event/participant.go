package eventgrpc

import (
	"context"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ParticipantService interface {
	ParticipateEvent(ctx context.Context, eventId string, userId int64) (*eventv1.EventObject, error)
	CancelParticipation(ctx context.Context, eventId string, userId int64) (*eventv1.EventObject, error)
	KickParticipant(ctx context.Context, dto *dtos.KickParticipant) (*eventv1.EventObject, error)
	BanParticipant(ctx context.Context, dto *dtos.BanParticipant) (*eventv1.EventObject, error)
}

func (s serverApi) ParticipateEvent(ctx context.Context, req *eventv1.EventActionRequest) (*eventv1.EventObject, error) {
	err := validate.EventActionRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.participant.ParticipateEvent(ctx, req.EventId, req.UserId)
	if err != nil {
		return nil, handleError(err)
	}

	return event, nil
}

func (s serverApi) CancelParticipation(ctx context.Context, req *eventv1.EventActionRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) KickParticipant(ctx context.Context, req *eventv1.KickParticipantRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) BanParticipant(ctx context.Context, req *eventv1.BanParticipantRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}
