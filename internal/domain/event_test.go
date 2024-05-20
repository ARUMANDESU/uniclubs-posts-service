package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventIsOwner(t *testing.T) {
	event := Event{
		OwnerId: 1,
	}

	assert.True(t, event.IsOwner(1))
	assert.False(t, event.IsOwner(2))
}

func TestEventIsOrganizer(t *testing.T) {
	event := Event{
		Organizers: []Organizer{
			{User: User{ID: 1}},
			{User: User{ID: 2}},
		},
	}

	assert.True(t, event.IsOrganizer(1))
	assert.False(t, event.IsOrganizer(3))
}

func TestEventGetOrganizerById(t *testing.T) {
	event := Event{
		Organizers: []Organizer{
			{User: User{ID: 1}},
			{User: User{ID: 2}},
		},
	}

	assert.NotNil(t, event.GetOrganizerById(1))
	assert.Nil(t, event.GetOrganizerById(3))
}

func TestEventAddOrganizer(t *testing.T) {
	event := Event{}
	organizer := Organizer{User: User{ID: 1}}
	organizer2 := Organizer{User: User{ID: 2}}

	event.AddOrganizer(organizer)

	assert.Equal(t, 1, len(event.Organizers))
	assert.Equal(t, organizer, event.Organizers[0])

	event.AddOrganizer(organizer2)

	assert.Equal(t, 2, len(event.Organizers))
	assert.Equal(t, organizer2, event.Organizers[1])
}

func TestEventRemoveOrganizer(t *testing.T) {
	t.Run("RemoveOrganizer", func(t *testing.T) {
		event := Event{
			OwnerId: 1,
			Organizers: []Organizer{
				{User: User{ID: 1}},
				{User: User{ID: 2}},
				{User: User{ID: 3}},
			},
		}

		err := event.RemoveOrganizer(2)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(event.Organizers))

		err = event.RemoveOrganizer(1)
		assert.ErrorIs(t, err, ErrUserIsEventOwner)
		assert.Equal(t, 2, len(event.Organizers))

		err = event.RemoveOrganizer(4)
		assert.NotNil(t, err)

		err = event.RemoveOrganizer(3)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(event.Organizers))
	})

	t.Run("RemoveOrganizer Empty", func(t *testing.T) {
		event := Event{}

		err := event.RemoveOrganizer(1)
		assert.ErrorIs(t, err, ErrOrganizersEmpty)
	})
}

func TestEventRemoveOrganizersByClubId(t *testing.T) {
	t.Run("RemoveOrganizersByClubId", func(t *testing.T) {
		event := Event{
			OwnerId: 1,
			ClubId:  2,
			Organizers: []Organizer{
				{User: User{ID: 1}, ClubId: 1},
				{User: User{ID: 3}, ClubId: 1},
				{User: User{ID: 2}, ClubId: 2},
			},
		}

		err := event.RemoveOrganizersByClubId(1)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(event.Organizers))
		assert.Equal(t, int64(1), event.Organizers[0].ID)
		assert.Equal(t, int64(2), event.Organizers[1].ID)
	})

	t.Run("RemoveOrganizersByClubId Empty", func(t *testing.T) {
		event := Event{}

		err := event.RemoveOrganizersByClubId(1)
		assert.ErrorIs(t, err, ErrOrganizersEmpty)

	})

}

func TestEventIsCollaborator(t *testing.T) {
	event := Event{
		CollaboratorClubs: []Club{
			{ID: 1},
			{ID: 2},
		},
	}

	assert.True(t, event.IsCollaborator(1))
	assert.False(t, event.IsCollaborator(3))
}

func TestEventAddCollaborator(t *testing.T) {
	event := Event{}
	club := Club{ID: 1}

	event.AddCollaborator(club)

	assert.Equal(t, 1, len(event.CollaboratorClubs))
	assert.Equal(t, club, event.CollaboratorClubs[0])
}

