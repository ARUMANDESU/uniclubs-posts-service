package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.000Z"

type Event struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Club               Club               `json:"club"`
	User               User               `json:"user"`
	CollaboratorClubs  []Club             `json:"collaborator_clubs,omitempty"`
	Organizers         []Organizer        `json:"organizers,omitempty"`
	Title              string             `json:"title,omitempty"`
	Description        string             `json:"description,omitempty"`
	Type               string             `json:"type,omitempty"`
	Status             string             `json:"status,omitempty"`
	Tags               string             `json:"tags,omitempty"`
	MaxParticipants    uint32             `json:"max_participants,omitempty"`
	ParticipantsCount  uint32             `json:"participants_count,omitempty"`
	LocationLink       string             `json:"location_link,omitempty"`
	LocationUniversity string             `json:"location_university,omitempty"`
	StartDate          time.Time          `json:"start_date"`
	EndDate            time.Time          `json:"end_date"`
	CoverImages        []CoverImage       `json:"cover_images,omitempty"`
	AttachedImages     []File             `json:"attached_images,omitempty"`
	AttachedFiles      []File             `json:"attached_files,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	DeletedAt          time.Time          `json:"deleted_at"`
}

type Organizer struct {
	User
	ClubId int64
}

func (o Organizer) ToProto() *eventv1.OrganizerObject {
	return &eventv1.OrganizerObject{
		Id:        o.ID,
		FirstName: o.FirstName,
		LastName:  o.LastName,
		AvatarUrl: o.AvatarURL,
		ClubId:    o.ClubId,
	}
}
func OrganizersToProto(organizers []Organizer) []*eventv1.OrganizerObject {
	convertedOrganizers := make([]*eventv1.OrganizerObject, len(organizers))
	for _, organizer := range organizers {
		convertedOrganizers = append(convertedOrganizers, organizer.ToProto())
	}
	return convertedOrganizers
}

func (e Event) ToProto() *eventv1.EventObject {
	return &eventv1.EventObject{
		Id:                 e.ID.String(),
		Club:               e.Club.ToProto(),
		User:               e.User.ToProto(),
		CollaboratorClubs:  ClubsToProto(e.CollaboratorClubs),
		Organizers:         OrganizersToProto(e.Organizers),
		Title:              e.Title,
		Description:        e.Description,
		Type:               e.Type,
		Tags:               e.Tags,
		MaxParticipants:    e.MaxParticipants,
		ParticipantsCount:  e.ParticipantsCount,
		LocationLink:       e.LocationLink,
		LocationUniversity: e.LocationUniversity,
		StartTime:          e.StartDate.Format(timeLayout),
		EndTime:            e.EndDate.Format(timeLayout),
		CoverImages:        CoverImagesToProto(e.CoverImages),
		AttachedImages:     FilesToProto(e.AttachedImages),
		AttachedFiles:      FilesToProto(e.AttachedFiles),
		CreatedAt:          e.CreatedAt.Format(timeLayout),
		UpdatedAt:          e.UpdatedAt.Format(timeLayout),
		DeletedAt:          e.DeletedAt.Format(timeLayout),
	}
}
