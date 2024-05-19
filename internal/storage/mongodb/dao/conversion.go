package dao

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDomainEvent(
	e Event,
) *domain.Event {
	collaboratorClubs := ToDomainClubs(e.CollaboratorClubs)
	organizers := ToDomainOrganizers(e.Organizers)

	return &domain.Event{
		ID:                 e.ID.Hex(),
		ClubId:             e.ClubId,
		OwnerId:            e.OwnerId,
		CollaboratorClubs:  collaboratorClubs,
		Organizers:         organizers,
		Title:              e.Title,
		Description:        e.Description,
		Type:               e.Type,
		Status:             e.Status,
		Tags:               e.Tags,
		MaxParticipants:    e.MaxParticipants,
		ParticipantsCount:  e.ParticipantsCount,
		LocationLink:       e.LocationLink,
		LocationUniversity: e.LocationUniversity,
		StartDate:          e.StartDate,
		EndDate:            e.EndDate,
		CoverImages:        ToDomainCoverImages(e.CoverImages),
		AttachedImages:     ToDomainFiles(e.AttachedImages),
		AttachedFiles:      ToDomainFiles(e.AttachedFiles),
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
		DeletedAt:          e.DeletedAt,
	}
}

func EventToModel(event *domain.Event) Event {
	objectID, _ := primitive.ObjectIDFromHex(event.ID)

	return Event{
		ID:                 objectID,
		ClubId:             event.ClubId,
		OwnerId:            event.OwnerId,
		CollaboratorClubs:  ToCollaboratorClubs(event.CollaboratorClubs),
		Organizers:         ToOrganizers(event.Organizers),
		Title:              event.Title,
		Description:        event.Description,
		Type:               event.Type,
		Status:             event.Status,
		Tags:               event.Tags,
		MaxParticipants:    event.MaxParticipants,
		ParticipantsCount:  event.ParticipantsCount,
		LocationLink:       event.LocationLink,
		LocationUniversity: event.LocationUniversity,
		StartDate:          event.StartDate,
		EndDate:            event.EndDate,
		CoverImages:        ToCoverImages(event.CoverImages),
		AttachedImages:     ToFiles(event.AttachedImages),
		AttachedFiles:      ToFiles(event.AttachedFiles),
		CreatedAt:          event.CreatedAt,
		UpdatedAt:          event.UpdatedAt,
		DeletedAt:          event.DeletedAt,
	}
}
