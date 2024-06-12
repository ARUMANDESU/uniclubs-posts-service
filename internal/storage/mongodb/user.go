package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.updateUser"

	updateFuncs := []func(context.Context, *domain.User) error{
		s.updateUserInEventCollection,
		s.updateUserInInviteCollection,
		s.updateUserInParticipantCollection,
		s.updateUserInBansCollection,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	for _, updateFunc := range updateFuncs {
		wg.Add(1)
		go func(updateFunc func(context.Context, *domain.User) error) {
			defer wg.Done()
			if err := updateFunc(ctx, user); err != nil {
				errChan <- fmt.Errorf("%s: %w", op, err)
			}
		}(updateFunc)
	}

	// Wait for both updates to complete
	wg.Wait()
	close(errChan)

	// Check for errors from goroutines
	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (s *Storage) updateUserInEventCollection(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.updateUserInEventCollection"

	_, err := s.eventsCollection.UpdateMany(ctx, bson.M{"organizers._id": user.ID}, bson.M{
		"$set": bson.M{
			"organizers.$[elem].first_name": user.FirstName,
			"organizers.$[elem].last_name":  user.LastName,
			"organizers.$[elem].barcode":    user.Barcode,
			"organizers.$[elem].avatar_url": user.AvatarURL,
		},
	}, options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem._id": user.ID}},
	}))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) updateUserInInviteCollection(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.updateUserInInviteCollection"

	_, err := s.invitesCollection.UpdateMany(ctx, bson.M{"user._id": user.ID}, bson.M{
		"$set": bson.M{
			"user.first_name": user.FirstName,
			"user.last_name":  user.LastName,
			"user.barcode":    user.Barcode,
			"user.avatar_url": user.AvatarURL,
		},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) updateUserInParticipantCollection(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.updateUserInParticipantCollection"

	_, err := s.participantsCollection.UpdateMany(ctx, bson.M{"user._id": user.ID}, bson.M{
		"$set": bson.M{
			"user.first_name": user.FirstName,
			"user.last_name":  user.LastName,
			"user.barcode":    user.Barcode,
			"user.avatar_url": user.AvatarURL,
		},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) updateUserInBansCollection(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.updateUserInBansCollection"

	_, err := s.bansCollection.UpdateMany(ctx, bson.M{"user._id": user.ID}, bson.M{
		"$set": bson.M{
			"user.first_name": user.FirstName,
			"user.last_name":  user.LastName,
			"user.barcode":    user.Barcode,
			"user.avatar_url": user.AvatarURL,
		},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
