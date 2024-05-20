package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClubToProto(t *testing.T) {
	club := Club{
		ID:      1,
		Name:    "Test Club",
		LogoURL: "http://example.com/logo.jpg",
	}

	protoClub := club.ToProto()

	assert.Equal(t, club.ID, protoClub.GetId())
	assert.Equal(t, club.Name, protoClub.GetName())
	assert.Equal(t, club.LogoURL, protoClub.GetLogoUrl())
}

func TestClubsToProto(t *testing.T) {
	clubs := []Club{
		{
			ID:      1,
			Name:    "Test Club 1",
			LogoURL: "http://example.com/logo1.jpg",
		},
		{
			ID:      2,
			Name:    "Test Club 2",
			LogoURL: "http://example.com/logo2.jpg",
		},
	}

	protoClubs := ClubsToProto(clubs)

	for i, protoClub := range protoClubs {
		assert.Equal(t, clubs[i].ID, protoClub.GetId())
		assert.Equal(t, clubs[i].Name, protoClub.GetName())
		assert.Equal(t, clubs[i].LogoURL, protoClub.GetLogoUrl())
	}
}

func TestClubFromProto(t *testing.T) {
	protoClub := &eventv1.ClubObject{
		Id:      1,
		Name:    "Test Club",
		LogoUrl: "http://example.com/logo.jpg",
	}

	club := ClubFromProto(protoClub)

	assert.Equal(t, protoClub.GetId(), club.ID)
	assert.Equal(t, protoClub.GetName(), club.Name)
	assert.Equal(t, protoClub.GetLogoUrl(), club.LogoURL)
}
