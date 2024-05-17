package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Invite struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	EventId primitive.ObjectID `json:"event_id,omitempty" bson:"event_id"`
	ClubId  int64              `json:"club_id,omitempty" bson:"club_id"`
}

type UserInvite struct {
	Invite `bson:",inline"`
	UserId int64 `json:"user_id,omitempty" bson:"user_id"`
}

func (s Storage) SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.UserInvite, error) {
	const op = "storage.mongodb.sendJoinRequestToUser"

	eventObjectId, err := primitive.ObjectIDFromHex(dto.EventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	invite := UserInvite{
		Invite: Invite{
			ID:      primitive.NewObjectID(),
			EventId: eventObjectId,
			ClubId:  dto.TargetClubId,
		},
		UserId: dto.TargetId,
	}

	_, err = s.inviteCollection.InsertOne(ctx, invite)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ToDomainUserInvite(invite), nil

}

func (s Storage) GetJoinRequests(ctx context.Context, eventId string) ([]domain.UserInvite, error) {
	const op = "storage.mongodb.getJoinRequests"

	find, err := s.inviteCollection.Find(ctx, bson.D{{"event_id", eventId}})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var invites []UserInvite
	err = find.All(ctx, &invites)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ToDomainUserInvites(invites), nil
}

func (s Storage) GetJoinRequestByUserId(ctx context.Context, userId int64) (*domain.UserInvite, error) {
	const op = "storage.mongodb.getJoinRequestByUserId"

	var invite UserInvite
	err := s.inviteCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&invite)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ToDomainUserInvite(invite), nil
}

func (s Storage) GetJoinRequestsById(ctx context.Context, requestId string) (*domain.UserInvite, error) {
	const op = "storage.mongodb.getJoinRequestsById"

	objectID, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var invite UserInvite
	err = s.inviteCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&invite)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ToDomainUserInvite(invite), nil
}
