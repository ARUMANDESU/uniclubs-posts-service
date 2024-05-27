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

func (s *Storage) GetEventParticipant(ctx context.Context, eventId string, userId int64) (*domain.Participant, error) {
	const op = "storage.mongodb.event.getEventParticipant"

	objectID, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var participant dao.Participant
	err = s.participantsCollection.FindOne(ctx, bson.M{"event_id": objectID, "user._id": userId}).Decode(&participant)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, storage.ErrParticipantNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ParticipantToDomain(participant), nil
}

func (s *Storage) AddEventParticipant(ctx context.Context, participant *domain.Participant) error {
	const op = "storage.mongodb.event.addEventParticipant"
	participantDAO, err := dao.ParticipantFromDomain(participant)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.participantsCollection.InsertOne(ctx, participantDAO)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteEventParticipant(ctx context.Context, eventId string, userId int64) error {
	const op = "storage.mongodb.event.deleteEventParticipant"

	objectID, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.participantsCollection.DeleteOne(ctx, bson.M{"event_id": objectID, "user._id": userId})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
