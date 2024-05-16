package validate

import (
	"errors"
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"log"
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
		validation.Field(&req.ClubId, validation.Required, validation.Min(0)),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
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

	startTimeValidation := validation.Date(time.RFC3339).
		Max(base.AddDate(10, 0, 0)).
		Min(base.AddDate(-6, 0, 0))
	endTimeValidation := validation.Date(time.RFC3339).
		Max(base.AddDate(10, 0, 0)).
		Min(base.AddDate(-6, 0, 0))

	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.ClubId, validation.Required, validation.Min(0)),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
		validation.Field(&req.Title, validation.Length(MinTitleLength, MaxTitleLength)),
		validation.Field(&req.Description, validation.Length(MinDescriptionLength, MaxDescriptionLength)),
		validation.Field(&req.Type, validation.In("university", "intra-club")),
		validation.Field(&req.Tags, validation.Each(validation.Length(MinTagsLength, MaxTagsLength))),
		validation.Field(&req.MaxParticipants, validation.Min(0), validation.Max(MaxParticipantsNumber)),
		validation.Field(&req.StartTime, startTimeValidation),
		validation.Field(&req.EndTime, endTimeValidation),
		validation.Field(&req.LocationLink, validation.Length(MinLocationLink, MaxLocationLink)),
		validation.Field(&req.LocationUniversity, validation.Length(MinLocationUniversity, MaxLocationUniversity)),
		validation.Field(&req.CoverImages, validation.By(coverImages)),
		validation.Field(&req.AttachedImages, validation.By(attachedFiles)),
		validation.Field(&req.AttachedFiles, validation.By(attachedFiles)),
	)

}

func attachedFiles(value interface{}) error {
	a, ok := value.([]*eventv1.FileObject)
	if !ok {
		log.Printf("file type: %T", value)
		return validation.NewInternalError(errors.New("attached files invalid type"))
	}

	for i, file := range a {
		err := validation.ValidateStruct(file,
			validation.Field(&file.Url, validation.Required, is.URL),
			validation.Field(&file.Name, validation.Required, validation.Length(MinAttachedFileNameLength, MaxAttachedFileNameLength)),
			validation.Field(&file.Type, validation.Required),
		)
		if err != nil {
			return validation.NewInternalError(fmt.Errorf("attached file %d: %w", i, err))
		}
	}

	return nil
}

func coverImages(value interface{}) error {
	c, ok := value.([]*eventv1.CoverImage)
	if !ok {
		log.Printf("image type: %T", value)
		return validation.NewInternalError(errors.New("cover images invalid type"))
	}

	for i, image := range c {
		err := validation.ValidateStruct(image,
			validation.Field(&image.Url, validation.Required, is.URL),
			validation.Field(&image.Name, validation.Required, validation.Length(MinAttachedFileNameLength, MaxAttachedFileNameLength)),
			validation.Field(&image.Type, validation.Required),
			validation.Field(&image.Position, validation.Required, validation.Min(0), validation.Max(MaxPosition)),
		)
		if err != nil {
			return validation.NewInternalError(fmt.Errorf("cover image %d: %w", i, err))
		}

	}

	return nil
}