func TestEventRemoveCollaborator(t *testing.T) {
	t.Run("RemoveCollaborator", func(t *testing.T) {
		event := Event{
			ClubId: 1,
			CollaboratorClubs: []Club{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
		}

		err := event.RemoveCollaborator(2)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(event.CollaboratorClubs))

		err = event.RemoveCollaborator(1)
		assert.ErrorIs(t, err, ErrClubIsEventOwner)
		assert.Equal(t, 2, len(event.CollaboratorClubs))

		err = event.RemoveCollaborator(4)
		assert.ErrorIs(t, err, ErrCollaboratorNotFound)

		err = event.RemoveCollaborator(3)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(event.CollaboratorClubs))
	})

	t.Run("RemoveCollaboratorEmpty", func(t *testing.T) {
		event := Event{}

		err := event.RemoveCollaborator(1)
		assert.ErrorIs(t, err, ErrCollaboratorsEmpty)

	})
}

func TestEventGetCollaboratorById(t *testing.T) {
	event := Event{
		CollaboratorClubs: []Club{
			{ID: 1},
			{ID: 2},
		},
	}

	assert.NotNil(t, event.GetCollaboratorById(1))
	assert.Equal(t, int64(1), event.GetCollaboratorById(1).ID)
	assert.Nil(t, event.GetCollaboratorById(3))
}

func TestEventToProto(t *testing.T) {
	event := Event{
		ID:                 "1",
		ClubId:             1,
		OwnerId:            1,
		CollaboratorClubs:  []Club{{ID: 1}},
		Organizers:         []Organizer{{User: User{ID: 1}}},
		Title:              "Test Event",
		Description:        "Test Description",
		Type:               "Test Type",
		Status:             EventStatusDraft,
		Tags:               []string{"Test Tag"},
		MaxParticipants:    1,
		ParticipantsCount:  1,
		LocationLink:       "http://example.com/location",
		LocationUniversity: "Test University",
		StartDate:          time.Now(),
		EndDate:            time.Now(),
		CoverImages:        []CoverImage{{File: File{Name: "Test Cover Image"}}},
		AttachedImages:     []File{{Name: "Test Attached Image"}},
		AttachedFiles:      []File{{Name: "Test Attached File"}},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		DeletedAt:          time.Now(),
	}

	protoEvent := event.ToProto()

	assert.Equal(t, event.ID, protoEvent.GetId())
	assert.Equal(t, event.ClubId, protoEvent.GetClubId())
	assert.Equal(t, event.OwnerId, protoEvent.GetOwnerId())
	assert.Equal(t, len(event.CollaboratorClubs), len(protoEvent.GetCollaboratorClubs()))
	assert.Equal(t, len(event.Organizers), len(protoEvent.GetOrganizers()))
	assert.Equal(t, event.Title, protoEvent.GetTitle())
	assert.Equal(t, event.Description, protoEvent.GetDescription())
	assert.Equal(t, event.Type, EventType(protoEvent.GetType()))
	assert.Equal(t, event.Status, EventStatus(protoEvent.GetStatus()))
	assert.Equal(t, len(event.Tags), len(protoEvent.GetTags()))
	assert.Equal(t, event.MaxParticipants, protoEvent.GetMaxParticipants())
	assert.Equal(t, event.ParticipantsCount, protoEvent.GetParticipantsCount())
	assert.Equal(t, event.LocationLink, protoEvent.GetLocationLink())
	assert.Equal(t, event.LocationUniversity, protoEvent.GetLocationUniversity())
	assert.Equal(t, event.StartDate.Format(timeLayout), protoEvent.GetStartDate())
	assert.Equal(t, event.EndDate.Format(timeLayout), protoEvent.GetEndDate())
	assert.Equal(t, len(event.CoverImages), len(protoEvent.GetCoverImages()))
	assert.Equal(t, len(event.AttachedImages), len(protoEvent.GetAttachedImages()))
	assert.Equal(t, len(event.AttachedFiles), len(protoEvent.GetAttachedFiles()))
	assert.Equal(t, event.CreatedAt.Format(timeLayout), protoEvent.GetCreatedAt())
	assert.Equal(t, event.UpdatedAt.Format(timeLayout), protoEvent.GetUpdatedAt())
	assert.Equal(t, event.DeletedAt.Format(timeLayout), protoEvent.GetDeletedAt())
}
