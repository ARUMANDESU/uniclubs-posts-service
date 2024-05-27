package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage/mongodb/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

func (s *Storage) BanParticipant(ctx context.Context, dto *dtos.BanParticipant) error {
	const op = "storage.mongodb.event.banParticipant"

	objectID, err := primitive.ObjectIDFromHex(dto.EventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	banRecord := dao.BanRecord{
		EventId:  objectID,
		UserId:   dto.ParticipantId,
		BannedAt: time.Now(),
		Reason:   dto.Reason,
		ByWhoId:  dto.UserId,
	}

	_, err = s.bansCollection.InsertOne(ctx, banRecord)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UnBanParticipant(ctx context.Context, eventId string, userId int64) error {
	const op = "storage.mongodb.event.unBanParticipant"

	objectID, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.bansCollection.DeleteOne(ctx, bson.M{"event_id": objectID, "user_id": userId})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetBanRecord(ctx context.Context, eventId string, userId int64) (*domain.BanRecord, error) {
	const op = "storage.mongodb.event.getBanRecord"

	objectID, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var banRecord dao.BanRecord
	err = s.bansCollection.FindOne(ctx, bson.M{"event_id": objectID, "user_id": userId}).Decode(&banRecord)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, storage.ErrBanRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.BanRecordToDomain(banRecord), nil
}
