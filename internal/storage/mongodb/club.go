package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s Storage) UpdateClub(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.updateClub"

	_, err := s.eventsCollection.UpdateMany(ctx, bson.M{"collaborator_clubs.id": club.ID}, bson.M{
		"$set": bson.M{
			"collaborator_clubs.$[elem].name":     club.Name,
			"collaborator_clubs.$[elem].logo_url": club.LogoURL,
		},
	}, options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.id": club.ID}},
	}))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
