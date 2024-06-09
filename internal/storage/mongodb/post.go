package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	dtos "github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage/mongodb/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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
	const op = "storage.mongodb.post.updatePost"

	daoPost := dao.PostFromDomain(post)

	lastUpdated := daoPost.UpdatedAt
	daoPost.UpdatedAt = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.M{"_id": daoPost.ID, "updated_at": lastUpdated}
	update := bson.M{"$set": daoPost}

	err := s.postsCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&daoPost)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrOptimisticLockingFailed)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.PostToDomain(daoPost), nil
}

func (s *Storage) DeletePost(ctx context.Context, postId string) (*domain.Post, error) {
	const op = "storage.mongodb.post.deletePost"

	objectID, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
	}

	var post dao.Post
	err = s.postsCollection.FindOneAndDelete(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.PostToDomain(&post), nil
}

func (s *Storage) HidePost(ctx context.Context, postId string) (*domain.Post, error) {
	// todo: implement me
	panic("implement me")
}

func (s *Storage) UnhidePost(ctx context.Context, postId string) (*domain.Post, error) {
	// todo: implement me
	panic("implement me")
}

func (s *Storage) GetPostById(ctx context.Context, postId string) (*domain.Post, error) {
	const op = "storage.mongodb.post.getPostById"

	objectID, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
	}

	var post dao.Post
	err = s.postsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.PostToDomain(&post), nil
}

func (s *Storage) ListPosts(ctx context.Context, filters *dtos.ListPostsRequest) ([]domain.Post, *domain.PaginationMetadata, error) {
	const op = "storage.mongodb.post.listPosts"

	filter := constructPostFilter(filters)

	totalRecords, err := s.postsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	if totalRecords == 0 {
		return nil, nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	opts := options.Find()
	//todo: make this work later
	//opts.SetSort(constructEventSortBy(filters.BaseFilter))
	opts.SetSkip(int64(filters.Offset()))
	opts.SetLimit(int64(filters.Limit()))

	cursor, err := s.postsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, handleError(op, err)
	}
	defer cursor.Close(ctx)

	var posts []dao.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	paginationMetadata := domain.CalculatePaginationMetadata(int32(totalRecords), filters.Page, filters.PageSize)

	return dao.PostsToDomain(posts), &paginationMetadata, nil
}

func constructPostFilter(filters *dtos.ListPostsRequest) bson.M {
	m := bson.M{}

	if filters.ClubId != 0 {
		m["club._id"] = filters.ClubId
	}

	if filters.Tags != nil && len(filters.Tags) > 0 {
		m["tags"] = bson.M{"$in": filters.Tags}
	}

	return m
}

func constructPostSortBy(filter *domain.BaseFilter) bson.M {
	sortBy := bson.M{}

	switch filter.SortBy {
	case domain.SortByDate:
		sortBy["created_at"] = constructEventSortOrder(filter.SortOrder)
	case domain.SortByParticipants:
		sortBy["participants"] = constructEventSortOrder(filter.SortOrder)
	case domain.SortByType:
		sortBy["type"] = constructEventSortOrder(filter.SortOrder)
	default:
		sortBy["start_date"] = 1
	}

	return sortBy
}
