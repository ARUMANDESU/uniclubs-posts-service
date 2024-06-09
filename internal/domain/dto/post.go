package dtos

import (
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
)

type CreatePostRequest struct {
	ClubId        int64               `json:"club_id"`
	UserId        int64               `json:"user_id"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	Tags          []string            `json:"tags"`
	CoverImages   []domain.CoverImage `json:"cover_images"`
	AttachedFiles []domain.File       `json:"attached_files"`
	Paths         map[string]bool     `json:"paths"`
}

type UpdatePostRequest struct {
	PostId        string              `json:"post_id"`
	UserId        int64               `json:"user_id"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	Tags          []string            `json:"tags"`
	CoverImages   []domain.CoverImage `json:"cover_images"`
	AttachedFiles []domain.File       `json:"attached_files"`
	Paths         map[string]bool     `json:"paths"`
}

type ActionRequest struct {
	PostId string `json:"post_id"`
	UserId int64  `json:"user_id"`
}

type ListPostsRequest struct {
	domain.BaseFilter
	ClubId int64           `json:"club_id"`
	Tags   []string        `json:"tags"`
	Paths  map[string]bool `json:"paths"`
}

func ToCreatePostRequest(post *postv1.CreatePostRequest) *CreatePostRequest {
	paths := make(map[string]bool)
	for _, path := range post.GetCreateMask().GetPaths() {
		paths[path] = true
	}

	return &CreatePostRequest{
		ClubId:        post.GetClubId(),
		UserId:        post.GetUserId(),
		Title:         post.GetTitle(),
		Description:   post.GetDescription(),
		Tags:          post.GetTags(),
		CoverImages:   domain.PbToCoverImages(post.GetCoverImages()),
		AttachedFiles: domain.PbToFiles(post.GetAttachedFiles()),
		Paths:         paths,
	}
}

func ToUpdatePostRequest(post *postv1.UpdatePostRequest) *UpdatePostRequest {
	paths := make(map[string]bool)
	for _, path := range post.GetUpdateMask().GetPaths() {
		paths[path] = true
	}

	return &UpdatePostRequest{
		PostId:        post.GetId(),
		UserId:        post.GetUserId(),
		Title:         post.GetTitle(),
		Description:   post.GetDescription(),
		Tags:          post.GetTags(),
		CoverImages:   domain.PbToCoverImages(post.GetCoverImages()),
		AttachedFiles: domain.PbToFiles(post.GetAttachedFiles()),
		Paths:         paths,
	}
}

func ToActionRequest(action *postv1.ActionRequest) *ActionRequest {
	return &ActionRequest{
		PostId: action.GetId(),
		UserId: action.GetUserId(),
	}
}

func ToListPostsRequest(l *postv1.ListPostsRequest) *ListPostsRequest {
	dto := &ListPostsRequest{
		BaseFilter: domain.BaseFilter{
			Page:     l.GetPage(),
			PageSize: l.GetPageSize(),
			Query:    l.GetQuery(),
			SortBy:   domain.SortBy(l.GetSortBy()),
		},
		ClubId: l.GetClubId(),
		Tags:   l.GetTags(),
	}

	paths := make(map[string]bool)
	for _, path := range l.GetListMask().GetPaths() {
		paths[path] = true
	}
	dto.Paths = paths

	if l.GetSortOrder() == "" {
		dto.BaseFilter.SortOrder = domain.SortOrderDesc
	} else {
		dto.BaseFilter.SortOrder = domain.SortOrder(l.GetSortOrder())
	}

	if l.GetPage() == 0 {
		dto.BaseFilter.Page = 1
	}
	if l.GetPageSize() == 0 {
		dto.BaseFilter.PageSize = 25
	}

	return dto
}
