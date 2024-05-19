package dao

import "github.com/arumandesu/uniclubs-posts-service/internal/domain"

type User struct {
	ID        int64  `json:"id,omitempty" bson:"_id"`
	FirstName string `json:"first_name,omitempty" bson:"first_name"`
	LastName  string `json:"last_name,omitempty" bson:"last_name"`
	Barcode   string `json:"barcode,omitempty" bson:"barcode"`
	AvatarURL string `json:"avatar_url,omitempty" bson:"avatar_url"`
}

type Organizer struct {
	User
	ClubId int64 `json:"club_id" bson:"club_id"`
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

func ToOrganizers(organizers []domain.Organizer) []Organizer {
	organizerIds := make([]Organizer, len(organizers))
	for i, organizer := range organizers {
		organizerIds[i] = Organizer{
			User:   UserFromDomainUser(organizer.User),
			ClubId: organizer.ClubId,
		}
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
		User:   ToDomainUser(organizer.User),
		ClubId: organizer.ClubId,
	}
}

func ToDomainOrganizers(organizersMongo []Organizer) []domain.Organizer {
	organizers := make([]domain.Organizer, len(organizersMongo))
	for i, organizer := range organizersMongo {
		organizers[i] = ToDomainOrganizer(organizer)
	}
	return organizers
}
