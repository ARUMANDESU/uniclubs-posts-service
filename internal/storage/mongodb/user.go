package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.updateUser"

	_, err := s.eventsCollection.UpdateMany(ctx, bson.M{"organizers.id": user.ID}, bson.M{
		"$set": bson.M{
			"organizers.$[elem].first_name": user.FirstName,
			"organizers.$[elem].last_name":  user.LastName,
			"organizers.$[elem].avatar_url": user.AvatarURL,
		},
	}, options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.id": user.ID}},
	}))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
