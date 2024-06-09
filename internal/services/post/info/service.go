package postinfo

import (
	"context"
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	postservice "github.com/arumandesu/uniclubs-posts-service/internal/services/post"
	"log/slog"
)

type Service struct {
	log          *slog.Logger
	postProvider PostProvider
	clubProvider ClubProvider
}

type PostProvider interface {
	GetPostById(ctx context.Context, postId string) (*domain.Post, error)
	ListPosts(ctx context.Context, filters *dtos.ListPostsRequest) ([]domain.Post, *domain.PaginationMetadata, error)
}

type ClubProvider interface {
	HasPermission(ctx context.Context, userId, clubId int64, permission clubv1.Permission) (bool, error)
}

func New(log *slog.Logger, postProvider PostProvider, clubProvider ClubProvider) *Service {
	return &Service{
		log:          log,
		postProvider: postProvider,
		clubProvider: clubProvider,
	}
}

func (s Service) GetPost(ctx context.Context, postId string, userId int64) (*domain.Post, error) {
	const op = "services.post.info.getPost"
	log := s.log.With(slog.String("op", op))

	post, err := s.postProvider.GetPostById(ctx, postId)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to get post", err)
	}

	return post, nil
}

func (s Service) ListPosts(ctx context.Context, filter *dtos.ListPostsRequest) ([]domain.Post, *domain.PaginationMetadata, error) {
	const op = "services.post.info.listPosts"
	log := s.log.With(slog.String("op", op))

	posts, metadata, err := s.postProvider.ListPosts(ctx, filter)
	if err != nil {
		return nil, nil, postservice.HandleError(log, "failed to list posts", err)
	}

	return posts, metadata, nil
}
