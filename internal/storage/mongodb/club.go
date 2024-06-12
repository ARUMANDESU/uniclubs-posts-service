package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

func (s *Storage) UpdateClub(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.updateClub"

	updateFuncs := []func(context.Context, *domain.Club) error{
		s.updateClubsInEventsCollection,
		s.updateClubInInviteCollection,
		s.updateClubInPostsCollection,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	for _, updateFunc := range updateFuncs {
		wg.Add(1)
		go func(updateFunc func(context.Context, *domain.Club) error) {
			defer wg.Done()
			if err := updateFunc(ctx, club); err != nil {
				errChan <- fmt.Errorf("%s: %w", op, err)
			}
		}(updateFunc)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (s *Storage) updateClubsInEventsCollection(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.updateClubsInEventsCollection"

	_, err := s.eventsCollection.UpdateMany(ctx, bson.M{"collaborator_clubs._id": club.ID}, bson.M{
		"$set": bson.M{
			"collaborator_clubs.$[elem].name":     club.Name,
			"collaborator_clubs.$[elem].logo_url": club.LogoURL,
		},
	}, options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem._id": club.ID}},
	}))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) updateClubInInviteCollection(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.updateClubInInviteCollection"

	_, err := s.invitesCollection.UpdateMany(ctx, bson.M{"club._id": club.ID}, bson.M{
		"$set": bson.M{
			"club.name":     club.Name,
			"club.logo_url": club.LogoURL,
		},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) updateClubInPostsCollection(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.updateClubInPostsCollection"

	_, err := s.postsCollection.UpdateMany(ctx, bson.M{"club._id": club.ID}, bson.M{
		"$set": bson.M{
			"club.name":     club.Name,
			"club.logo_url": club.LogoURL,
		},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
