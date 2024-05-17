package event

import (
	"context"
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CollaboratorService interface {
	SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.Event, error)
}

func (s serverApi) AddCollaborator(ctx context.Context, req *eventv1.AddCollaboratorRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) RemoveCollaborator(ctx context.Context, req *eventv1.RemoveCollaboratorRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) AddOrganizer(ctx context.Context, req *eventv1.AddOrganizerRequest) (*empty.Empty, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0), validation.NotIn(req.OrganizerId).Error("organizer_id must be different from user_id")),
		validation.Field(&req.OrganizerId, validation.Required, validation.Min(0)),
		validation.Field(&req.OrganizerClubId, validation.Required, validation.Min(0)),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.collaborator.SendJoinRequestToUser(ctx, dto.AddOrganizerRequestToUserToDTO(req))
	if err != nil {
		switch {
		case errors.Is(err, event.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, event.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, event.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, event.ErrInviteAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, event.ErrUserIsFromAnotherClub):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}

func (s serverApi) RemoveOrganizer(ctx context.Context, req *eventv1.RemoveOrganizerRequest) (*empty.Empty, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0), validation.NotIn(req.OrganizerId).Error("organizer_id must be different from user_id")),
		validation.Field(&req.OrganizerId, validation.Required, validation.Min(0)),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//TODO implement me
	panic("implement me")
}

func (s serverApi) HandleInviteUser(ctx context.Context, request *eventv1.HandleInviteUserRequest) (*emptypb.Empty, error) {
	err := validation.ValidateStruct(request,
		validation.Field(&request.UserId, validation.Required, validation.Min(0)),
		validation.Field(&request.InviteId, validation.Required, validation.Min(0)),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//TODO implement me
	panic("implement me")
}

func (s serverApi) RevokeInviteUser(ctx context.Context, request *eventv1.RevokeInviteUserRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
