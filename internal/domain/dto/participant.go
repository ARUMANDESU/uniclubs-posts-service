package dtos

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
)

type KickParticipant struct {
	EventId       string `json:"event_id"`
	UserId        int64  `json:"user_id"`
	ParticipantId int64  `json:"participant_id"`
}

func ProtoToKickParticipant(req *eventv1.KickParticipantRequest) *KickParticipant {
	return &KickParticipant{
		EventId:       req.GetEventId(),
		UserId:        req.GetUserId(),
		ParticipantId: req.GetParticipantId(),
	}
}

type BanParticipant struct {
	EventId     string      `json:"event_id"`
	UserId      int64       `json:"user_id"`
	Participant domain.User `json:"participant"`
	Reason      string      `json:"reason"`
}

func ProtoToBanParticipant(req *eventv1.BanParticipantRequest) *BanParticipant {
	return &BanParticipant{
		EventId:     req.GetEventId(),
		UserId:      req.GetUserId(),
		Participant: domain.User{ID: req.GetParticipantId()},
		Reason:      req.GetReason(),
	}
}

type ListParticipants struct {
	EventId string            `json:"event_id"`
	Filter  domain.BaseFilter `json:"filter"`
}

func ProtoToListParticipants(req *eventv1.ListParticipantsRequest) *ListParticipants {
	dto := &ListParticipants{
		EventId: req.GetEventId(),
		Filter: domain.BaseFilter{
			Page:     req.GetPageNumber(),
			PageSize: req.GetPageSize(),
			Query:    req.GetQuery(),
			SortBy:   domain.SortBy(req.GetSortBy()),
		},
	}

	if req.GetSortOrder() == "" {
		dto.Filter.SortOrder = domain.SortOrderDesc
	} else {
		dto.Filter.SortOrder = domain.SortOrder(req.GetSortOrder())
	}

	if req.GetPageNumber() == 0 {
		dto.Filter.Page = 1
	}

	if req.GetPageSize() == 0 {
		dto.Filter.PageSize = 10
	}

	return dto
}

type ListBans struct {
	EventId string            `json:"event_id"`
	Filter  domain.BaseFilter `json:"filter"`
	UserId  int64             `json:"user_id"`
}

func ProtoToListBans(req *eventv1.ListBannedParticipantsRequest) *ListBans {
	dto := &ListBans{
		EventId: req.GetEventId(),
		Filter: domain.BaseFilter{
			Page:     req.GetPageNumber(),
			PageSize: req.GetPageSize(),
			Query:    req.GetQuery(),
		},
		UserId: req.GetUserId(),
	}

	if req.GetPageNumber() == 0 {
		dto.Filter.Page = 1
	}

	if req.GetPageSize() == 0 {
		dto.Filter.PageSize = 10
	}

	return dto
}

type UnbanParticipant struct {
	EventId       string `json:"event_id"`
	UserId        int64  `json:"user_id"`
	ParticipantId int64  `json:"participant_id"`
}

func ProtoToUnbanParticipant(req *eventv1.UnbanParticipantRequest) *UnbanParticipant {
	return &UnbanParticipant{
		EventId:       req.GetEventId(),
		UserId:        req.GetUserId(),
		ParticipantId: req.GetParticipantId(),
	}
}
