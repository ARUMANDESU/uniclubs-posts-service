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
	"strings"
	"time"
)

func (s *Storage) CreateEvent(ctx context.Context, club *domain.Club, user *domain.User) (*domain.Event, error) {
	const op = "storage.mongodb.event.createEvent"

	event := dao.Event{
		ID:        primitive.NewObjectID(),
		ClubId:    club.ID,
		OwnerId:   user.ID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Status:    domain.EventStatusDraft,
	}

	event.AddOrganizer(dao.OrganizerFromDomainUser(*user, club.ID))
	event.AddCollaboratorClub(dao.ClubFromDomainClub(*club))

	insertResult, err := s.eventsCollection.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if insertResult.InsertedID == nil {
		return nil, fmt.Errorf("%s: no inserted id", op)
	}

	insertedID := insertResult.InsertedID.(primitive.ObjectID)

	err = s.eventsCollection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainEvent(event), nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*domain.Event, error) {
	const op = "storage.mongodb.event.getEvent"

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var event dao.Event
	err = s.eventsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&event)
	if err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrEventNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainEvent(event), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	const op = "storage.mongodb.event.updateEvent"

	eventModel := dao.EventToModel(event)

	lastUpdated := event.UpdatedAt
	event.UpdatedAt = time.Now()

	_, err := s.eventsCollection.UpdateOne(ctx, bson.M{"_id": eventModel.ID, "updated_at": lastUpdated}, bson.M{"$set": eventModel})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrOptimisticLockingFailed)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.eventsCollection.FindOne(ctx, bson.M{"_id": eventModel.ID}).Decode(&eventModel)
	if err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrEventNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainEvent(eventModel), nil
}

func (s *Storage) DeleteEventById(ctx context.Context, eventId string) error {
	const op = "storage.mongodb.event.deleteEventById"

	objectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.eventsCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("%s: %w", storage.ErrEventNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
