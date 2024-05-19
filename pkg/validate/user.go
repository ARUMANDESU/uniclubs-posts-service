package validate

import (
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	validation "github.com/go-ozzo/ozzo-validation"
)

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
