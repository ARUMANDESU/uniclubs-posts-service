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

func ParticipantToDomain(participant Participant) *domain.Participant {
	return &domain.Participant{
		ID:       participant.Id.Hex(),
		EventId:  participant.EventId.Hex(),
		User:     ToDomainUser(participant.User),
		JoinedAt: participant.JoinedAt,
	}
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
