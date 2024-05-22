package dtos

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"strings"
)

type UpdateEvent struct {
	EventId            string              `json:"event_id"`
	UserId             int64               `json:"user_id"`
	Title              string              `json:"title"`
	Description        string              `json:"description"`
	Type               domain.EventType    `json:"type"`
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

type SendJoinRequestToUser struct {
	EventId      string      `json:"event_id"`
	UserId       int64       `json:"user_id"`
	Target       domain.User `json:"target"`
	TargetClubId int64       `json:"target_club_id"`
}

type SendJoinRequestToClub struct {
	EventId string      `json:"event_id"`
	UserId  int64       `json:"user_id"`
	Club    domain.Club `json:"club"`
}

type AcceptJoinRequestClub struct {
	InviteId string      `json:"invite_id"`
	ClubId   int64       `json:"club_id"`
	User     domain.User `json:"user"`
}

type RejectEvent struct {
	EventId string      `json:"event_id"`
	User    domain.User `json:"user"`
	Reason  string      `json:"reason"`
}

func UpdateToDTO(event *eventv1.UpdateEventRequest) *UpdateEvent {
	tags := event.GetTags()
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return &UpdateEvent{
		EventId:            event.GetEventId(),
		UserId:             event.GetUserId(),
		Title:              strings.Trim(event.GetTitle(), " "),
		Description:        strings.Trim(event.GetDescription(), " "),
		Type:               domain.EventType(event.GetType()),
		Tags:               tags,
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

func AddOrganizerRequestToUserToDTO(event *eventv1.AddOrganizerRequest) *SendJoinRequestToUser {
	return &SendJoinRequestToUser{
		EventId:      event.GetEventId(),
		UserId:       event.GetUserId(),
		Target:       domain.UserFromProto(event.GetTarget()),
		TargetClubId: event.GetTargetClubId(),
	}
}

func AddCollaboratorRequestToClubToDTO(event *eventv1.AddCollaboratorRequest) *SendJoinRequestToClub {
	return &SendJoinRequestToClub{
		EventId: event.GetEventId(),
		UserId:  event.GetUserId(),
		Club:    domain.ClubFromProto(event.GetClub()),
	}
}

func AcceptJoinRequestClubToDTO(event *eventv1.HandleInviteClubRequest) *AcceptJoinRequestClub {
	return &AcceptJoinRequestClub{
		InviteId: event.GetInviteId(),
		ClubId:   event.GetClubId(),
		User:     domain.UserFromProto(event.GetUser()),
	}
}

func RejectEventToDTO(event *eventv1.RejectEventRequest) *RejectEvent {
	return &RejectEvent{
		EventId: event.GetEventId(),
		User:    domain.UserFromProto(event.GetUser()),
		Reason:  event.GetReason(),
	}
}
