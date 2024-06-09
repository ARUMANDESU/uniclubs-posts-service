package postgrpc

import (
	"context"
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ManagementService interface {
	CreatePost(ctx context.Context, dto *dtos.CreatePostRequest) (*domain.Post, error)
	UpdatePost(ctx context.Context, dto *dtos.UpdatePostRequest) (*domain.Post, error)
	DeletePost(ctx context.Context, dto *dtos.ActionRequest) (*domain.Post, error)
	HidePost(ctx context.Context, dto *dtos.ActionRequest) (*domain.Post, error)
	UnhidePost(ctx context.Context, dto *dtos.ActionRequest) (*domain.Post, error)
}

func (s serverApi) CreatePost(ctx context.Context, req *postv1.CreatePostRequest) (*postv1.PostObject, error) {
	err := validate.CreatePostRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	post, err := s.management.CreatePost(ctx, dtos.ToCreatePostRequest(req))
	if err != nil {
		return nil, handleServiceError(err)
	}

	return domain.PostToPb(post), nil
}

func (s serverApi) UpdatePost(ctx context.Context, req *postv1.UpdatePostRequest) (*postv1.PostObject, error) {
	err := validate.UpdatePostRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	post, err := s.management.UpdatePost(ctx, dtos.ToUpdatePostRequest(req))
	if err != nil {
		return nil, handleServiceError(err)
	}

	return domain.PostToPb(post), nil
}

func (s serverApi) DeletePost(ctx context.Context, req *postv1.ActionRequest) (*postv1.PostObject, error) {
	err := validate.PostActionRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	post, err := s.management.DeletePost(ctx, dtos.ToActionRequest(req))
	if err != nil {
		return nil, handleServiceError(err)
	}

	return domain.PostToPb(post), nil
}

func (s serverApi) HidePost(ctx context.Context, req *postv1.ActionRequest) (*postv1.PostObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s serverApi) UnhidePost(ctx context.Context, req *postv1.ActionRequest) (*postv1.PostObject, error) {
	//TODO implement me
	panic("implement me")
}
