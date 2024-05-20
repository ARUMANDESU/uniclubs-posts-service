package validate

import (
	"testing"

	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
)

func TestAddCollaborator(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		req := &eventv1.AddCollaboratorRequest{
			EventId: "1",
			UserId:  1,
			Club:    &eventv1.ClubObject{Id: 1, Name: "Test Club", LogoUrl: "http://example.com/logo.jpg"},
		}

		err := AddCollaborator(req)
		assert.Nil(t, err)
	})

	t.Run("InvalidRequestMissingEventId", func(t *testing.T) {
		req := &eventv1.AddCollaboratorRequest{
			UserId: 1,
			Club:   &eventv1.ClubObject{Id: 1, Name: "Test Club", LogoUrl: "http://example.com/logo.jpg"},
		}

		err := AddCollaborator(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingUserId", func(t *testing.T) {
		req := &eventv1.AddCollaboratorRequest{
			EventId: "1",
			Club:    &eventv1.ClubObject{Id: 1, Name: "Test Club", LogoUrl: "http://example.com/logo.jpg"},
		}

		err := AddCollaborator(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingClub", func(t *testing.T) {
		req := &eventv1.AddCollaboratorRequest{
			EventId: "1",
			UserId:  1,
		}

		err := AddCollaborator(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		req := "invalid"

		err := AddCollaborator(req)
		assert.NotNil(t, err)
	})
}

func TestAddOrganizer(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		req := &eventv1.AddOrganizerRequest{
			EventId:      "1",
			UserId:       1,
			Target:       &eventv1.UserObject{Id: 2, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
			TargetClubId: 1,
		}

		err := AddOrganizer(req)
		assert.Nil(t, err)
	})

	t.Run("InvalidRequestSameUserIdAndTargetId", func(t *testing.T) {
		req := &eventv1.AddOrganizerRequest{
			EventId:      "1",
			UserId:       1,
			Target:       &eventv1.UserObject{Id: 1, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
			TargetClubId: 1,
		}

		err := AddOrganizer(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingUserId", func(t *testing.T) {
		req := &eventv1.AddOrganizerRequest{
			EventId:      "1",
			Target:       &eventv1.UserObject{Id: 2, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
			TargetClubId: 1,
		}

		err := AddOrganizer(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingTarget", func(t *testing.T) {
		req := &eventv1.AddOrganizerRequest{
			EventId:      "1",
			UserId:       1,
			TargetClubId: 1,
		}

		err := AddOrganizer(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingTargetClubId", func(t *testing.T) {
		req := &eventv1.AddOrganizerRequest{
			EventId: "1",
			UserId:  1,
			Target:  &eventv1.UserObject{Id: 2, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
		}

		err := AddOrganizer(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		req := "invalid"

		err := AddOrganizer(req)
		assert.NotNil(t, err)
	})
}

func TestHandleInviteUser(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		req := &eventv1.HandleInviteUserRequest{
			InviteId: "1",
			UserId:   1,
		}

		err := HandleInviteUser(req)
		assert.Nil(t, err)
	})

	t.Run("InvalidRequestMissingInviteId", func(t *testing.T) {
		req := &eventv1.HandleInviteUserRequest{
			UserId: 1,
		}

		err := HandleInviteUser(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingUserId", func(t *testing.T) {
		req := &eventv1.HandleInviteUserRequest{
			InviteId: "1",
		}

		err := HandleInviteUser(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		req := "invalid"

		err := HandleInviteUser(req)
		assert.NotNil(t, err)
	})
}

func TestHandleInviteClub(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		req := &eventv1.HandleInviteClubRequest{
			InviteId: "1",
			ClubId:   1,
			User:     &eventv1.UserObject{Id: 1, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
		}

		err := HandleInviteClub(req)
		assert.Nil(t, err)
	})

	t.Run("InvalidRequestMissingInviteId", func(t *testing.T) {
		req := &eventv1.HandleInviteClubRequest{
			ClubId: 1,
			User:   &eventv1.UserObject{Id: 1, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
		}

		err := HandleInviteClub(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingClubId", func(t *testing.T) {
		req := &eventv1.HandleInviteClubRequest{
			InviteId: "1",
			User:     &eventv1.UserObject{Id: 1, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
		}

		err := HandleInviteClub(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingUser", func(t *testing.T) {
		req := &eventv1.HandleInviteClubRequest{
			InviteId: "1",
			ClubId:   1,
		}

		err := HandleInviteClub(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		req := "invalid"

		err := HandleInviteClub(req)
		assert.NotNil(t, err)
	})
}

func TestRevokeInvite(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		req := &eventv1.RevokeInviteRequest{
			InviteId: "1",
			UserId:   1,
		}

		err := RevokeInvite(req)
		assert.Nil(t, err)
	})

	t.Run("InvalidRequestMissingInviteId", func(t *testing.T) {
		req := &eventv1.RevokeInviteRequest{
			UserId: 1,
		}

		err := RevokeInvite(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingUserId", func(t *testing.T) {
		req := &eventv1.RevokeInviteRequest{
			InviteId: "1",
		}

		err := RevokeInvite(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		req := "invalid"

		err := RevokeInvite(req)
		assert.NotNil(t, err)
	})
}
