package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_ToOrganizer(t *testing.T) {
	user := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Barcode:   "123456",
		AvatarURL: "http://example.com/avatar.jpg",
	}

	organizer := user.ToOrganizer(2, 3)

	assert.Equal(t, user.ID, organizer.ID)
	assert.Equal(t, user.FirstName, organizer.FirstName)
	assert.Equal(t, user.LastName, organizer.LastName)
	assert.Equal(t, user.Barcode, organizer.Barcode)
	assert.Equal(t, user.AvatarURL, organizer.AvatarURL)
	assert.Equal(t, int64(2), organizer.ClubId)
	assert.Equal(t, int64(3), organizer.ByWhoId)
}

func TestUser_ToProto(t *testing.T) {
	user := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Barcode:   "123456",
		AvatarURL: "http://example.com/avatar.jpg",
	}

	protoUser := user.ToProto()

	assert.Equal(t, user.ID, protoUser.GetId())
	assert.Equal(t, user.FirstName, protoUser.GetFirstName())
	assert.Equal(t, user.LastName, protoUser.GetLastName())
	assert.Equal(t, user.Barcode, protoUser.GetBarcode())
	assert.Equal(t, user.AvatarURL, protoUser.GetAvatarUrl())
}

func TestUserFromProto(t *testing.T) {
	protoUser := &eventv1.UserObject{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
		Barcode:   "123456",
		AvatarUrl: "http://example.com/avatar.jpg",
	}

	user := UserFromProto(protoUser)

	assert.Equal(t, protoUser.GetId(), user.ID)
	assert.Equal(t, protoUser.GetFirstName(), user.FirstName)
	assert.Equal(t, protoUser.GetLastName(), user.LastName)
	assert.Equal(t, protoUser.GetBarcode(), user.Barcode)
	assert.Equal(t, protoUser.GetAvatarUrl(), user.AvatarURL)
}

func TestOrganizer_IsByWho(t *testing.T) {
	organizer := Organizer{
		User: User{
			ID: 1,
		},
		ByWhoId: 2,
	}

	assert.True(t, organizer.IsByWho(2))
	assert.False(t, organizer.IsByWho(3))
}

func TestOrganizers_ToProto(t *testing.T) {
	organizers := []Organizer{
		{
			User: User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Barcode:   "123456",
				AvatarURL: "http://example.com/avatar.jpg",
			},
			ClubId:  2,
			ByWhoId: 3,
		},
		{
			User: User{
				ID:        4,
				FirstName: "Jane",
				LastName:  "Doe",
				Barcode:   "654321",
				AvatarURL: "http://example.com/avatar2.jpg",
			},
			ClubId:  5,
			ByWhoId: 6,
		},
	}

	protoOrganizers := OrganizersToProto(organizers)

	for i, protoOrganizer := range protoOrganizers {
		assert.Equal(t, organizers[i].ID, protoOrganizer.GetId())
		assert.Equal(t, organizers[i].FirstName, protoOrganizer.GetFirstName())
		assert.Equal(t, organizers[i].LastName, protoOrganizer.GetLastName())
		assert.Equal(t, organizers[i].AvatarURL, protoOrganizer.GetAvatarUrl())
		assert.Equal(t, organizers[i].ClubId, protoOrganizer.GetClubId())
		assert.Equal(t, organizers[i].ByWhoId, protoOrganizer.GetByWhoId())
	}
}

func TestOrganizer_ToProto(t *testing.T) {
	organizer := Organizer{
		User: User{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Barcode:   "123456",
			AvatarURL: "http://example.com/avatar.jpg",
		},
		ClubId:  2,
		ByWhoId: 3,
	}

	protoOrganizer := organizer.ToProto()
	assert.Equal(t, organizer.ID, protoOrganizer.GetId())
	assert.Equal(t, organizer.FirstName, protoOrganizer.GetFirstName())
	assert.Equal(t, organizer.LastName, protoOrganizer.GetLastName())
	assert.Equal(t, organizer.AvatarURL, protoOrganizer.GetAvatarUrl())
	assert.Equal(t, organizer.ClubId, protoOrganizer.GetClubId())
}
