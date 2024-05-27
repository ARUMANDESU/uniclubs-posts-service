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

type (
	ManagementService interface {
		CreateEvent(ctx context.Context, club domain.Club, user domain.User) (*domain.Event, error)
		UpdateEvent(ctx context.Context, dto *dtos.UpdateEvent) (*domain.Event, error)
		DeleteEvent(ctx context.Context, dto *dtos.DeleteEvent) (*domain.Event, error)

		PublishEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error)
		UnpublishEvent(ctx context.Context, eventId string, userId int64) (*domain.Event, error)
		SendToReview(ctx context.Context, eventId string, userId int64) (*domain.Event, error)
		RevokeReview(ctx context.Context, eventId string, userId int64) (*domain.Event, error)
		ApproveEvent(ctx context.Context, eventId string, user domain.User) (*domain.Event, error)
		RejectEvent(ctx context.Context, dto *dtos.RejectEvent) (*domain.Event, error)
	}
)

func (s serverApi) CreateEvent(ctx context.Context, req *eventv1.CreateEventRequest) (*eventv1.EventObject, error) {
	if err := validate.CreateEvent(req); err != nil {
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
	if err := validate.UpdateEvent(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dto, err := dtos.UpdateToDTO(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.UpdateEvent(ctx, dto)
	if err != nil {
		return nil, handleError(err)
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

	event, err := s.management.DeleteEvent(ctx, dtos.DeleteEventToDTO(req))
	if err != nil {
		return nil, handleError(err)
	}

	return event.ToProto(), nil
}

func (s serverApi) PublishEvent(ctx context.Context, req *eventv1.EventActionRequest) (*eventv1.EventObject, error) {
	if err := validate.EventActionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.PublishEvent(ctx, req.GetEventId(), req.GetUserId())
	if err != nil {
		return nil, handleError(err)
	}

	return event.ToProto(), nil
}

func (s serverApi) UnpublishEvent(ctx context.Context, req *eventv1.EventActionRequest) (*eventv1.EventObject, error) {
	if err := validate.EventActionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.UnpublishEvent(ctx, req.GetEventId(), req.GetUserId())
	if err != nil {
		return nil, handleError(err)
	}

	return event.ToProto(), nil
}

func (s serverApi) ApproveEvent(ctx context.Context, req *eventv1.ApproveEventRequest) (*eventv1.EventObject, error) {
	if err := validate.ApproveEventRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.ApproveEvent(ctx, req.GetEventId(), domain.UserFromProto(req.GetUser()))
	if err != nil {
		return nil, handleError(err)
	}

	return event.ToProto(), nil
}

func (s serverApi) RejectEvent(ctx context.Context, req *eventv1.RejectEventRequest) (*eventv1.EventObject, error) {
	if err := validate.RejectEventRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.RejectEvent(ctx, dtos.RejectEventToDTO(req))
	if err != nil {
		return nil, handleError(err)
	}

	return event.ToProto(), nil
}

func (s serverApi) SendToReview(ctx context.Context, req *eventv1.EventActionRequest) (*eventv1.EventObject, error) {

	if err := validate.EventActionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.SendToReview(ctx, req.GetEventId(), req.GetUserId())
	if err != nil {
		return nil, handleError(err)
	}

	return event.ToProto(), nil
}

func (s serverApi) RevokeReview(ctx context.Context, req *eventv1.EventActionRequest) (*eventv1.EventObject, error) {
	if err := validate.EventActionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.management.RevokeReview(ctx, req.GetEventId(), req.GetUserId())
	if err != nil {
		return nil, handleError(err)
	}

	return event.ToProto(), nil
}

func handleError(err error) error {
	switch {
	case errors.Is(err, eventservice.ErrEventNotFound),
		errors.Is(err, eventservice.ErrClubNotExists),
		errors.Is(err, eventservice.ErrParticipantNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, eventservice.ErrInvalidID),
		errors.Is(err, eventservice.ErrEventInvalidFields):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, eventservice.ErrUserIsNotEventOwner),
		errors.Is(err, eventservice.ErrUserIsFromAnotherClub),
		errors.Is(err, eventservice.ErrUserIsNotEventOrganizer),
		errors.Is(err, eventservice.ErrPermissionsDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, eventservice.ErrEventUpdateConflict):
		return status.Error(codes.Aborted, err.Error())
	case errors.Is(err, eventservice.ErrEventIsFull),
		errors.Is(err, eventservice.ErrAlreadyParticipating),
		errors.Is(err, eventservice.ErrInvalidEventStatus):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
