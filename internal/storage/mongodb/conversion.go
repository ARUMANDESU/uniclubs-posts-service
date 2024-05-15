package mongodb

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
)

func ToDomainEvent(event Event, user domain.User, club domain.Club, organizers []domain.Organizer, collaboratorClubs []domain.Club) *domain.Event {
	return &domain.Event{
		ID:                 event.ID.Hex(),
		Club:               club,
		User:               user,
		CollaboratorClubs:  collaboratorClubs,
		Organizers:         organizers,
		Title:              event.Title,
		Description:        event.Description,
		Type:               event.Type,
		Status:             event.Status,
		Tags:               event.Tags,
		MaxParticipants:    event.MaxParticipants,
		ParticipantsCount:  event.ParticipantsCount,
		LocationLink:       event.LocationLink,
		LocationUniversity: event.LocationUniversity,
		StartDate:          event.StartDate,
		EndDate:            event.EndDate,
		CoverImages:        ToDomainCoverImages(event.CoverImages),
		AttachedImages:     ToDomainFiles(event.AttachedImages),
		AttachedFiles:      ToDomainFiles(event.AttachedFiles),
		CreatedAt:          event.CreatedAt,
		UpdatedAt:          event.UpdatedAt,
		DeletedAt:          event.DeletedAt,
	}
}

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

func ToDomainOrganizers(organizersMongo []Organizer) []domain.Organizer {
	organizers := make([]domain.Organizer, len(organizersMongo))
	for i, organizer := range organizersMongo {
		organizers[i] = domain.Organizer{
			User:   ToDomainUser(organizer.User),
			ClubId: organizer.ClubId,
		}
	}
	return organizers
}

func ToDomainFiles(filesMongo []FileMongo) []domain.File {
	files := make([]domain.File, len(filesMongo))
	for i, file := range filesMongo {
		files[i] = domain.File{
			Url:  file.URL,
			Name: file.Name,
			Type: file.Type,
		}
	}
	return files
}

func ToDomainCoverImages(coverImagesMongo []CoverImageMongo) []domain.CoverImage {
	coverImages := make([]domain.CoverImage, len(coverImagesMongo))
	for i, coverImage := range coverImagesMongo {
		coverImages[i] = domain.CoverImage{
			File: domain.File{
				Url:  coverImage.URL,
				Name: coverImage.Name,
				Type: coverImage.Type,
			},
			Position: coverImage.Position,
		}
	}
	return coverImages
}

func ToDomainClub(club Club) domain.Club {
	return domain.Club{
		ID:      club.ID,
		Name:    club.Name,
		LogoURL: club.LogoURL,
	}
}

func ToDomainClubs(clubsMongo []Club) []domain.Club {
	clubs := make([]domain.Club, len(clubsMongo))
	for i, club := range clubsMongo {
		clubs[i] = ToDomainClub(club)
	}
	return clubs
}
