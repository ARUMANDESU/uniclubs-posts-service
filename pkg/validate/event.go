package validate

import (
	"errors"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

const (
	MinTitleLength = 0
	MaxTitleLength = 500

	MinDescriptionLength = 0
	MaxDescriptionLength = 35000

	MinTagsLength = 2
	MaxTagsLength = 75

	MinLocationLink = 0
	MaxLocationLink = 2500

	MinLocationUniversity = 0
	MaxLocationUniversity = 250

	MinAttachedFileNameLength = 0
	MaxAttachedFileNameLength = 250

	MaxPosition           = 20
	MaxParticipantsNumber = 100000

	MaxReasonLength = 2000
)

func CreateEvent(value interface{}) error {
	req, ok := value.(*eventv1.CreateEventRequest)
	if !ok {
		return validation.NewInternalError(errors.New("create event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.Club, validation.Required, validation.By(club)),
		validation.Field(&req.User, validation.Required, validation.By(user)),
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

	startDateValidation := validation.Date(domain.TimeLayout).
		Max(base.AddDate(10, 0, 0)).
		Min(base.AddDate(-6, 0, 0))
	endDateValidation := validation.Date(domain.TimeLayout).
		Max(base.AddDate(10, 0, 0)).
		Min(base.AddDate(-6, 0, 0))

	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(0)),
		validation.Field(&req.Title, validation.Length(MinTitleLength, MaxTitleLength)),
		validation.Field(&req.Description, validation.Length(MinDescriptionLength, MaxDescriptionLength)),
		validation.Field(&req.Type, validation.In(domain.EventTypeUniversity.String(), domain.EventTypeIntraClub.String()).Error("event type must be UNIVERSITY or INTRA_CLUB")),
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

func ListEvents(value interface{}) error {
	req, ok := value.(*eventv1.ListEventsRequest)
	if !ok {
		return validation.NewInternalError(errors.New("list events invalid type"))
	}

	validSortBy := []any{domain.SortByDate.String(), domain.SortByParticipants.String(), domain.SortByType.String()}

	return validation.ValidateStruct(req,
		validation.Field(&req.Query, validation.Length(0, 1000)),
		validation.Field(&req.SortBy, validation.In(validSortBy...)),
		validation.Field(&req.SortOrder,
			validation.In(domain.SortOrderAsc.String(), domain.SortOrderDesc.String()),
		),
		validation.Field(&req.PageNumber, validation.Required),
		validation.Field(&req.PageSize, validation.Required),
		validation.Field(&req.Filter, validation.By(eventFilter)),
	)
}

func EventActionRequest(value interface{}) error {
	req, ok := value.(*eventv1.EventActionRequest)
	if !ok {
		return validation.NewInternalError(errors.New("event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(1)),
	)
}

func PublishEvent(value interface{}) error {
	req, ok := value.(*domain.Event)
	if !ok {
		return validation.NewInternalError(errors.New("publish event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.Status,
			validation.Required,
			validation.When(
				req.Type == domain.EventTypeUniversity,
				validation.In(domain.EventStatusApproved).Error("event must be APPROVED"),
			),
			validation.When(
				req.Type == domain.EventTypeIntraClub,
				validation.In(domain.EventStatusApproved, domain.EventStatusDraft).Error("event status must be APPROVED or DRAFT"),
			),
		),
		validation.Field(&req.Title,
			validation.Required.Error("event title is required"),
			validation.Length(3, MaxTitleLength).Error("event title must be between 3 and 500 characters"),
		),
		validation.Field(&req.Type,
			validation.Required.Error("event type is required"),
			validation.In(domain.EventTypeUniversity, domain.EventTypeIntraClub),
		),
		validation.Field(&req.StartDate, validation.Required.Error("start date is required")),
		validation.Field(&req.EndDate, validation.Required.Error("end date is required")),
		validation.Field(&req.CoverImages,
			validation.Required.Error("cover image is required"),
			validation.Length(1, 100).Error("cover image is required, add at least one image"),
		),
	)
}

func SendToReview(value interface{}) error {
	req, ok := value.(*domain.Event)
	if !ok {
		return validation.NewInternalError(errors.New("event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.Title,
			validation.Required.Error("event title is required"),
			validation.Length(3, MaxTitleLength).Error("event title must be between 3 and 500 characters"),
		),
		validation.Field(&req.Type,
			validation.Required.Error("event type is required"),
			validation.In(domain.EventTypeUniversity, domain.EventTypeIntraClub),
		),
		validation.Field(&req.StartDate, validation.Required.Error("start date is required")),
		validation.Field(&req.EndDate, validation.Required.Error("end date is required")),
		validation.Field(&req.CoverImages,
			validation.Required.Error("cover image is required"),
			validation.Length(1, 100).Error("cover image is required, add at least one image"),
		),
	)
}

func ApproveEventRequest(value interface{}) error {
	req, ok := value.(*eventv1.ApproveEventRequest)
	if !ok {
		return validation.NewInternalError(errors.New("approve event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.User, validation.Required, validation.By(user)),
	)
}

func RejectEventRequest(value interface{}) error {
	req, ok := value.(*eventv1.RejectEventRequest)
	if !ok {
		return validation.NewInternalError(errors.New("reject event invalid type"))
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.EventId, validation.Required),
		validation.Field(&req.User, validation.Required, validation.By(user)),
		validation.Field(&req.Reason, validation.Required, validation.Length(0, MaxReasonLength)),
	)
}
