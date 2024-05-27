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

func (s *Storage) CreateJoinRequestToUser(ctx context.Context, dto *dtos.SendJoinRequestToUser) (*domain.UserInvite, error) {
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

func (s *Storage) GetUserJoinRequests(ctx context.Context, eventId string) ([]domain.UserInvite, error) {
	const op = "storage.mongodb.getJoinRequests"

	find, err := s.invitesCollection.Find(ctx, bson.D{{"event_id", eventId}})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer find.Close(ctx)

	var invites []dao.OrganizerInvite
	err = find.All(ctx, &invites)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainUserInvites(invites), nil
}

func (s *Storage) GetJoinRequestByUserId(ctx context.Context, eventId string, userId int64) (*domain.UserInvite, error) {
	const op = "storage.mongodb.getJoinRequestByUserId"

	eventObjectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{
		"event_id": eventObjectId,
		"user.id":  userId,
	}

	var invite dao.OrganizerInvite
	err = s.invitesCollection.FindOne(ctx, filter).Decode(&invite)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainUserInvite(invite), nil
}

func (s *Storage) GetJoinRequestsByUserInviteId(ctx context.Context, requestId string) (*domain.UserInvite, error) {
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

func (s *Storage) DeleteInvite(ctx context.Context, inviteId string) error {
	const op = "storage.mongodb.deleteJoinRequest"

	objectID, err := primitive.ObjectIDFromHex(inviteId)
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

/*

Club Invite Storage

*/

func (s *Storage) CreateJoinRequestToClub(ctx context.Context, dto *dtos.SendJoinRequestToClub) (*domain.Invite, error) {
	const op = "storage.mongodb.sendJoinRequestToClub"

	eventObjectId, err := primitive.ObjectIDFromHex(dto.EventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	invite := dao.ClubInvite{
		ID:      primitive.NewObjectID(),
		EventId: eventObjectId,
		Club:    dao.ClubFromDomainClub(dto.Club),
	}

	_, err = s.invitesCollection.InsertOne(ctx, invite)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainInvite(invite), nil
}

func (s *Storage) GetClubJoinRequests(ctx context.Context, eventId string) ([]domain.Invite, error) {
	const op = "storage.mongodb.getClubJoinRequests"

	find, err := s.invitesCollection.Find(ctx, bson.D{{"event_id", eventId}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer find.Close(ctx)

	var invites []dao.ClubInvite
	err = find.All(ctx, &invites)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainInvites(invites), nil
}

func (s *Storage) GetJoinRequestsByClubInviteId(ctx context.Context, inviteId string) (*domain.Invite, error) {
	const op = "storage.mongodb.getJoinRequestsById"

	objectID, err := primitive.ObjectIDFromHex(inviteId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var invite dao.ClubInvite
	err = s.invitesCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&invite)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainInvite(invite), nil
}

func (s *Storage) GetJoinRequestByClubId(ctx context.Context, eventId string, clubId int64) (*domain.Invite, error) {
	const op = "storage.mongodb.getJoinRequestByClubId"

	eventObjectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{
		"event_id": eventObjectId,
		"club._id": clubId,
	}

	var invite dao.ClubInvite
	err = s.invitesCollection.FindOne(ctx, filter).Decode(&invite)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainInvite(invite), nil
}

func (s *Storage) GetUserInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.UserInvite, error) {
	const op = "storage.mongodb.getUserInvites"

	filter := bson.M{}
	if dto.EventId != "" {
		eventObjectId, err := primitive.ObjectIDFromHex(dto.EventId)
		if err != nil {
			if errors.Is(err, primitive.ErrInvalidHex) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		filter["event_id"] = eventObjectId
	}
	if dto.UserId != 0 {
		filter["user._id"] = dto.UserId
	}
	if dto.ClubId != 0 {
		filter["club_id"] = dto.ClubId
	}

	filter["user"] = bson.M{"$exists": true}

	find, err := s.invitesCollection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer find.Close(ctx)

	var invites []dao.OrganizerInvite
	err = find.All(ctx, &invites)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainUserInvites(invites), nil
}

func (s *Storage) GetClubInvites(ctx context.Context, dto *dtos.GetInvites) ([]domain.Invite, error) {
	const op = "storage.mongodb.getClubInvites"

	filter := bson.M{}
	if dto.EventId != "" {
		eventObjectId, err := primitive.ObjectIDFromHex(dto.EventId)
		if err != nil {
			if errors.Is(err, primitive.ErrInvalidHex) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		filter["event_id"] = eventObjectId
	}
	if dto.UserId != 0 {
		filter["by_who_id"] = dto.UserId
	}
	if dto.ClubId != 0 {
		filter["club._id"] = dto.ClubId
	}
	filter["club"] = bson.M{"$exists": true}

	find, err := s.invitesCollection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInviteNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer find.Close(ctx)

	var invites []dao.ClubInvite
	err = find.All(ctx, &invites)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dao.ToDomainInvites(invites), nil
}
