package validate

import (
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func AddCollaborator(value any) error {
	req, ok := value.(*eventv1.AddCollaboratorRequest)
	if !ok {
		return validation.NewInternalError(errors.New("add collaborator invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required.Error("missing event_id")),
		validation.Field(&req.UserId, validation.Required.Error("missing user_id")),
		validation.Field(&req.Club, validation.Required.Error("missing club"), validation.By(club)),
	)
}

func AddOrganizer(value interface{}) error {
	req, ok := value.(*eventv1.AddOrganizerRequest)
	if !ok {
		return validation.NewInternalError(errors.New("add organizer invalid type"))
	}
	if req.Target == nil {
		return validation.NewInternalError(errors.New("missing target user"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required.Error("missing event_id")),
		validation.Field(&req.Target, validation.Required.Error("missing target user"), validation.By(user)),
		validation.Field(&req.TargetClubId, validation.Required.Error("missing target's clubId"), validation.Min(0)),
		validation.Field(&req.UserId,
			validation.Required.Error("missing userId"),
			validation.Min(0),
			validation.When(req.Target != nil, validation.NotIn(req.Target.Id).Error("target_id must be different from user_id")),
		),
	)
}

func HandleInviteUser(value interface{}) error {
	req, ok := value.(*eventv1.HandleInviteUserRequest)
	if !ok {
		return validation.NewInternalError(errors.New("handle invite user invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.InviteId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
	)

}

func HandleInviteClub(value interface{}) error {
	req, ok := value.(*eventv1.HandleInviteClubRequest)
	if !ok {
		return validation.NewInternalError(errors.New("handle invite club invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.InviteId, validation.Required),
		validation.Field(&req.ClubId, validation.Required, validation.Min(0)),
		validation.Field(&req.User, validation.Required, validation.By(user)),
	)
}

func RevokeInvite(value interface{}) error {
	req, ok := value.(*eventv1.RevokeInviteRequest)
	if !ok {
		return validation.NewInternalError(errors.New("revoke invite invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.InviteId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
	)
}
