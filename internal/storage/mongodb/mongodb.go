package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage struct {
	client                 *mongo.Client
	eventsCollection       *mongo.Collection
	invitesCollection      *mongo.Collection
	participantsCollection *mongo.Collection
	bansCollection         *mongo.Collection
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
	eventsCollection := db.Collection("events")
	inviteCollection := db.Collection("invites")
	participantsCollection := db.Collection("participants")
	bansCollection := db.Collection("bans")

	// Create text index on the 'name' field
	eventIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "description", Value: "text"},
			{Key: "tags", Value: "text"},
		},
	}
	_, err = eventsCollection.Indexes().CreateOne(ctx, eventIndex)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	participantsIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "event_id", Value: 1},
			{Key: "user._id", Value: 1},
			{Key: "user.first_name", Value: "text"},
			{Key: "user.last_name", Value: "text"},
		},
	}
	_, err = participantsCollection.Indexes().CreateOne(ctx, participantsIndex)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		client:                 client,
		eventsCollection:       eventsCollection,
		invitesCollection:      inviteCollection,
		participantsCollection: participantsCollection,
		bansCollection:         bansCollection,
	}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	const op = "storage.mongodb.close"

	if err := s.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
