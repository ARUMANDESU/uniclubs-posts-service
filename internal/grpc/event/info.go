package eventgrpc

import (
	"context"
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InfoService interface {
	GetEvent(ctx context.Context, eventId string, userId int64) (*dtos.GetEvent, error)
	ListEvents(ctx context.Context, filters domain.EventsFilter) ([]domain.Event, *domain.PaginationMetadata, error)
	GetUserInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.UserInvite, error)
	GetClubInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.Invite, error)
	ListParticipants(ctx context.Context, dto *dtos.ListParticipants) ([]domain.Participant, *domain.PaginationMetadata, error)
}

func (s serverApi) GetEvent(ctx context.Context, req *eventv1.GetEventRequest) (*eventv1.GetEventResponse, error) {
	err := validate.GetEvent(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dto, err := s.info.GetEvent(ctx, req.GetEventId(), req.GetUserId())
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

	return &eventv1.GetEventResponse{
		Event:             dto.Event.ToProto(),
		UserStatus:        eventv1.UserStatus(dto.UserStatus),
		ParticipantStatus: eventv1.ParticipantStatus(dto.ParticipantStatus),
	}, nil
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
	err := validate.ListParticipants(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dto := dtos.ProtoToListParticipants(request)
	participants, pagination, err := s.info.ListParticipants(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrParticipantNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &eventv1.ListParticipantsResponse{
		Participants: domain.ParticipantsToProto(participants),
		Metadata:     pagination.ToProto(),
	}, nil
}

func (s serverApi) GetClubInvites(ctx context.Context, req *eventv1.GetInvitesRequest) (*eventv1.GetClubInvitesResponse, error) {
	err := validate.GetInvitesRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dto := dtos.ProtoToGetInvites(req)
	invites, err := s.info.GetClubInvites(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrInviteNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &eventv1.GetClubInvitesResponse{
		Invites: domain.ClubInvitesToProto(invites),
	}, nil
}

func (s serverApi) GetOrganizerInvites(ctx context.Context, req *eventv1.GetInvitesRequest) (*eventv1.GetOrganizerInvitesResponse, error) {
	err := validate.GetInvitesRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dto := dtos.ProtoToGetInvites(req)
	invites, err := s.info.GetUserInvites(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrInviteNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &eventv1.GetOrganizerInvitesResponse{
		Invites: domain.UserInvitesToProto(invites),
	}, nil
}
