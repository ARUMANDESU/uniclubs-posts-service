package validate

import (
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"log"
)

func user(value interface{}) error {
	u, ok := value.(*uniclubs_posts_service_v1_eventv1.UserObject)
	if !ok {
		return validation.NewInternalError(errors.New("user invalid type"))
	}
	return validation.ValidateStruct(u,
		validation.Field(&u.Id, validation.Required, validation.Min(1)),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Barcode, validation.Required),
	)
}

func club(value interface{}) error {
	c, ok := value.(*uniclubs_posts_service_v1_eventv1.ClubObject)
	if !ok {
		return validation.NewInternalError(errors.New("club invalid type"))
	}
	return validation.ValidateStruct(c,
		validation.Field(&c.Id, validation.Required, validation.Min(1)),
		validation.Field(&c.Name, validation.Required),
	)
}

func attachedFiles(value interface{}) error {
	a, ok := value.([]*uniclubs_posts_service_v1_eventv1.FileObject)
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
	c, ok := value.([]*uniclubs_posts_service_v1_eventv1.CoverImage)
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
