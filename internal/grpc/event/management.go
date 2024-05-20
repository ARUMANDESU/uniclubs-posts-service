package eventgrpc

import (
	"context"
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ManagementService interface {
	CreateEvent(ctx context.Context, club domain.Club, user domain.User) (*domain.Event, error)
	UpdateEvent(ctx context.Context, dto *dto.UpdateEvent) (*domain.Event, error)
	DeleteEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error)
}

func (s serverApi) CreateEvent(ctx context.Context, req *eventv1.CreateEventRequest) (*eventv1.EventObject, error) {
	err := validate.CreateEvent(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.CreateEvent(ctx, domain.ClubFromProto(req.GetClub()), domain.UserFromProto(req.GetUser()))
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) UpdateEvent(ctx context.Context, req *eventv1.UpdateEventRequest) (*eventv1.EventObject, error) {
	err := validate.UpdateEvent(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.UpdateEvent(ctx, dto.UpdateToDTO(req))
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrEventUpdateConflict):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, eventservice.ErrUserIsNotEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) DeleteEvent(ctx context.Context, req *eventv1.DeleteEventRequest) (*eventv1.EventObject, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.DeleteEvent(ctx, req.GetEventId(), req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrUserIsNotEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) PublishEvent(ctx context.Context, req *eventv1.PublishEventRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) UnpublishEvent(ctx context.Context, req *eventv1.PublishEventRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}
