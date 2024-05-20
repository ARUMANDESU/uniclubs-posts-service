package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"time"
)

const timeLayout = "2006-01-02 15:04:05.999999999 -0700 MST"

type EventStatus string

const (
	EventStatusDraft     = "DRAFT"
	EventStatusPublished = "PUBLISHED"
	EventStatusCanceled  = "CANCELED"
	EventStatusArchived  = "ARCHIVED"
)

func (s EventStatus) String() string {
	return string(s)
}

type Event struct {
	ID                 string       `json:"id"`
	ClubId             int64        `json:"club_id"`
	OwnerId            int64        `json:"owner_id"`
	CollaboratorClubs  []Club       `json:"collaborator_clubs"`
	Organizers         []Organizer  `json:"organizers"`
	Title              string       `json:"title,omitempty"`
	Description        string       `json:"description,omitempty"`
	Type               string       `json:"type,omitempty"`
	Status             EventStatus  `json:"status,omitempty"`
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

func (e *Event) IsOwner(userId int64) bool {
	return e.OwnerId == userId
}

func (e *Event) IsOrganizer(userId int64) bool {
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

func (e *Event) RemoveOrganizer(organizerId int64) error {
	if len(e.Organizers) == 0 {
		return ErrOrganizersEmpty
	}

	if organizerId == e.OwnerId {
		return ErrUserIsEventOwner
	}

	for i, organizer := range e.Organizers {
		if organizer.ID == organizerId {
			e.Organizers = append(e.Organizers[:i], e.Organizers[i+1:]...)
			return nil
		}
	}

	return ErrOrganizerNotFound
}

func (e *Event) RemoveOrganizersByClubId(clubId int64) error {
	if len(e.Organizers) == 0 {
		return ErrOrganizersEmpty
	}
	var newOrganizers []Organizer
	for _, organizer := range e.Organizers {
		if organizer.ID == e.OwnerId || organizer.ClubId != clubId {
			newOrganizers = append(newOrganizers, organizer)
		}
	}

	e.Organizers = newOrganizers
	return nil
}

func (e *Event) IsCollaborator(clubId int64) bool {
	for _, club := range e.CollaboratorClubs {
		if club.ID == clubId {
			return true
		}
	}
	return false
}

func (e *Event) AddCollaborator(club Club) {
	e.CollaboratorClubs = append(e.CollaboratorClubs, club)
}

func (e *Event) GetCollaboratorById(clubId int64) *Club {
	for _, collaborator := range e.CollaboratorClubs {
		if clubId == collaborator.ID {
			return &collaborator
		}
	}

	return nil
}

func (e *Event) RemoveCollaborator(clubId int64) error {
	if len(e.CollaboratorClubs) == 0 {
		return ErrCollaboratorsEmpty
	}
	if clubId == e.ClubId {
		return ErrClubIsEventOwner
	}

	for i, collaborator := range e.CollaboratorClubs {
		if collaborator.ID == clubId {
			e.CollaboratorClubs = append(e.CollaboratorClubs[:i], e.CollaboratorClubs[i+1:]...)
			return nil
		}
	}

	return ErrCollaboratorNotFound
}

func (e *Event) ToProto() *eventv1.EventObject {
	return &eventv1.EventObject{
		Id:                 e.ID,
		ClubId:             e.ClubId,
		OwnerId:            e.OwnerId,
		CollaboratorClubs:  ClubsToProto(e.CollaboratorClubs),
		Organizers:         OrganizersToProto(e.Organizers),
		Title:              e.Title,
		Description:        e.Description,
		Type:               e.Type,
		Tags:               e.Tags,
		Status:             e.Status.String(),
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

func EventsToProto(events []Event) []*eventv1.EventObject {
	result := make([]*eventv1.EventObject, 0, len(events))
	for _, event := range events {
		result = append(result, event.ToProto())
	}
	return result
}
