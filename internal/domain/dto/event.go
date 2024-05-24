package dtos

import (
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"strings"
	"time"
)

type UpdateEvent struct {
	EventId               string              `json:"event_id"`
	UserId                int64               `json:"user_id"`
	Title                 string              `json:"title"`
	Description           string              `json:"description"`
	Type                  domain.EventType    `json:"type"`
	Tags                  []string            `json:"tags"`
	MaxParticipants       uint32              `json:"max_participants"`
	LocationLink          string              `json:"location_link"`
	LocationUniversity    string              `json:"location_university"`
	StartDate             time.Time           `json:"start_date"`
	EndDate               time.Time           `json:"end_date"`
	CoverImages           []domain.CoverImage `json:"cover_images"`
	AttachedImages        []domain.File       `json:"attached_images"`
	AttachedFiles         []domain.File       `json:"attached_files"`
	IsHiddenForNonMembers bool                `json:"is_hidden_for_non_members"`
	Paths                 []string
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

type DeleteEvent struct {
	EventId string `json:"event_id"`
	UserId  int64  `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
}

func (u *UpdateEvent) HasUnchangeableFields() bool {
	/*
	* unchangeable are fields except
	* tags, max_participants, location_link,
	* location_university, start_date, end_date
	* is_hidden_for_non_members
	 */
	if len(u.Paths) == 0 {
		return true
	}
	for _, path := range u.Paths {
		switch path {
		case "tags":
		case "max_participants":
		case "location_link":
		case "location_university":
		case "start_date":
		case "end_date":
		case "is_hidden_for_non_members":
		default:
			return true
		}
	}

	return false
}

func UpdateToDTO(event *eventv1.UpdateEventRequest) (*UpdateEvent, error) {
	const op = "dtos.UpdateToDTO"

	paths := make(map[string]bool)
	for _, path := range event.GetUpdateMask().GetPaths() {
		paths[path] = true
	}

	var startDate, endDate time.Time
	var err error

	if paths["start_date"] {
		startDate, err = time.Parse(domain.TimeLayout, event.StartDate)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to parse end date: %w", op, err)
		}
	}

	if paths["end_date"] {
		endDate, err = time.Parse(domain.TimeLayout, event.EndDate)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to parse end date: %w", op, err)
		}
	}

	tags := event.GetTags()
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	return &UpdateEvent{
		EventId:               event.GetEventId(),
		UserId:                event.GetUserId(),
		Title:                 strings.Trim(event.GetTitle(), " "),
		Description:           strings.Trim(event.GetDescription(), " "),
		Type:                  domain.EventType(event.GetType()),
		Tags:                  tags,
		MaxParticipants:       uint32(event.GetMaxParticipants()),
		LocationLink:          event.GetLocationLink(),
		LocationUniversity:    event.GetLocationUniversity(),
		StartDate:             startDate,
		EndDate:               endDate,
		CoverImages:           domain.ProtoToCoverImages(event.GetCoverImages()),
		AttachedImages:        domain.ProtoToFiles(event.GetAttachedImages()),
		AttachedFiles:         domain.ProtoToFiles(event.GetAttachedFiles()),
		IsHiddenForNonMembers: event.GetIsHiddenForNonMembers(),
		Paths:                 event.GetUpdateMask().GetPaths(),
	}, nil
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

func DeleteEventToDTO(event *eventv1.DeleteEventRequest) *DeleteEvent {
	return &DeleteEvent{
		EventId: event.GetEventId(),
		UserId:  event.GetUserId(),
		IsAdmin: event.GetIsAdmin(),
	}
}
