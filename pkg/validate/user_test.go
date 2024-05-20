package validate

import (
	"testing"

	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
)

func TestRemoveOrganizerValidation(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		req := &eventv1.RemoveOrganizerRequest{
			EventId:     "1",
			UserId:      1,
			OrganizerId: 2,
		}

		err := RemoveOrganizer(req)
		assert.Nil(t, err)
	})

	t.Run("Invalid request same userId And organizerId", func(t *testing.T) {
		req := &eventv1.RemoveOrganizerRequest{
			EventId:     "1",
			UserId:      1,
			OrganizerId: 1,
		}

		err := RemoveOrganizer(req)
		assert.NotNil(t, err)
	})

	t.Run("Invalid request: missing userId", func(t *testing.T) {
		req := &eventv1.RemoveOrganizerRequest{
			EventId:     "1",
			OrganizerId: 2,
		}

		err := RemoveOrganizer(req)
		assert.NotNil(t, err)
	})

	t.Run("Invalid request: missing organizerId", func(t *testing.T) {
		req := &eventv1.RemoveOrganizerRequest{
			EventId: "1",
			UserId:  1,
		}

		err := RemoveOrganizer(req)
		assert.NotNil(t, err)
	})

	t.Run("Invalid request: missing eventId", func(t *testing.T) {
		req := &eventv1.RemoveOrganizerRequest{
			UserId:      1,
			OrganizerId: 2,
		}

		err := RemoveOrganizer(req)
		assert.NotNil(t, err)
	})
}
