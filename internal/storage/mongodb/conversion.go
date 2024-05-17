package mongodb

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func DomainToModel(event *domain.Event) Event {
	objectID, _ := primitive.ObjectIDFromHex(event.ID)

	return Event{
		ID:                  objectID,
		ClubId:              event.Club.ID,
		UserId:              event.User.ID,
		CollaboratorClubIds: ToCollaboratorClubIds(event.CollaboratorClubs),
		OrganizerIds:        ToOrganizerIds(event.Organizers),
		Title:               event.Title,
		Description:         event.Description,
		Type:                event.Type,
		Status:              event.Status,
		Tags:                event.Tags,
		MaxParticipants:     event.MaxParticipants,
		ParticipantsCount:   event.ParticipantsCount,
		LocationLink:        event.LocationLink,
		LocationUniversity:  event.LocationUniversity,
		StartDate:           event.StartDate,
		EndDate:             event.EndDate,
		CoverImages:         ToCoverImages(event.CoverImages),
		AttachedImages:      ToFiles(event.AttachedImages),
		AttachedFiles:       ToFiles(event.AttachedFiles),
		CreatedAt:           event.CreatedAt,
		UpdatedAt:           event.UpdatedAt,
		DeletedAt:           event.DeletedAt,
	}
}

func ToCoverImages(coverImages []domain.CoverImage) []CoverImageMongo {
	coverImagesMongo := make([]CoverImageMongo, len(coverImages))
	for i, coverImage := range coverImages {
		coverImagesMongo[i] = CoverImageMongo{
			FileMongo: FileMongo{
				URL:  coverImage.Url,
				Name: coverImage.Name,
				Type: coverImage.Type,
			},
			Position: coverImage.Position,
		}
	}
	return coverImagesMongo
}

func ToFiles(files []domain.File) []FileMongo {
	filesMongo := make([]FileMongo, len(files))
	for i, file := range files {
		filesMongo[i] = FileMongo{
			URL:  file.Url,
			Name: file.Name,
			Type: file.Type,
		}
	}
	return filesMongo
}

func ToCollaboratorClubIds(clubs []domain.Club) []int64 {
	clubIds := make([]int64, len(clubs))
	for i, club := range clubs {
		clubIds[i] = club.ID
	}
	return clubIds
}

func ToOrganizerIds(organizers []domain.Organizer) []int64 {
	organizerIds := make([]int64, len(organizers))
	for i, organizer := range organizers {
		organizerIds[i] = organizer.User.ID
	}
	return organizerIds
}

func ToDomainInvite(invite Invite) domain.Invite {
	return domain.Invite{
		ID:      invite.ID.Hex(),
		EventId: invite.EventId.Hex(),
		ClubId:  invite.ClubId,
	}
}

func ToDomainUserInvite(userInvite UserInvite) *domain.UserInvite {
	return &domain.UserInvite{
		Invite: ToDomainInvite(userInvite.Invite),
		UserId: userInvite.UserId,
	}
}

func ToDomainInvites(invites []Invite) []domain.Invite {
	domainInvites := make([]domain.Invite, len(invites))
	for i, invite := range invites {
		domainInvites[i] = ToDomainInvite(invite)
	}
	return domainInvites
}

func ToDomainUserInvites(userInvites []UserInvite) []domain.UserInvite {
	domainUserInvites := make([]domain.UserInvite, len(userInvites))
	for i, userInvite := range userInvites {
		domainUserInvites[i] = *ToDomainUserInvite(userInvite)
	}
	return domainUserInvites
}
