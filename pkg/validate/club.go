package validate

import (
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	validation "github.com/go-ozzo/ozzo-validation"
)

func RemoveCollaborator(value any) error {
	req, ok := value.(*eventv1.RemoveCollaboratorRequest)
	if !ok {
		return validation.NewInternalError(errors.New("remove organizer invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
		validation.Field(&req.ClubId, validation.Required, validation.Min(0)),
	)
}
