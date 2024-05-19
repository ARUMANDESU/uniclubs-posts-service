package validate

import (
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

const (
	MinTitleLength            = 0
	MaxTitleLength            = 500
	MinDescriptionLength      = 0
	MaxDescriptionLength      = 35000
	MinTagsLength             = 2
	MaxTagsLength             = 75
	MinLocationLink           = 0
	MaxLocationLink           = 2500
	MinLocationUniversity     = 0
	MaxLocationUniversity     = 250
	MinAttachedFileNameLength = 0
	MaxAttachedFileNameLength = 250
	MaxPosition               = 20
	MaxParticipantsNumber     = 100000
)

func CreateEvent(value interface{}) error {
	req, ok := value.(*eventv1.CreateEventRequest)
	if !ok {
		return validation.NewInternalError(errors.New("create event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.Club, validation.By(club)),
		validation.Field(&req.User, validation.By(user)),
	)
}

func GetEvent(value interface{}) error {
	req, ok := value.(*eventv1.GetEventRequest)
	if !ok {
		return validation.NewInternalError(errors.New("get event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Min(0)),
	)
}

func UpdateEvent(value interface{}) error {
	req, ok := value.(*eventv1.UpdateEventRequest)
	if !ok {
		return validation.NewInternalError(errors.New("update event invalid type"))
	}
	base := time.Now()

	startDateValidation := validation.Date(time.RFC3339).
		Max(base.AddDate(10, 0, 0)).
		Min(base.AddDate(-6, 0, 0))
	endDateValidation := validation.Date(time.RFC3339).
		Max(base.AddDate(10, 0, 0)).
		Min(base.AddDate(-6, 0, 0))

	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
		validation.Field(&req.Title, validation.Length(MinTitleLength, MaxTitleLength)),
		validation.Field(&req.Description, validation.Length(MinDescriptionLength, MaxDescriptionLength)),
		validation.Field(&req.Type, validation.In("university", "intra-club")),
		validation.Field(&req.Tags, validation.Each(validation.Length(MinTagsLength, MaxTagsLength))),
		validation.Field(&req.MaxParticipants, validation.Min(0), validation.Max(MaxParticipantsNumber)),
		validation.Field(&req.StartDate, startDateValidation),
		validation.Field(&req.EndDate, endDateValidation),
		validation.Field(&req.LocationLink, validation.Length(MinLocationLink, MaxLocationLink)),
		validation.Field(&req.LocationUniversity, validation.Length(MinLocationUniversity, MaxLocationUniversity)),
		validation.Field(&req.CoverImages, validation.By(coverImages)),
		validation.Field(&req.AttachedImages, validation.By(attachedFiles)),
		validation.Field(&req.AttachedFiles, validation.By(attachedFiles)),
	)
}

func AddOrganizer(value interface{}) error {
	req, ok := value.(*eventv1.AddOrganizerRequest)
	if !ok {
		return validation.NewInternalError(errors.New("add organizer invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0), validation.NotIn(req.Target.Id).Error("target_id must be different from user_id")),
		validation.Field(&req.Target, validation.By(user)),
		validation.Field(&req.TargetClubId, validation.Required, validation.Min(0)),
	)
}

func RemoveOrganizer(value interface{}) error {
	req, ok := value.(*eventv1.RemoveOrganizerRequest)
	if !ok {
		return validation.NewInternalError(errors.New("remove organizer invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0), validation.NotIn(req.OrganizerId).Error("organizer_id must be different from user_id")),
		validation.Field(&req.OrganizerId, validation.Required, validation.Min(0)),
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
