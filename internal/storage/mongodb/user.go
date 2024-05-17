package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID        int64  `json:"id,omitempty" bson:"_id"`
	FirstName string `json:"first_name,omitempty" bson:"first_name"`
	LastName  string `json:"last_name,omitempty" bson:"last_name"`
	Barcode   string `json:"barcode,omitempty" bson:"barcode"`
	AvatarURL string `json:"avatar_url,omitempty" bson:"avatar_url"`
}

type Organizer struct {
	User
	ClubId int64 `json:"club_id" bson:"club_id"`
}

func (s Storage) SaveUser(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.saveUser"

	_, err := s.userCollection.InsertOne(ctx, user)
	if err != nil {

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s Storage) GetUserByID(ctx context.Context, userID int64) (user *domain.User, err error) {
	const op = "storage.mongodb.getUserByID"

	err = s.userCollection.FindOne(ctx, bson.D{{"_id", userID}}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	const op = "storage.mongodb.updateUser"

	_, err := s.userCollection.ReplaceOne(ctx, bson.D{{"_id", user.ID}}, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s Storage) DeleteUserByID(ctx context.Context, userID int64) error {
	const op = "storage.mongodb.deleteUserByID"

	_, err := s.userCollection.DeleteOne(ctx, bson.D{{"_id", userID}})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
