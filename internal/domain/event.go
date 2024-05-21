package domain

import (
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05.999999999 -0700 MST"

type EventStatus string
type EventType string

const (
	EventStatusDraft      EventStatus = "DRAFT"
	EventStatusPending    EventStatus = "PENDING"
	EventStatusApproved   EventStatus = "APPROVED"
	EventStatusRejected   EventStatus = "REJECTED"
	EventStatusInProgress EventStatus = "IN_PROGRESS"
	EventStatusFinished   EventStatus = "FINISHED"
	EventStatusCanceled   EventStatus = "CANCELED"
	EventStatusArchived   EventStatus = "ARCHIVED"

	EventTypeUniversity EventType = "UNIVERSITY"
	EventTypeIntraClub  EventType = "INTRA_CLUB"
)

func (s EventStatus) String() string {
	return string(s)
}

func (t EventType) String() string {
	return string(t)
}

type Event struct {
	ID                 string          `json:"id"`
	ClubId             int64           `json:"club_id"`
	OwnerId            int64           `json:"owner_id"`
	CollaboratorClubs  []Club          `json:"collaborator_clubs"`
	Organizers         []Organizer     `json:"organizers"`
	Title              string          `json:"title,omitempty"`
	Description        string          `json:"description,omitempty"`
	Type               EventType       `json:"type,omitempty"`
	Status             EventStatus     `json:"status,omitempty"`
	Tags               []string        `json:"tags,omitempty"`
	MaxParticipants    uint32          `json:"max_participants,omitempty"`
	ParticipantsCount  uint32          `json:"participants_count,omitempty"`
	LocationLink       string          `json:"location_link,omitempty"`
	LocationUniversity string          `json:"location_university,omitempty"`
	StartDate          time.Time       `json:"start_date"`
	EndDate            time.Time       `json:"end_date"`
	CoverImages        []CoverImage    `json:"cover_images,omitempty"`
	AttachedImages     []File          `json:"attached_images,omitempty"`
	AttachedFiles      []File          `json:"attached_files,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          time.Time       `json:"deleted_at"`
	PublishedAt        time.Time       `json:"published_at"`
	ApproveMetadata    ApproveMetadata `json:"approve_metadata"`
	RejectMetadata     RejectMetadata  `json:"reject_metadata"`
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

func (e *Event) ChangeStatus(status EventStatus) {
	e.Status = status
}

func (e *Event) canPublish() error {
	if e.Status == EventStatusApproved {
		return nil
	}
	if e.Type == EventTypeIntraClub {
		return nil
	}

	return ErrEventIsNotApproved
}

func (e *Event) Publish() error {
	if err := e.canPublish(); err != nil {
		return err
	}

	e.ChangeStatus(EventStatusInProgress)
	e.PublishedAt = time.Now()
	return nil
}

func (e *Event) SendToReview() error {
	statusErrors := map[EventStatus]error{
		EventStatusPending:    fmt.Errorf("event already in review status"),
		EventStatusInProgress: fmt.Errorf("event already in progress"),
		EventStatusApproved:   fmt.Errorf("event already approved"),
		EventStatusArchived:   fmt.Errorf("event archived"),
		EventStatusCanceled:   fmt.Errorf("event canceled"),
		EventStatusFinished:   fmt.Errorf("event finished"),
	}

	if e.Type == EventTypeIntraClub {
		return fmt.Errorf("intra club events do not need review")
	}

	if err, ok := statusErrors[e.Status]; ok {
		return err
	}

	e.ChangeStatus(EventStatusPending)
	return nil
}

func (e *Event) RevokeReview() error {
	if e.Status != EventStatusPending {
		return fmt.Errorf("event is not in review status")
	}

	e.ChangeStatus(EventStatusDraft)
	return nil
}

func (e *Event) Approve(user User) error {
	if e.Status != EventStatusPending {
		return fmt.Errorf("event is not in review status")
	}

	e.ChangeStatus(EventStatusApproved)
	e.ApproveMetadata = ApproveMetadata{
		ApprovedBy: user,
		ApprovedAt: time.Now(),
	}
	return nil
}

func (e *Event) Reject(user User, reason string) error {
	if e.Status != EventStatusPending {
		return fmt.Errorf("event is not in review status")
	}

	e.ChangeStatus(EventStatusRejected)
	e.RejectMetadata = RejectMetadata{
		RejectedBy: user,
		RejectedAt: time.Now(),
		Reason:     reason,
	}
	return nil
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
		Type:               e.Type.String(),
		Tags:               e.Tags,
		Status:             e.Status.String(),
		MaxParticipants:    e.MaxParticipants,
		ParticipantsCount:  e.ParticipantsCount,
		LocationLink:       e.LocationLink,
		LocationUniversity: e.LocationUniversity,
		StartDate:          e.StartDate.Format(TimeLayout),
		EndDate:            e.EndDate.Format(TimeLayout),
		CoverImages:        CoverImagesToProto(e.CoverImages),
		AttachedImages:     FilesToProto(e.AttachedImages),
		AttachedFiles:      FilesToProto(e.AttachedFiles),
		CreatedAt:          e.CreatedAt.Format(TimeLayout),
		UpdatedAt:          e.UpdatedAt.Format(TimeLayout),
		DeletedAt:          e.DeletedAt.Format(TimeLayout),
		//PublishedAt:        e.PublishedAt.Format(TimeLayout),
		ApproveMetadata: e.ApproveMetadata.ToProto(),
		RejectMetadata:  e.RejectMetadata.ToProto(),
	}
}

func EventsToProto(events []Event) []*eventv1.EventObject {
	result := make([]*eventv1.EventObject, 0, len(events))
	for _, event := range events {
		result = append(result, event.ToProto())
	}
	return result
}
