package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/config"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage struct {
	client           *mongo.Client
	userCollection   *mongo.Collection
	clubsCollection  *mongo.Collection
	eventsCollection *mongo.Collection
}

func New(ctx context.Context, cfg config.MongoDB) (*Storage, error) {
	const op = "storage.mongodb.new"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()
	if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := client.Database(cfg.DatabaseName)
	usersCollection := db.Collection("users")
	clubsCollection := db.Collection("clubs")
	eventsCollection := db.Collection("events")

	return &Storage{
		client:           client,
		userCollection:   usersCollection,
		clubsCollection:  clubsCollection,
		eventsCollection: eventsCollection,
	}, nil
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
