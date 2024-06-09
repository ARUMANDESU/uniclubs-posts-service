package domain

import (
	posts "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts"
)

type Club struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logo_url"`
}

func (c Club) ToProto() *posts.ClubObject {
	return &posts.ClubObject{
		Id:      c.ID,
		Name:    c.Name,
		LogoUrl: c.LogoURL,
	}
}

func ClubsToProto(clubs []Club) []*posts.ClubObject {
	convertedClubs := make([]*posts.ClubObject, len(clubs))
	for i, club := range clubs {
		convertedClubs[i] = club.ToProto()
	}
	return convertedClubs
}

func ClubFromProto(club *posts.ClubObject) Club {
	return Club{
		ID:      club.GetId(),
		Name:    club.GetName(),
		LogoURL: club.GetLogoUrl(),
	}
}
