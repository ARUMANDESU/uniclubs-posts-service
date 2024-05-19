package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	ID                 primitive.ObjectID `bson:"_id"`
	ClubId             int64              `bson:"club_id"`
	OwnerId            int64              `bson:"owner_id"`
	CollaboratorClubs  []Club             `bson:"collaborator_clubs"`
	Organizers         []Organizer        `bson:"organizers"`
	Title              string             `bson:"title,omitempty"`
	Description        string             `bson:"description,omitempty"`
	Type               string             `bson:"type,omitempty"`
	Status             string             `bson:"status,omitempty"`
	IsApprove          bool               `bson:"is_approve"`
	Tags               []string           `bson:"tags,omitempty"`
	ParticipantIds     []int64            `bson:"participant_ids,omitempty"`
	MaxParticipants    uint32             `bson:"max_participants,omitempty"`
	ParticipantsCount  uint32             `bson:"participants_count,omitempty"`
	LocationLink       string             `bson:"location_link,omitempty"`
	LocationUniversity string             `bson:"location_university,omitempty"`
	StartDate          time.Time          `bson:"start_date,omitempty"`
	EndDate            time.Time          `bson:"end_date,omitempty"`
	CoverImages        []CoverImageMongo  `bson:"cover_images,omitempty"`
	AttachedImages     []FileMongo        `bson:"attached_images,omitempty"`
	AttachedFiles      []FileMongo        `bson:"attached_files,omitempty"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
	DeletedAt          time.Time          `bson:"deleted_at"`
}

func (e *Event) AddOrganizer(organizer Organizer) {
	e.Organizers = append(e.Organizers, organizer)
}

func (e *Event) AddCollaboratorClub(club Club) {
	e.CollaboratorClubs = append(e.CollaboratorClubs, club)
}
