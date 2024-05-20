package eventgrpc

import (
	"context"
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InfoService interface {
	GetEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error)
}

func (s serverApi) GetEvent(ctx context.Context, req *eventv1.GetEventRequest) (*eventv1.EventObject, error) {
	err := validate.GetEvent(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.info.GetEvent(ctx, req.GetEventId(), req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) ListEvents(ctx context.Context, req *eventv1.ListEventsRequest) (*eventv1.ListEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}
