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

type InfoService interface {
	GetPost(ctx context.Context, postId string, userId int64) (*domain.Post, error)
	ListPosts(ctx context.Context, filter *dtos.ListPostsRequest) ([]domain.Post, *domain.PaginationMetadata, error)
}

func (s serverApi) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.PostObject, error) {
	err := validate.GetPostRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	post, err := s.info.GetPost(ctx, req.GetId(), req.GetUserId())
	if err != nil {
		return nil, handleServiceError(err)
	}

	return domain.PostToPb(post), nil
}
func (s serverApi) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	err := validate.ListPostsRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	posts, metadata, err := s.info.ListPosts(ctx, dtos.ToListPostsRequest(req))
	if err != nil {
		return nil, handleServiceError(err)
	}

	return &postv1.ListPostsResponse{
		Posts:    domain.PostsToPb(posts),
		Metadata: metadata.ToProto(),
	}, nil
}
