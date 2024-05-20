package validate

import (
	"testing"

	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
)

func TestRemoveCollaborator(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		req := &eventv1.RemoveCollaboratorRequest{
			EventId: "1",
			UserId:  1,
			ClubId:  1,
		}

		err := RemoveCollaborator(req)
		assert.Nil(t, err)
	})

	t.Run("InvalidRequestMissingEventId", func(t *testing.T) {
		req := &eventv1.RemoveCollaboratorRequest{
			UserId: 1,
			ClubId: 1,
		}

		err := RemoveCollaborator(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingUserId", func(t *testing.T) {
		req := &eventv1.RemoveCollaboratorRequest{
			EventId: "1",
			ClubId:  1,
		}

		err := RemoveCollaborator(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidRequestMissingClubId", func(t *testing.T) {
		req := &eventv1.RemoveCollaboratorRequest{
			EventId: "1",
			UserId:  1,
		}

		err := RemoveCollaborator(req)
		assert.NotNil(t, err)
	})

	t.Run("InvalidType", func(t *testing.T) {
		req := "invalid"

		err := RemoveCollaborator(req)
		assert.NotNil(t, err)
	})
}
