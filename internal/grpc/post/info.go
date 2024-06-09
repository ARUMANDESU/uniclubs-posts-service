package postgrpc

import (
	"context"
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
)

type InfoService interface {
}

func (s serverApi) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.PostObject, error) {
	//TODO implement me
	panic("implement me")
}
func (s serverApi) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	//TODO implement me
	panic("implement me")
}
