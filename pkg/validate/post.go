package validate

import (
	"errors"
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	MinPostDescriptionLength = 0
	MaxPostDescriptionLength = 3500

	MinPostTitleLength = 0
	MaxPostTitleLength = 255

	MaxPostTagsCount = 15
)

func CreatePostRequest(value interface{}) error {
	req, ok := value.(*postv1.CreatePostRequest)
	if !ok {
		return validation.NewInternalError(errors.New("list bans invalid type"))
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.ClubId, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(1)),
		validation.Field(&req.Title, validation.Required, validation.Length(MinPostTitleLength, MaxPostTitleLength)),
		validation.Field(&req.Description, validation.Length(MinPostDescriptionLength, MaxPostDescriptionLength)),
		validation.Field(&req.Tags, validation.Length(0, MaxPostTagsCount)),
		validation.Field(&req.CoverImages, validation.By(coverImages)),
		validation.Field(&req.AttachedFiles, validation.By(attachedFiles)),
		validation.Field(&req.CreateMask, validation.Required, validation.By(createPostMask)),
	)
}

func UpdatePostRequest(value interface{}) error {
	req, ok := value.(*postv1.UpdatePostRequest)
	if !ok {
		return validation.NewInternalError(errors.New("list bans invalid type"))
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Id, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(1)),
		validation.Field(&req.Title, validation.Required, validation.Length(MinPostTitleLength, MaxPostTitleLength)),
		validation.Field(&req.Description, validation.Required, validation.Length(MinPostDescriptionLength, MaxPostDescriptionLength)),
		validation.Field(&req.Tags, validation.Required, validation.Length(0, MaxPostTagsCount)),
		validation.Field(&req.CoverImages, validation.By(coverImages)),
		validation.Field(&req.AttachedFiles, validation.By(attachedFiles)),
		validation.Field(&req.UpdateMask, validation.Required, validation.By(updatePostMask)),
	)
}

func PostActionRequest(value interface{}) error {
	req, ok := value.(*postv1.ActionRequest)
	if !ok {
		return validation.NewInternalError(errors.New("list bans invalid type"))
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Id, validation.Required),
		validation.Field(&req.UserId, validation.Required, validation.Min(1)),
	)
}

func GetPostRequest(value interface{}) error {
	req, ok := value.(*postv1.GetPostRequest)
	if !ok {
		return validation.NewInternalError(errors.New("list bans invalid type"))
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Id, validation.Required),
		validation.Field(&req.UserId, validation.Min(0)),
	)
}

func ListPostsRequest(value interface{}) error {
	req, ok := value.(*postv1.ListPostsRequest)
	if !ok {
		return validation.NewInternalError(errors.New("list bans invalid type"))
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Page, validation.Min(0)),
		validation.Field(&req.PageSize, validation.Min(0)),
		validation.Field(&req.Query, validation.Length(0, 255)),
		validation.Field(&req.SortBy, validation.In("created_at", "title")),
		validation.Field(&req.SortOrder, validation.In("asc", "desc")),
		validation.Field(&req.ClubId, validation.Min(0)),
		validation.Field(&req.Tags, validation.Length(0, MaxPostTagsCount)),
	)
}
