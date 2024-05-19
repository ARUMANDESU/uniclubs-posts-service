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
	client            *mongo.Client
	eventsCollection  *mongo.Collection
	invitesCollection *mongo.Collection
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

	return &Storage{
		client:            client,
		eventsCollection:  eventsCollection,
		invitesCollection: inviteCollection,
	}, nil
}
