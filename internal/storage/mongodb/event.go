package mongodb

import (
	"context"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	ID                  primitive.ObjectID `bson:"_id"`
	ClubId              int64              `bson:"club_id"`
	UserId              int64              `bson:"user_id"`
	CollaboratorClubIds []int64            `bson:"collaborator_club_ids,omitempty"`
	OrganizerIds        []int64            `bson:"organizer_ids,omitempty"`
	Title               string             `bson:"title,omitempty"`
	Description         string             `bson:"description,omitempty"`
	Type                string             `bson:"type,omitempty"`
	Status              string             `bson:"status,omitempty"`
	Tags                []string           `bson:"tags,omitempty"`
	ParticipantIds      []int64            `bson:"participant_ids,omitempty"`
	MaxParticipants     uint32             `bson:"max_participants,omitempty"`
	ParticipantsCount   uint32             `bson:"participants_count,omitempty"`
	LocationLink        string             `bson:"location_link,omitempty"`
	LocationUniversity  string             `bson:"location_university,omitempty"`
	StartDate           time.Time          `bson:"start_date,omitempty"`
	EndDate             time.Time          `bson:"end_date,omitempty"`
	CoverImages         []CoverImageMongo  `bson:"cover_images,omitempty"`
	AttachedImages      []FileMongo        `bson:"attached_images,omitempty"`
	AttachedFiles       []FileMongo        `bson:"attached_files,omitempty"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at"`
	DeletedAt           time.Time          `bson:"deleted_at"`
}

type FileMongo struct {
	URL  string `bson:"url"`
	Name string `bson:"name"`
	Type string `bson:"type"`
}

type CoverImageMongo struct {
	FileMongo
	Position uint32 `bson:"position"`
}

func (s Storage) CreateEvent(ctx context.Context, clubId int64, userId int64) (*domain.Event, error) {
	const op = "storage.mongodb.event.createEvent"

	input := Event{
		ID:        primitive.NewObjectID(),
		ClubId:    clubId,
		UserId:    userId,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Status:    "DRAFT",
	}

	insertResult, err := s.eventsCollection.InsertOne(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if insertResult.InsertedID == nil {
		return nil, fmt.Errorf("%s: no inserted id", op)
	}

	insertedID := insertResult.InsertedID.(primitive.ObjectID)
	var event domain.Event
	err = s.eventsCollection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &event, nil
}

func (s Storage) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	//TODO implement me
	panic("implement me")
}
