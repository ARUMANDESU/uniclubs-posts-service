package validate

import (
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-protos/gen/go/posts"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log"
)

func user(value interface{}) error {
	u, ok := value.(*eventv1.UserObject)
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
	c, ok := value.(*posts.ClubObject)
	if !ok {
		return validation.NewInternalError(errors.New("club invalid type"))
	}
	return validation.ValidateStruct(c,
		validation.Field(&c.Id, validation.Required, validation.Min(1)),
		validation.Field(&c.Name, validation.Required),
	)
}

func attachedFiles(value interface{}) error {
	a, ok := value.([]*posts.FileObject)
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
	c, ok := value.([]*posts.CoverImage)
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

func eventFilter(value interface{}) error {
	req, ok := value.(*eventv1.EventFilter)
	if !ok {
		return validation.NewInternalError(errors.New("list events filter invalid type"))
	}
	if req == nil {
		return nil
	}

	validStatuses := []any{
		domain.EventStatusDraft.String(),
		domain.EventStatusPending.String(),
		domain.EventStatusApproved.String(),
		domain.EventStatusRejected.String(),
		domain.EventStatusInProgress.String(),
		domain.EventStatusFinished.String(),
		domain.EventStatusCanceled.String(),
		domain.EventStatusArchived.String(),
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.UserId, validation.Min(0)),
		validation.Field(&req.ClubId, validation.Min(0)),
		validation.Field(&req.Tags, validation.Each(validation.Length(MinTagsLength, MaxTagsLength))),
		validation.Field(&req.FromDate, validation.Date(domain.TimeLayout)),
		validation.Field(&req.TillDate, validation.Date(domain.TimeLayout)),
		validation.Field(&req.Status,
			validation.Each(
				validation.In(validStatuses...).
					Error(fmt.Sprintf("event status must be valid: %v", validStatuses)),
			),
		),
	)
}

func updateMask(value interface{}) error {
	req, ok := value.(*fieldmaskpb.FieldMask)
	if !ok {
		return validation.NewInternalError(errors.New("update event invalid type"))
	}

	valideUpdateFields := []any{"title", "description", "type", "tags",
		"max_participants", "location_link", "location_university",
		"start_date", "end_date", "cover_images", "attached_images",
		"attached_files", "is_hidden_for_non_members",
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Paths,
			validation.Required,
			validation.Each(validation.In(valideUpdateFields...).
				Error(fmt.Sprintf("update fields must be one of %v", valideUpdateFields)),
			),
		),
	)
}

func createPostMask(value interface{}) error {
	req, ok := value.(*fieldmaskpb.FieldMask)
	if !ok {
		return validation.NewInternalError(errors.New("create post mask invalid type"))
	}

	validCreateFields := []any{"title", "description", "tags", "cover_images", "attached_files"}

	return validation.ValidateStruct(req,
		validation.Field(&req.Paths,
			validation.Required,
			validation.Each(validation.In(validCreateFields...).
				Error(fmt.Sprintf("create fields must be one of %v", validCreateFields)),
			),
		),
	)
}

func updatePostMask(value interface{}) error {
	req, ok := value.(*fieldmaskpb.FieldMask)
	if !ok {
		return validation.NewInternalError(errors.New("update post mask invalid type"))
	}

	validUpdateFields := []any{"title", "description", "tags", "cover_images", "attached_files"}

	return validation.ValidateStruct(req,
		validation.Field(&req.Paths,
			validation.Required,
			validation.Each(validation.In(validUpdateFields...).
				Error(fmt.Sprintf("update fields must be one of %v", validUpdateFields)),
			),
		),
	)
}
