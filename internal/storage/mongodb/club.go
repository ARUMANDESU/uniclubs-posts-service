package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type Club struct {
	ID      int64  `json:"id,omitempty" bson:"_id"`
	Name    string `json:"name,omitempty" bson:"name"`
	LogoURL string `json:"logo_url,omitempty" bson:"logo_url"`
}

func (s Storage) SaveClub(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.saveClub"

	_, err := s.clubsCollection.InsertOne(ctx, club)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s Storage) UpdateClub(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.updateClub"

	_, err := s.clubsCollection.ReplaceOne(ctx, bson.D{{"_id", club.ID}}, club)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s Storage) GetClubByID(ctx context.Context, clubID int64) (club *domain.Club, err error) {
	const op = "storage.mongodb.getClubByID"

	err = s.clubsCollection.FindOne(ctx, bson.D{{"_id", clubID}}).Decode(&club)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return club, nil
}

func (s Storage) DeleteClub(ctx context.Context, club *domain.Club) error {
	const op = "storage.mongodb.deleteClub"

	_, err := s.clubsCollection.DeleteOne(ctx, bson.D{{"_id", club.ID}})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
