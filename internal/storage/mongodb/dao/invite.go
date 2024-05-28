package dao

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClubInvite struct {
	ID      primitive.ObjectID `bson:"_id"`
	EventId primitive.ObjectID `bson:"event_id"`
	Club    Club               `bson:"club"`
}

type OrganizerInvite struct {
	ID      primitive.ObjectID `bson:"_id"`
	EventId primitive.ObjectID `bson:"event_id"`
	ClubId  int64              `bson:"club_id"`
	ByWhoId int64              `bson:"by_who_id"`
	User    User               `bson:"user"`
}

func ToDomainInvite(i ClubInvite) *domain.Invite {
	return &domain.Invite{
		ID: i.ID.Hex(),
		Event: domain.Event{
			ID: i.EventId.Hex(),
		},
		Club: ToDomainClub(i.Club),
	}
}

func ToDomainUserInvite(u OrganizerInvite) *domain.UserInvite {
	return &domain.UserInvite{
		ID: u.ID.Hex(),
		Event: domain.Event{
			ID: u.EventId.Hex(),
		},
		ClubId:  u.ClubId,
		ByWhoId: u.ByWhoId,
		User:    ToDomainUser(u.User),
	}
}

func ToDomainInvites(invites []ClubInvite) []domain.Invite {
	domainInvites := make([]domain.Invite, len(invites))
	for i, invite := range invites {
		domainInvites[i] = *ToDomainInvite(invite)
	}
	return domainInvites
}

func ToDomainUserInvites(userInvites []OrganizerInvite) []domain.UserInvite {
	domainUserInvites := make([]domain.UserInvite, len(userInvites))
	for i, userInvite := range userInvites {
		domainUserInvites[i] = *ToDomainUserInvite(userInvite)
	}
	return domainUserInvites
}

func ToDomainClubInvite(i ClubInvite, event Event) domain.Invite {
	return domain.Invite{
		ID:    i.ID.Hex(),
		Event: *ToDomainEvent(event),
		Club:  ToDomainClub(i.Club),
	}
}

func ToDomainOrgInvite(u OrganizerInvite, event Event) domain.UserInvite {
	return domain.UserInvite{
		ID:      u.ID.Hex(),
		Event:   *ToDomainEvent(event),
		ClubId:  u.ClubId,
		ByWhoId: u.ByWhoId,
		User:    ToDomainUser(u.User),
	}
}
