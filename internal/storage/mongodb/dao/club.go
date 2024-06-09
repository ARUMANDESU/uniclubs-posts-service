package dao

import "github.com/arumandesu/uniclubs-posts-service/internal/domain"

type Club struct {
	ID      int64  `json:"id,omitempty" bson:"_id"`
	Name    string `json:"name,omitempty" bson:"name"`
	LogoURL string `json:"logo_url,omitempty" bson:"logo_url"`
}

// Into dao

func ClubFromDomain(club domain.Club) Club {
	return Club{
		ID:      club.ID,
		Name:    club.Name,
		LogoURL: club.LogoURL,
	}
}

func ToCollaboratorClubs(clubs []domain.Club) []Club {
	clubIds := make([]Club, len(clubs))
	for i, club := range clubs {
		clubIds[i] = ClubFromDomain(club)
	}
	return clubIds
}

// From dao to domain

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
