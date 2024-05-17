package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.000Z"

const (
	EventStatusDraft     = "DRAFT"
	EventStatusPublished = "PUBLISHED"
	EventStatusCanceled  = "CANCELED"
	EventStatusArchived  = "ARCHIVED"
)

type Event struct {
	ID                 string       `json:"id"`
	Club               Club         `json:"club"`
	User               User         `json:"user"`
	CollaboratorClubs  []Club       `json:"collaborator_clubs,omitempty"`
	Organizers         []Organizer  `json:"organizers,omitempty"`
	Title              string       `json:"title,omitempty"`
	Description        string       `json:"description,omitempty"`
	Type               string       `json:"type,omitempty"`
	Status             string       `json:"status,omitempty"`
	Tags               []string     `json:"tags,omitempty"`
	MaxParticipants    uint32       `json:"max_participants,omitempty"`
	ParticipantsCount  uint32       `json:"participants_count,omitempty"`
	LocationLink       string       `json:"location_link,omitempty"`
	LocationUniversity string       `json:"location_university,omitempty"`
	StartDate          time.Time    `json:"start_date"`
	EndDate            time.Time    `json:"end_date"`
	CoverImages        []CoverImage `json:"cover_images,omitempty"`
	AttachedImages     []File       `json:"attached_images,omitempty"`
	AttachedFiles      []File       `json:"attached_files,omitempty"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	DeletedAt          time.Time    `json:"deleted_at"`
}

func (e *Event) IsOrganizer(userId int64) bool {
	if e.User.ID == userId {
		return true
	}

	for _, organizer := range e.Organizers {
		if organizer.ID == userId {
			return true
		}
	}

	return false
}

func (e *Event) GetOrganizerById(userId int64) *Organizer {
	for _, organizer := range e.Organizers {
		if organizer.ID == userId {
			return &organizer
		}
	}

	return nil
}

func (e *Event) AddOrganizer(organizer Organizer) {
	e.Organizers = append(e.Organizers, organizer)
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
	for i, organizer := range organizers {
		convertedOrganizers[i] = organizer.ToProto()
	}
	return convertedOrganizers
}

func (e *Event) ToProto() *eventv1.EventObject {
	return &eventv1.EventObject{
		Id:                 e.ID,
		Club:               e.Club.ToProto(),
		User:               e.User.ToProto(),
		CollaboratorClubs:  ClubsToProto(e.CollaboratorClubs),
		Organizers:         OrganizersToProto(e.Organizers),
		Title:              e.Title,
		Description:        e.Description,
		Type:               e.Type,
		Tags:               e.Tags,
		Status:             e.Status,
		MaxParticipants:    e.MaxParticipants,
		ParticipantsCount:  e.ParticipantsCount,
		LocationLink:       e.LocationLink,
		LocationUniversity: e.LocationUniversity,
		StartDate:          e.StartDate.Format(timeLayout),
		EndDate:            e.EndDate.Format(timeLayout),
		CoverImages:        CoverImagesToProto(e.CoverImages),
		AttachedImages:     FilesToProto(e.AttachedImages),
		AttachedFiles:      FilesToProto(e.AttachedFiles),
		CreatedAt:          e.CreatedAt.Format(timeLayout),
		UpdatedAt:          e.UpdatedAt.Format(timeLayout),
		DeletedAt:          e.DeletedAt.Format(timeLayout),
	}
}
