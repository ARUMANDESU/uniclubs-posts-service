package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"time"
)

type UserStatus int32
type ParticipantStatus int32

const (
	UserStatusUnknown   UserStatus = 0
	UserStatusOrganizer UserStatus = 1
	UserStatusOwner     UserStatus = 2

	ParticipantStatusUnknown ParticipantStatus = 0
	ParticipantStatusJoined  ParticipantStatus = 2
	ParticipantStatusBanned  ParticipantStatus = 4
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Barcode   string `json:"barcode"`
	AvatarURL string `json:"avatar_url"`
}

type Organizer struct {
	User    `json:",inline"`
	ClubId  int64 `json:"club_id"`
	ByWhoId int64 `json:"by_who_id"`
}

type Participant struct {
	ID       string `json:"id"`
	EventId  string `json:"event_id"`
	User     `json:",inline"`
	JoinedAt time.Time `json:"joined_at"`
}

func (u User) ToOrganizer(clubId, byWhoId int64) Organizer {
	return Organizer{
		User:    u,
		ClubId:  clubId,
		ByWhoId: byWhoId,
	}
}

func (u User) ToProto() *eventv1.UserObject {
	return &eventv1.UserObject{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Barcode:   u.Barcode,
		AvatarUrl: u.AvatarURL,
	}
}

func UserFromProto(user *eventv1.UserObject) User {
	return User{
		ID:        user.GetId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Barcode:   user.GetBarcode(),
		AvatarURL: user.GetAvatarUrl(),
	}
}

func (o Organizer) ToProto() *eventv1.OrganizerObject {
	return &eventv1.OrganizerObject{
		Id:        o.ID,
		FirstName: o.FirstName,
		LastName:  o.LastName,
		AvatarUrl: o.AvatarURL,
		ClubId:    o.ClubId,
		ByWhoId:   o.ByWhoId,
	}
}

func (o Organizer) IsByWho(userId int64) bool {
	return o.ByWhoId == userId
}

func OrganizersToProto(organizers []Organizer) []*eventv1.OrganizerObject {
	convertedOrganizers := make([]*eventv1.OrganizerObject, len(organizers))
	for i, organizer := range organizers {
		convertedOrganizers[i] = organizer.ToProto()
	}
	return convertedOrganizers
}

func ParticipantsToProto(participants []Participant) []*eventv1.UserObject { // todo: change later to ParticipantObject
	convertedParticipants := make([]*eventv1.UserObject, len(participants))
	for i, participant := range participants {
		convertedParticipants[i] = participant.ToProto()
	}
	return convertedParticipants
}
