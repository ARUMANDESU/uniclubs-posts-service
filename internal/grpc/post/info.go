package postgrpc

import (
	"context"
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InfoService interface {
	GetPost(ctx context.Context, postId string, userId int64) (*domain.Post, error)
}

func (s serverApi) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.PostObject, error) {
	err := validate.GetPostRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	post, err := s.Services.info.GetPost(ctx, req.GetId(), req.GetUserId())
	if err != nil {
		return nil, handleServiceError(err)
	}

	return domain.PostToPb(post), nil
}
func (s serverApi) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	//TODO implement me
	panic("implement me")
}
