package dao

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	ID                    primitive.ObjectID `bson:"_id"`
	ClubId                int64              `bson:"club_id"`
	OwnerId               int64              `bson:"owner_id"`
	CollaboratorClubs     []Club             `bson:"collaborator_clubs"`
	Organizers            []Organizer        `bson:"organizers"`
	Title                 string             `bson:"title,omitempty"`
	Description           string             `bson:"description,omitempty"`
	Type                  string             `bson:"type,omitempty"`
	Status                string             `bson:"status,omitempty"`
	Tags                  []string           `bson:"tags,omitempty"`
	ParticipantIds        []int64            `bson:"participant_ids,omitempty"`
	MaxParticipants       uint32             `bson:"max_participants,omitempty"`
	ParticipantsCount     uint32             `bson:"participants_count,minsize"`
	LocationLink          string             `bson:"location_link,omitempty"`
	LocationUniversity    string             `bson:"location_university,omitempty"`
	StartDate             time.Time          `bson:"start_date,omitempty"`
	EndDate               time.Time          `bson:"end_date,omitempty"`
	CoverImages           []CoverImageMongo  `bson:"cover_images,omitempty"`
	AttachedImages        []FileMongo        `bson:"attached_images,omitempty"`
	AttachedFiles         []FileMongo        `bson:"attached_files,omitempty"`
	CreatedAt             time.Time          `bson:"created_at"`
	UpdatedAt             time.Time          `bson:"updated_at"`
	DeletedAt             time.Time          `bson:"deleted_at,omitempty"`
	PublishedAt           time.Time          `bson:"published_at,omitempty"`
	ApproveMetadata       ApproveMetadata    `json:"approve_metadata,omitempty"`
	RejectMetadata        RejectMetadata     `json:"reject_metadata,omitempty"`
	IsHiddenForNonMembers bool               `bson:"is_hidden_for_non_members"`
}

func (e *Event) AddOrganizer(organizer Organizer) {
	e.Organizers = append(e.Organizers, organizer)
}

func (e *Event) AddCollaboratorClub(club Club) {
	e.CollaboratorClubs = append(e.CollaboratorClubs, club)
}

func ToDomainEvent(
	e Event,
) *domain.Event {
	collaboratorClubs := ToDomainClubs(e.CollaboratorClubs)
	organizers := ToDomainOrganizers(e.Organizers)

	return &domain.Event{
		ID:                    e.ID.Hex(),
		ClubId:                e.ClubId,
		OwnerId:               e.OwnerId,
		CollaboratorClubs:     collaboratorClubs,
		Organizers:            organizers,
		Title:                 e.Title,
		Description:           e.Description,
		Type:                  domain.EventType(e.Type),
		Status:                domain.EventStatus(e.Status),
		Tags:                  e.Tags,
		MaxParticipants:       e.MaxParticipants,
		ParticipantsCount:     e.ParticipantsCount,
		LocationLink:          e.LocationLink,
		LocationUniversity:    e.LocationUniversity,
		StartDate:             e.StartDate,
		EndDate:               e.EndDate,
		CoverImages:           ToDomainCoverImages(e.CoverImages),
		AttachedImages:        ToDomainFiles(e.AttachedImages),
		AttachedFiles:         ToDomainFiles(e.AttachedFiles),
		CreatedAt:             e.CreatedAt,
		UpdatedAt:             e.UpdatedAt,
		DeletedAt:             e.DeletedAt,
		PublishedAt:           e.PublishedAt,
		ApproveMetadata:       e.ApproveMetadata.ToDomain(),
		RejectMetadata:        e.RejectMetadata.ToDomain(),
		IsHiddenForNonMembers: e.IsHiddenForNonMembers,
	}
}

func EventToModel(event *domain.Event) Event {
	objectID, _ := primitive.ObjectIDFromHex(event.ID)

	return Event{
		ID:                    objectID,
		ClubId:                event.ClubId,
		OwnerId:               event.OwnerId,
		CollaboratorClubs:     ToCollaboratorClubs(event.CollaboratorClubs),
		Organizers:            ToOrganizers(event.Organizers),
		Title:                 event.Title,
		Description:           event.Description,
		Type:                  event.Type.String(),
		Status:                event.Status.String(),
		Tags:                  event.Tags,
		MaxParticipants:       event.MaxParticipants,
		ParticipantsCount:     event.ParticipantsCount,
		LocationLink:          event.LocationLink,
		LocationUniversity:    event.LocationUniversity,
		StartDate:             event.StartDate,
		EndDate:               event.EndDate,
		CoverImages:           ToCoverImages(event.CoverImages),
		AttachedImages:        ToFiles(event.AttachedImages),
		AttachedFiles:         ToFiles(event.AttachedFiles),
		CreatedAt:             event.CreatedAt,
		UpdatedAt:             event.UpdatedAt,
		DeletedAt:             event.DeletedAt,
		PublishedAt:           event.PublishedAt,
		ApproveMetadata:       ToApproveMetadata(event.ApproveMetadata),
		RejectMetadata:        ToRejectMetadata(event.RejectMetadata),
		IsHiddenForNonMembers: event.IsHiddenForNonMembers,
	}
}

func ToDomainEvents(events []Event) []domain.Event {
	models := make([]domain.Event, 0, len(events))
	for _, event := range events {
		models = append(models, *ToDomainEvent(event))
	}

	return models
}
