package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage struct {
	client           *mongo.Client
	userCollection   *mongo.Collection
	clubsCollection  *mongo.Collection
	eventsCollection *mongo.Collection
	inviteCollection *mongo.Collection
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
	inviteCollection := db.Collection("invites")

	return &Storage{
		client:           client,
		userCollection:   usersCollection,
		clubsCollection:  clubsCollection,
		eventsCollection: eventsCollection,
		inviteCollection: inviteCollection,
	}, nil
}
