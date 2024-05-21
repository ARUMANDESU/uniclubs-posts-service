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
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func (s *Storage) CreateEvent(ctx context.Context, club domain.Club, user domain.User) (*domain.Event, error) {
	const op = "storage.mongodb.event.createEvent"

	event := dao.Event{
		ID:        primitive.NewObjectID(),
		ClubId:    club.ID,
		OwnerId:   user.ID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Status:    domain.EventStatusDraft.String(),
	}

	event.AddOrganizer(dao.OrganizerFromDomainUser(user, club.ID))
	event.AddCollaboratorClub(dao.ClubFromDomainClub(club))

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
			return fmt.Errorf("%s: %w", op, storage.ErrEventNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	//todo: delete invites related to the event

	return nil
}

func (s *Storage) ListEvents(ctx context.Context, filters domain.Filters) ([]domain.Event, *domain.PaginationMetadata, error) {
	const op = "storage.mongodb.event.listEvents"

	var filter bson.M
	if filters.Query != "" {
		filter = bson.M{"$text": bson.M{"$search": filters.Query}}
	} else {
		filter = bson.M{}
	}

	if filters.ClubId != 0 {
		filter["club_id"] = filters.ClubId
	}

	if filters.UserId != 0 {
		filter["$or"] = []bson.M{
			{"owner_id": filters.UserId},
			{"organizers.id": filters.UserId},
		}
	}

	if filters.Tags != nil && len(filters.Tags) > 0 {
		filter["tags"] = bson.M{"$in": filters.Tags}
	}

	if !filters.FromDate.IsZero() {
		filter["created_at"] = bson.M{"$gte": filters.FromDate}
	}

	if !filters.ToDate.IsZero() {
		filter["created_at"] = bson.M{"$lte": filters.ToDate}
	}

	if filters.Status != "" {
		filter["status"] = filters.Status
	}

	totalRecords, err := s.eventsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	opts := options.Find()
	opts.SetSkip(int64(filters.Offset()))
	opts.SetLimit(int64(filters.Limit()))

	cursor, err := s.eventsCollection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, fmt.Errorf("%s: %w", op, storage.ErrEventNotFound)
		}
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	var events []dao.Event
	if err = cursor.All(ctx, &events); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	paginationMetadata := domain.CalculatePaginationMetadata(int32(totalRecords), filters.Page, filters.PageSize)

	return dao.ToDomainEvents(events), &paginationMetadata, nil

}
