package dtos

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"testing"
	"time"
)

func TestUpdateToDTO(t *testing.T) {

	startDate, _ := time.Parse(domain.TimeLayout, "2024-05-12 15:04:05.999999999 -0700 MST")
	endDate, _ := time.Parse(domain.TimeLayout, "2024-06-02 15:04:05.999999999 -0700 MST")
	paths := []string{
		"title",
		"description",
		"type",
		"tags",
		"max_participants",
		"location_link",
		"location_university",
		"start_date",
		"end_date",
		"cover_images",
		"attached_images",
		"attached_files",
	}

	event := &eventv1.UpdateEventRequest{
		EventId:               "1",
		UserId:                1,
		Title:                 "Test Title",
		Description:           "Test Description",
		Type:                  "Test Type",
		Tags:                  []string{"Test Tag"},
		MaxParticipants:       1,
		LocationLink:          "http://example.com/location",
		LocationUniversity:    "Test University",
		StartDate:             "2024-05-12 15:04:05.999999999 -0700 MST",
		EndDate:               "2024-06-02 15:04:05.999999999 -0700 MST",
		CoverImages:           []*eventv1.CoverImage{{Name: "Test Cover Image", Url: "http://example.com/cover.jpg", Type: "jpg", Position: 1}},
		AttachedImages:        []*eventv1.FileObject{{Name: "Test Attached Image", Url: "http://example.com/image.jpg", Type: "jpg"}},
		AttachedFiles:         []*eventv1.FileObject{{Name: "Test Attached File", Url: "http://example.com/file.pdf", Type: "pdf"}},
		IsHiddenForNonMembers: false,
		UpdateMask:            &fieldmaskpb.FieldMask{Paths: paths},
	}

	expectedEvent := &UpdateEvent{
		EventId:            "1",
		UserId:             1,
		Title:              "Test Title",
		Description:        "Test Description",
		Type:               "Test Type",
		Tags:               []string{"Test Tag"},
		MaxParticipants:    1,
		LocationLink:       "http://example.com/location",
		LocationUniversity: "Test University",
		StartDate:          startDate,
		EndDate:            endDate,
		CoverImages: []domain.CoverImage{
			{
				File: domain.File{
					Name: "Test Cover Image",
					Url:  "http://example.com/cover.jpg",
					Type: "jpg",
				},
				Position: 1,
			},
		},
		AttachedImages: []domain.File{
			{
				Name: "Test Attached Image",
				Url:  "http://example.com/image.jpg",
				Type: "jpg",
			},
		},
		AttachedFiles: []domain.File{
			{
				Name: "Test Attached File",
				Url:  "http://example.com/file.pdf",
				Type: "pdf",
			},
		},
		IsHiddenForNonMembers: false,
		Paths:                 paths,
	}

	dto, err := UpdateToDTO(event)
	require.NoError(t, err)
	assert.Equal(t, expectedEvent, dto)
}

func TestUpdateToDTO_FailPath(t *testing.T) {
	tests := []struct {
		name  string
		event *eventv1.UpdateEventRequest
	}{
		{
			name: "InvalidStartDate",
			event: &eventv1.UpdateEventRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"start_date"}},
				StartDate:  "invalid date",
			},
		},
		{
			name: "InvalidEndDate",
			event: &eventv1.UpdateEventRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"end_date"}},
				EndDate:    "invalid date",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UpdateToDTO(tt.event)
			assert.Error(t, err)
		})
	}
}
func TestAddOrganizerRequestToUserToDTO(t *testing.T) {
	event := &eventv1.AddOrganizerRequest{
		EventId:      "1",
		UserId:       1,
		Target:       &eventv1.UserObject{Id: 1, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
		TargetClubId: 1,
	}

	dto := AddOrganizerRequestToUserToDTO(event)

	assert.Equal(t, event.GetEventId(), dto.EventId)
	assert.Equal(t, event.GetUserId(), dto.UserId)
	assert.Equal(t, event.GetTarget().GetId(), dto.Target.ID)
	assert.Equal(t, event.GetTargetClubId(), dto.TargetClubId)
}

func TestAddCollaboratorRequestToClubToDTO(t *testing.T) {
	event := &eventv1.AddCollaboratorRequest{
		EventId: "1",
		UserId:  1,
		Club:    &eventv1.ClubObject{Id: 1, Name: "Test Club", LogoUrl: "http://example.com/logo.jpg"},
	}

	dto := AddCollaboratorRequestToClubToDTO(event)

	assert.Equal(t, event.GetEventId(), dto.EventId)
	assert.Equal(t, event.GetUserId(), dto.UserId)
	assert.Equal(t, event.GetClub().GetId(), dto.Club.ID)
}

func TestAcceptJoinRequestClubToDTO(t *testing.T) {
	event := &eventv1.HandleInviteClubRequest{
		InviteId: "1",
		ClubId:   1,
		User:     &eventv1.UserObject{Id: 1, FirstName: "John", LastName: "Doe", Barcode: "123456", AvatarUrl: "http://example.com/avatar.jpg"},
	}

	dto := AcceptJoinRequestClubToDTO(event)

	assert.Equal(t, event.GetInviteId(), dto.InviteId)
	assert.Equal(t, event.GetClubId(), dto.ClubId)
	assert.Equal(t, event.GetUser().GetId(), dto.User.ID)
}
