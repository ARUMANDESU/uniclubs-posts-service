package dao

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Participant struct {
	Id       primitive.ObjectID `bson:"_id"`
	EventId  primitive.ObjectID `bson:"event_id"`
	User     User               `bson:"user"`
	JoinedAt time.Time          `bson:"joined_at,omitempty"`
}

type BanRecord struct {
	EventId  primitive.ObjectID `bson:"event_id"`
	UserId   int64              `bson:"user_id"`
	BannedAt time.Time          `bson:"banned_at,omitempty"`
	Reason   string             `bson:"reason,omitempty"`
	ByWhoId  int64              `bson:"by_who_id,omitempty"`
}

func ParticipantToDomain(participant Participant) *domain.Participant {
	return &domain.Participant{
		ID:       participant.Id.Hex(),
		EventId:  participant.EventId.Hex(),
		User:     ToDomainUser(participant.User),
		JoinedAt: participant.JoinedAt,
	}
}

func ParticipantsToDomain(participants []Participant) []domain.Participant {
	var result []domain.Participant
	for _, participant := range participants {
		result = append(result, *ParticipantToDomain(participant))
	}
	return result

}

func ParticipantFromDomain(participant *domain.Participant) (*Participant, error) {
	eventId, err := primitive.ObjectIDFromHex(participant.EventId)
	if err != nil {
		return nil, err
	}
	participantRecordId, err := primitive.ObjectIDFromHex(participant.ID)
	if err != nil {
		return nil, err
	}

	return &Participant{
		Id:       participantRecordId,
		EventId:  eventId,
		User:     UserFromDomainUser(participant.User),
		JoinedAt: participant.JoinedAt,
	}, nil
}

func BanRecordToDomain(banRecord BanRecord) *domain.BanRecord {
	return &domain.BanRecord{
		EventId:  banRecord.EventId.Hex(),
		UserId:   banRecord.UserId,
		BannedAt: banRecord.BannedAt,
		Reason:   banRecord.Reason,
		ByWhoId:  banRecord.ByWhoId,
	}
}
