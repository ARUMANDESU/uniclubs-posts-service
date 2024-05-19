package dao

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
)

type User struct {
	ID        int64  `bson:"_id"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Barcode   string `bson:"barcode"`
	AvatarURL string `bson:"avatar_url,omitempty"`
}

type Organizer struct {
	User    `bson:",inline"`
	ClubId  int64 `bson:"club_id"`
	ByWhoId int64 `bson:"by_who_id,omitempty"`
}

// Into dao

func UserFromDomainUser(user domain.User) User {
	return User{
		ID:        user.ID,
		Barcode:   user.Barcode,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
	}
}

func OrganizerFromDomainUser(user domain.User, clubId int64) Organizer {
	return Organizer{
		User:   UserFromDomainUser(user),
		ClubId: clubId,
	}
}

func ToOrganizer(organizer domain.Organizer) Organizer {
	return Organizer{
		User:    UserFromDomainUser(organizer.User),
		ClubId:  organizer.ClubId,
		ByWhoId: organizer.ByWhoId,
	}
}

func ToOrganizers(organizers []domain.Organizer) []Organizer {
	organizerIds := make([]Organizer, len(organizers))
	for i, organizer := range organizers {
		organizerIds[i] = ToOrganizer(organizer)
	}
	return organizerIds
}

// From dao to domain

func ToDomainUser(user User) domain.User {
	return domain.User{
		ID:        user.ID,
		Barcode:   user.Barcode,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
	}
}

func ToDomainUsers(usersMongo []User) []domain.User {
	users := make([]domain.User, len(usersMongo))
	for i, user := range usersMongo {
		users[i] = ToDomainUser(user)
	}
	return users
}

func ToDomainOrganizer(organizer Organizer) domain.Organizer {
	return domain.Organizer{
		User:    ToDomainUser(organizer.User),
		ClubId:  organizer.ClubId,
		ByWhoId: organizer.ByWhoId,
	}
}

func ToDomainOrganizers(organizersMongo []Organizer) []domain.Organizer {
	organizers := make([]domain.Organizer, len(organizersMongo))
	for i, organizer := range organizersMongo {
		organizers[i] = ToDomainOrganizer(organizer)
	}
	return organizers
}
