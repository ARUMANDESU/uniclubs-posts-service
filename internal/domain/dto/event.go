package dto

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"time"
)

const timeLayout = "2006-01-02T15:04:05Z07:00"

type CreateEventDTO struct {
	ClubId             int64                `json:"club_id" bson:"club_id"`
	UserId             int64                `json:"user_id" bson:"user_id"`
	CollaboratorClubs  []int64              `json:"collaborator_clubs" bson:"collaborator_clubs,omitempty"`
	Title              string               `json:"title" bson:"title,omitempty"`
	Description        string               `json:"description" bson:"description,omitempty"`
	Type               string               `json:"type" bson:"type,omitempty"`
	Tags               string               `json:"tags" bson:"tags,omitempty"`
	MaxParticipants    uint32               `json:"max_participants" bson:"max_participants,omitempty"`
	StartDate          time.Time            `json:"start_date" bson:"start_date,omitempty"`
	EndDate            time.Time            `json:"end_date" bson:"end_date,omitempty"`
	LocationLink       string               `json:"location_link" bson:"location_link,omitempty"`
	LocationUniversity string               `json:"location_university" bson:"location_university,omitempty"`
	CoverImages        []*domain.CoverImage `json:"cover_images" bson:"cover_images,omitempty"`
	AttachedImages     []*domain.File       `json:"attached_images" bson:"attached_images,omitempty"`
	AttachedFiles      []*domain.File       `json:"attached_files" bson:"attached_files,omitempty"`
	CreatedAt          time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at" bson:"updated_at"`
}

type GetEventDTO struct {
	ID     int64 `json:"id" bson:"id"`
	UserId int64 `json:"user_id" bson:"user_id"`
}

func ProtoToCreateEventDTO(proto *eventv1.CreateEventRequest) *CreateEventDTO {
	startTime, _ := time.Parse(timeLayout, proto.GetStartTime())
	endTime, _ := time.Parse(timeLayout, proto.GetEndTime())

	return &CreateEventDTO{
		ClubId:             proto.GetClubId(),
		UserId:             proto.GetUserId(),
		CollaboratorClubs:  proto.GetCollaboratorClubs(),
		Title:              proto.GetTitle(),
		Description:        proto.GetDescription(),
		Type:               proto.GetType(),
		Tags:               proto.GetTags(),
		MaxParticipants:    uint32(proto.GetMaxParticipants()),
		StartDate:          startTime,
		EndDate:            endTime,
		LocationLink:       proto.GetLocationLink(),
		LocationUniversity: proto.GetLocationUniversity(),
		CoverImages:        domain.ProtoToCoverImages(proto.GetCoverImages()),
		AttachedImages:     domain.ProtoToFiles(proto.GetAttachedImages()),
		AttachedFiles:      domain.ProtoToFiles(proto.GetAttachedFiles()),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
