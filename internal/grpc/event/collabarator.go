package event

import (
	"context"
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrganizerService interface {
	SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.Event, error)
	AcceptUserJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error)
	RejectUserJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error)
	KickOrganizer(ctx context.Context, eventId string, userId, targetId int64) error
	RevokeInviteOrganizer(ctx context.Context, inviteId string, userId int64) error
}

type CollaboratorService interface {
	SendJoinRequestToClub(ctx context.Context, dto *dto.SendJoinRequestToClub) (*domain.Event, error)
	AcceptClubJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error)
	RejectClubJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error)
	KickClub(ctx context.Context, userId, targetId int64) error
	RevokeInviteClub(ctx context.Context, inviteId string, userId int64) error
}

func (s serverApi) AddCollaborator(ctx context.Context, req *eventv1.AddCollaboratorRequest) (*empty.Empty, error) {
	err := validate.AddCollaborator(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.collaborator.SendJoinRequestToClub(ctx, dto.AddCollaboratorRequestToClubToDTO(req))
	if err != nil {
		switch {
		case errors.Is(err, event.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, event.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, event.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, event.ErrInviteAlreadyExists), errors.Is(err, event.ErrClubAlreadyCollaborator):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}

func (s serverApi) RemoveCollaborator(ctx context.Context, req *eventv1.RemoveCollaboratorRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) HandleInviteClub(ctx context.Context, req *eventv1.HandleInviteClubRequest) (*eventv1.EventObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) RevokeInviteClub(ctx context.Context, request *eventv1.RevokeInviteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) AddOrganizer(ctx context.Context, req *eventv1.AddOrganizerRequest) (*empty.Empty, error) {
	err := validate.AddOrganizer(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.organizer.SendJoinRequestToUser(ctx, dto.AddOrganizerRequestToUserToDTO(req))
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
	err := validate.RemoveOrganizer(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.organizer.KickOrganizer(ctx, req.GetEventId(), req.GetUserId(), req.GetOrganizerId())
	if err != nil {
		switch {
		case errors.Is(err, event.ErrEventNotFound), errors.Is(err, event.ErrUserIsNotEventOrganizer):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, event.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, event.ErrPermissionsDenied), errors.Is(err, event.ErrUserIsEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, event.ErrEventUpdateConflict):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}

func (s serverApi) HandleInviteUser(ctx context.Context, req *eventv1.HandleInviteUserRequest) (*eventv1.EventObject, error) {
	err := validate.HandleInviteUser(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var res domain.Event
	if req.GetAction() == eventv1.HandleInvite_Action_ACCEPT {
		res, err = s.organizer.AcceptUserJoinRequest(ctx, req.InviteId, req.UserId)
	} else if req.GetAction() == eventv1.HandleInvite_Action_REJECT {
		res, err = s.organizer.RejectUserJoinRequest(ctx, req.InviteId, req.UserId)
	}

	if err != nil {
		switch {
		case errors.Is(err, event.ErrInviteNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, event.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, event.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, event.ErrUserAlreadyOrganizer):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, event.ErrEventUpdateConflict):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return res.ToProto(), nil
}

func (s serverApi) RevokeInviteUser(ctx context.Context, req *eventv1.RevokeInviteRequest) (*emptypb.Empty, error) {
	err := validate.RevokeInviteUser(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.organizer.RevokeInviteOrganizer(ctx, req.GetInviteId(), req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, event.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, event.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, event.ErrPermissionsDenied), errors.Is(err, event.ErrUserIsEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}
