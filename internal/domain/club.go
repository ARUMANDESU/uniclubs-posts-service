package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
)

type Club struct {
	ID      int64  `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	LogoURL string `json:"logo_url" bson:"logo_url"`
}

func (c Club) ToProto() *eventv1.ClubObject {
	return &eventv1.ClubObject{
		Id:      c.ID,
		Name:    c.Name,
		LogoUrl: c.LogoURL,
	}
}

func ClubsToProto(clubs []Club) []*eventv1.ClubObject {
	convertedClubs := make([]*eventv1.ClubObject, len(clubs))
	for _, club := range clubs {
		convertedClubs = append(convertedClubs, club.ToProto())
	}
	return convertedClubs
}
