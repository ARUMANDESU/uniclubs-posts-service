package eventgrpc

import (
	"context"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ParticipantService interface {
	ParticipateEvent(ctx context.Context, eventId string, userId int64) (*eventv1.EventObject, error)
	CancelParticipation(ctx context.Context, eventId string, userId int64) (*eventv1.EventObject, error)
	KickParticipant(ctx context.Context, dto *dtos.KickParticipant) error
	BanParticipant(ctx context.Context, dto *dtos.BanParticipant) (*eventv1.EventObject, error)
	UnbanParticipant(ctx context.Context, dto *dtos.UnbanParticipant) (*eventv1.EventObject, error)
}

func (s serverApi) ParticipateEvent(ctx context.Context, req *eventv1.EventActionRequest) (*emptypb.Empty, error) {
	err := validate.EventActionRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.participant.ParticipateEvent(ctx, req.EventId, req.UserId)
	if err != nil {
		return nil, handleError(err)
	}

	return nil, nil
}

func (s serverApi) CancelParticipation(ctx context.Context, req *eventv1.EventActionRequest) (*emptypb.Empty, error) {
	err := validate.EventActionRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.participant.CancelParticipation(ctx, req.EventId, req.UserId)
	if err != nil {
		return nil, handleError(err)
	}

	return nil, nil
}

func (s serverApi) KickParticipant(ctx context.Context, req *eventv1.KickParticipantRequest) (*emptypb.Empty, error) {
	err := validate.KickParticipantRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.participant.KickParticipant(ctx, dtos.ProtoToKickParticipant(req))
	if err != nil {
		return nil, handleError(err)
	}

	return nil, nil
}

func (s serverApi) BanParticipant(ctx context.Context, req *eventv1.BanParticipantRequest) (*emptypb.Empty, error) {
	err := validate.BanParticipantRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.participant.BanParticipant(ctx, dtos.ProtoToBanParticipant(req))
	if err != nil {
		return nil, handleError(err)
	}

	return nil, nil
}

func (s serverApi) UnbanParticipant(ctx context.Context, request *eventv1.UnbanParticipantRequest) (*emptypb.Empty, error) {
	err := validate.UnbanParticipantRequest(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.participant.UnbanParticipant(ctx, dtos.ProtoToUnbanParticipant(request))
	if err != nil {
		return nil, handleError(err)
	}

	return nil, nil
}
