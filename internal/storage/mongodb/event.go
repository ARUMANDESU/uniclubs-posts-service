package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Storage) CreateEvent(ctx context.Context, dto *dto.CreateEventDTO) (*domain.Event, error) {
	const op = "storage.mongodb.event.createEvent"

	insertResult, err := s.eventsCollection.InsertOne(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if insertResult.InsertedID == nil {
		return nil, fmt.Errorf("%s: no inserted id", op)
	}

	insertedID := insertResult.InsertedID.(primitive.ObjectID)
	var event domain.Event
	err = s.eventsCollection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &event, nil
}

func (s Storage) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	//TODO implement me
	panic("implement me")
}
