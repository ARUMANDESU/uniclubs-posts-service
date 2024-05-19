package domain

import eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"

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
