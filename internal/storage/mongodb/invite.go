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
)

func (s Storage) SendJoinRequestToUser(ctx context.Context, dto *dto.SendJoinRequestToUser) (*domain.UserInvite, error) {
	const op = "storage.mongodb.sendJoinRequestToUser"

	eventObjectId, err := primitive.ObjectIDFromHex(dto.EventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	invite := dao.OrganizerInvite{
		ID:      primitive.NewObjectID(),
		EventId: eventObjectId,
		ClubId:  dto.TargetClubId,
		ByWhoId: dto.UserId,
		User:    dao.UserFromDomainUser(dto.Target),
	}

	_, err = s.invitesCollection.InsertOne(ctx, invite)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainUserInvite(invite), nil

}

func (s Storage) GetJoinRequests(ctx context.Context, eventId string) ([]domain.UserInvite, error) {
	const op = "storage.mongodb.getJoinRequests"

	find, err := s.invitesCollection.Find(ctx, bson.D{{"event_id", eventId}})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var invites []dao.OrganizerInvite
	err = find.All(ctx, &invites)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainUserInvites(invites), nil
}

func (s Storage) GetJoinRequestByUserId(ctx context.Context, userId int64) (*domain.UserInvite, error) {
	const op = "storage.mongodb.getJoinRequestByUserId"

	var invite dao.OrganizerInvite
	err := s.invitesCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&invite)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainUserInvite(invite), nil
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

	var invite dao.OrganizerInvite
	err = s.invitesCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&invite)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainUserInvite(invite), nil
}

func (s Storage) DeleteJoinRequest(ctx context.Context, requestId string) error {
	const op = "storage.mongodb.deleteJoinRequest"

	objectID, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.invitesCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
