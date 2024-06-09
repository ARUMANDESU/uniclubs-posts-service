package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage/mongodb/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	const op = "storage.mongodb.post.createPost"

	daoPost := dao.PostFromDomain(post)

	insertResult, err := s.postsCollection.InsertOne(ctx, daoPost)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if insertResult.InsertedID == nil {
		return nil, fmt.Errorf("%s: no inserted id", op)
	}

	insertedID := insertResult.InsertedID.(primitive.ObjectID)

	err = s.postsCollection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&daoPost)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.PostToDomain(daoPost), nil
}

func (s *Storage) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	// todo: implement me
	panic("implement me")
}

func (s *Storage) DeletePost(ctx context.Context, postId string) (*domain.Post, error) {
	// todo: implement me
	panic("implement me")
}

func (s *Storage) HidePost(ctx context.Context, postId string) (*domain.Post, error) {
	// todo: implement me
	panic("implement me")
}

func (s *Storage) UnhidePost(ctx context.Context, postId string) (*domain.Post, error) {
	// todo: implement me
	panic("implement me")
}
