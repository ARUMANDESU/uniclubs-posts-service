package dto

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
)

type UpdateEvent struct {
	EventId            string              `json:"event_id"`
	UserId             int64               `json:"user_id"`
	Title              string              `json:"title"`
	Description        string              `json:"description"`
	Type               string              `json:"type"`
	Tags               []string            `json:"tags"`
	MaxParticipants    uint32              `json:"max_participants"`
	LocationLink       string              `json:"location_link"`
	LocationUniversity string              `json:"location_university"`
	StartDate          string              `json:"start_date"`
	EndDate            string              `json:"end_date"`
	CoverImages        []domain.CoverImage `json:"cover_images"`
	AttachedImages     []domain.File       `json:"attached_images"`
	AttachedFiles      []domain.File       `json:"attached_files"`
	Paths              []string
}

func UpdateToDTO(event *eventv1.UpdateEventRequest) *UpdateEvent {
	return &UpdateEvent{
		EventId:            event.GetEventId(),
		UserId:             event.GetUserId(),
		Title:              event.GetTitle(),
		Description:        event.GetDescription(),
		Type:               event.GetType(),
		Tags:               event.GetTags(),
		MaxParticipants:    uint32(event.GetMaxParticipants()),
		LocationLink:       event.GetLocationLink(),
		LocationUniversity: event.GetLocationUniversity(),
		StartDate:          event.GetStartDate(),
		EndDate:            event.GetEndDate(),
		CoverImages:        domain.ProtoToCoverImages(event.GetCoverImages()),
		AttachedImages:     domain.ProtoToFiles(event.GetAttachedImages()),
		AttachedFiles:      domain.ProtoToFiles(event.GetAttachedFiles()),
		Paths:              event.GetUpdateMask().GetPaths(),
	}
}
