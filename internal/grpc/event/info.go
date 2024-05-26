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
	ListEvents(ctx context.Context, filters domain.Filters) ([]domain.Event, *domain.PaginationMetadata, error)
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
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) ListEvents(ctx context.Context, req *eventv1.ListEventsRequest) (*eventv1.ListEventsResponse, error) {
	err := validate.ListEvents(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filters := domain.ProtoToFilers(req)
	events, pagination, err := s.info.ListEvents(ctx, filters)
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &eventv1.ListEventsResponse{
		Events:   domain.EventsToProto(events),
		Metadata: pagination.ToProto(),
	}, nil
}

func (s serverApi) ListParticipatedEvents(ctx context.Context, request *eventv1.ListParticipatedEventsRequest) (*eventv1.ListEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) ListParticipants(ctx context.Context, request *eventv1.ListParticipantsRequest) (*eventv1.ListParticipantsResponse, error) {
	//TODO implement me
	panic("implement me")
}
