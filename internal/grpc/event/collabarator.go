package eventgrpc

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
	SendJoinRequestToUser(ctx context.Context, dto *dtos.SendJoinRequestToUser) (*domain.Event, error)
	AcceptUserJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error)
	RejectUserJoinRequest(ctx context.Context, inviteId string, userId int64) (domain.Event, error)
	KickOrganizer(ctx context.Context, eventId string, userId, targetId int64) (*domain.Event, error)
	RevokeInviteOrganizer(ctx context.Context, inviteId string, userId int64) error
}

type CollaboratorService interface {
	SendJoinRequestToClub(ctx context.Context, dto *dtos.SendJoinRequestToClub) (*domain.Event, error)
	AcceptClubJoinRequest(ctx context.Context, dto *dtos.AcceptJoinRequestClub) (domain.Event, error)
	RejectClubJoinRequest(ctx context.Context, inviteId string, clubId int64) (domain.Event, error)
	KickClub(ctx context.Context, eventId string, userId, clubId int64) (*domain.Event, error)
	RevokeInviteClub(ctx context.Context, inviteId string, userId int64) error
}

func (s serverApi) AddCollaborator(ctx context.Context, req *eventv1.AddCollaboratorRequest) (*empty.Empty, error) {
	err := validate.AddCollaborator(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.collaborator.SendJoinRequestToClub(ctx, dtos.AddCollaboratorRequestToClubToDTO(req))
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, eventservice.ErrInviteAlreadyExists), errors.Is(err, eventservice.ErrClubAlreadyCollaborator):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}

func (s serverApi) RemoveCollaborator(ctx context.Context, req *eventv1.RemoveCollaboratorRequest) (*eventv1.EventObject, error) {
	err := validate.RemoveCollaborator(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.collaborator.KickClub(ctx, req.GetEventId(), req.GetUserId(), req.GetClubId())
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrInviteNotFound), errors.Is(err, eventservice.ErrCollaboratorNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID), errors.Is(err, eventservice.ErrClubMismatch):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied), errors.Is(err, eventservice.ErrClubIsEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, eventservice.ErrEventUpdateConflict):
			return nil, status.Error(codes.Aborted, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) HandleInviteClub(ctx context.Context, req *eventv1.HandleInviteClubRequest) (*eventv1.EventObject, error) {
	err := validate.HandleInviteClub(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var event domain.Event
	if req.GetAction() == eventv1.HandleInvite_Action_ACCEPT {
		event, err = s.collaborator.AcceptClubJoinRequest(ctx, dtos.AcceptJoinRequestClubToDTO(req))
	} else if req.GetAction() == eventv1.HandleInvite_Action_REJECT {
		event, err = s.collaborator.RejectClubJoinRequest(ctx, req.InviteId, req.ClubId)
	}

	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrInviteNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID), errors.Is(err, eventservice.ErrClubMismatch):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, eventservice.ErrEventUpdateConflict):
			return nil, status.Error(codes.Aborted, err.Error())
		case errors.Is(err, eventservice.ErrUserAlreadyOrganizer):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) RevokeInviteClub(ctx context.Context, req *eventv1.RevokeInviteRequest) (*emptypb.Empty, error) {
	err := validate.RevokeInvite(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.collaborator.RevokeInviteClub(ctx, req.GetInviteId(), req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound), errors.Is(err, eventservice.ErrInviteNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied), errors.Is(err, eventservice.ErrUserIsEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}

func (s serverApi) AddOrganizer(ctx context.Context, req *eventv1.AddOrganizerRequest) (*empty.Empty, error) {
	err := validate.AddOrganizer(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.organizer.SendJoinRequestToUser(ctx, dtos.AddOrganizerRequestToUserToDTO(req))
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, eventservice.ErrInviteAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, eventservice.ErrUserIsFromAnotherClub):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}

func (s serverApi) RemoveOrganizer(ctx context.Context, req *eventv1.RemoveOrganizerRequest) (*eventv1.EventObject, error) {
	err := validate.RemoveOrganizer(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.organizer.KickOrganizer(ctx, req.GetEventId(), req.GetUserId(), req.GetOrganizerId())
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound),
			errors.Is(err, eventservice.ErrUserIsNotEventOrganizer),
			errors.Is(err, eventservice.ErrOrganizerNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied), errors.Is(err, eventservice.ErrUserIsEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, eventservice.ErrEventUpdateConflict):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) HandleInviteUser(ctx context.Context, req *eventv1.HandleInviteUserRequest) (*eventv1.EventObject, error) {
	err := validate.HandleInviteUser(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var event domain.Event
	if req.GetAction() == eventv1.HandleInvite_Action_ACCEPT {
		event, err = s.organizer.AcceptUserJoinRequest(ctx, req.InviteId, req.UserId)
	} else if req.GetAction() == eventv1.HandleInvite_Action_REJECT {
		event, err = s.organizer.RejectUserJoinRequest(ctx, req.InviteId, req.UserId)
	}

	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrInviteNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, eventservice.ErrUserAlreadyOrganizer):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, eventservice.ErrEventUpdateConflict):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return event.ToProto(), nil
}

func (s serverApi) RevokeInviteUser(ctx context.Context, req *eventv1.RevokeInviteRequest) (*emptypb.Empty, error) {
	err := validate.RevokeInvite(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.organizer.RevokeInviteOrganizer(ctx, req.GetInviteId(), req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, eventservice.ErrEventNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, eventservice.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, eventservice.ErrPermissionsDenied), errors.Is(err, eventservice.ErrUserIsEventOwner):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &empty.Empty{}, nil
}
