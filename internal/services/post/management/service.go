package postmanagement

import (
	"context"
	"fmt"
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	postservice "github.com/arumandesu/uniclubs-posts-service/internal/services/post"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"time"
)

type Service struct {
	log          *slog.Logger
	postStorage  PostStorage
	clubProvider ClubProvider
}

type PostStorage interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	DeletePost(ctx context.Context, postId string) (*domain.Post, error)
	HidePost(ctx context.Context, postId string) (*domain.Post, error)
	UnhidePost(ctx context.Context, postId string) (*domain.Post, error)
	GetPostById(ctx context.Context, postId string) (*domain.Post, error)
}

type ClubProvider interface {
	GetClubById(ctx context.Context, clubId int64) (*domain.Club, error)
	HasPermission(ctx context.Context, userId, clubId int64, permission clubv1.Permission) (bool, error)
}

func New(log *slog.Logger, postStorage PostStorage, clubProvider ClubProvider) *Service {
	return &Service{
		log:          log,
		postStorage:  postStorage,
		clubProvider: clubProvider,
	}
}

func (s Service) CreatePost(ctx context.Context, dto *dtos.CreatePostRequest) (*domain.Post, error) {
	const op = "services.post.management.createPost"
	log := s.log.With(slog.String("op", op))

	hasPermission, err := s.clubProvider.HasPermission(ctx, dto.UserId, dto.ClubId, clubv1.Permission_PERMISSION_MANAGE_POSTS)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to check permission", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("%w: user %d does not have permission to manage posts in club %d", postservice.ErrPermissionDenied, dto.UserId, dto.ClubId)
	}

	club, err := s.clubProvider.GetClubById(ctx, dto.ClubId)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to get club by id", err)
	}

	now := time.Now()

	post := &domain.Post{
		ID:        primitive.NewObjectID().Hex(),
		Club:      *club,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if dto.Paths["title"] {
		post.Title = dto.Title
	}
	if dto.Paths["description"] {
		post.Description = dto.Description
	}
	if dto.Paths["tags"] {
		post.Tags = dto.Tags
	}
	if dto.Paths["cover_images"] {
		post.CoverImages = dto.CoverImages
	}
	if dto.Paths["attached_files"] {
		post.AttachedFiles = dto.AttachedFiles
	}

	post, err = s.postStorage.CreatePost(ctx, post)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to create post", err)
	}

	return post, nil
}

func (s Service) UpdatePost(ctx context.Context, dto *dtos.UpdatePostRequest) (*domain.Post, error) {
	const op = "services.post.management.updatePost"
	log := s.log.With(slog.String("op", op))

	post, err := s.postStorage.GetPostById(ctx, dto.PostId)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to get post by id", err)
	}

	hasPermission, err := s.clubProvider.HasPermission(ctx, dto.UserId, post.Club.ID, clubv1.Permission_PERMISSION_MANAGE_POSTS)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to check permission", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("%w: user %d does not have permission to manage posts in club %d", postservice.ErrPermissionDenied, dto.UserId, post.Club.ID)
	}

	if dto.Paths["title"] {
		post.Title = dto.Title
	}
	if dto.Paths["description"] {
		post.Description = dto.Description
	}
	if dto.Paths["tags"] {
		post.Tags = dto.Tags
	}
	if dto.Paths["cover_images"] {
		post.CoverImages = dto.CoverImages
	}
	if dto.Paths["attached_files"] {
		post.AttachedFiles = dto.AttachedFiles
	}

	post, err = s.postStorage.UpdatePost(ctx, post)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to update post", err)
	}

	return post, nil
}

func (s Service) DeletePost(ctx context.Context, dto *dtos.ActionRequest) (*domain.Post, error) {
	const op = "services.post.management.deletePost"
	log := s.log.With(slog.String("op", op))

	post, err := s.postStorage.GetPostById(ctx, dto.PostId)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to get post by id", err)
	}

	hasPermission, err := s.clubProvider.HasPermission(ctx, dto.UserId, post.Club.ID, clubv1.Permission_PERMISSION_MANAGE_POSTS)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to check permission", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("%w: user %d does not have permission to manage posts in club %d", postservice.ErrPermissionDenied, dto.UserId, post.Club.ID)
	}

	post, err = s.postStorage.DeletePost(ctx, dto.PostId)
	if err != nil {
		return nil, postservice.HandleError(log, "failed to delete post", err)
	}

	return post, nil
}

func (s Service) HidePost(ctx context.Context, dto *dtos.ActionRequest) (*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) UnhidePost(ctx context.Context, dto *dtos.ActionRequest) (*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}
