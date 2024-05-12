package event

import (
	"context"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ManagementService interface {
	CreateEvent(ctx context.Context, dto *dto.CreateEventDTO) (*domain.Event, error)
	GetEvent(ctx context.Context, dto *dto.GetEventDTO) (*domain.Event, error)
}

func (s serverApi) CreateEvent(ctx context.Context, req *eventv1.CreateEventRequest) (*eventv1.EventObject, error) {
	err := validate.CreateEvent(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.CreateEvent(ctx, dto.ProtoToCreateEventDTO(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return event.ToProto(), nil
}

func (s serverApi) GetEvent(ctx context.Context, req *eventv1.GetEventRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) UpdateEvent(ctx context.Context, req *eventv1.UpdateEventRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}
