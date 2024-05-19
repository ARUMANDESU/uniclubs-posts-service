package validate

import (
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	validation "github.com/go-ozzo/ozzo-validation"
)

func AddCollaborator(value any) error {
	req, ok := value.(*eventv1.AddCollaboratorRequest)
	if !ok {
		return validation.NewInternalError(errors.New("add collaborator invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required),
		validation.Field(&req.Club, validation.By(club)),
	)
}
