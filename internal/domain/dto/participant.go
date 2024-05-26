package dtos

import eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"

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
	EventId       string `json:"event_id"`
	UserId        int64  `json:"user_id"`
	ParticipantId int64  `json:"participant_id"`
	Reason        string `json:"reason"`
}

func ProtoToBanParticipant(req *eventv1.BanParticipantRequest) *BanParticipant {
	return &BanParticipant{
		EventId:       req.GetEventId(),
		UserId:        req.GetUserId(),
		ParticipantId: req.GetParticipantId(),
		Reason:        req.GetReason(),
	}
}
