package eventmanagement

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event/management/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"testing"
	"time"
)

type suite struct {
	ManagementService Service
	mockStorage       *mocks.EventStorage
}

func newSuite(t *testing.T) *suite {
	t.Helper()

	mockStorage := mocks.NewEventStorage(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	return &suite{
		ManagementService: New(log, mockStorage),
		mockStorage:       mockStorage,
	}
}

func TestService_CreateEvent_HappyPath(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	club := domain.Club{}
	user := domain.User{}

	suite.mockStorage.On(
		"CreateEvent",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("domain.Club"),
		mock.AnythingOfType("domain.User"),
	).Return(&domain.Event{}, nil)

	event, err := suite.ManagementService.CreateEvent(ctx, club, user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.NotNil(t, event)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_CreateEvent_StorageError(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	club := domain.Club{}
	user := domain.User{}
	expectedErr := errors.New("storage error")

	suite.mockStorage.On(
		"CreateEvent",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("domain.Club"),
		mock.AnythingOfType("domain.User"),
	).Return(nil, expectedErr)

	_, err := suite.ManagementService.CreateEvent(ctx, club, user)
	require.ErrorIs(t, err, expectedErr)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_UpdateEvent_HappyPath(t *testing.T) {

	dto := &dtos.UpdateEvent{
		EventId:            "event_id",
		UserId:             1,
		Title:              "updated title",
		Description:        "updated description",
		Type:               domain.EventTypeUniversity,
		Tags:               []string{"updated tag1", "updated tag2"},
		MaxParticipants:    100,
		LocationLink:       "updated location link",
		LocationUniversity: "updated location university",
		StartDate:          "2006-01-02T15:04:05Z",
		EndDate:            "2006-01-02T15:04:05Z",
		CoverImages: []domain.CoverImage{
			{
				File: domain.File{
					Name: "new cover image",
					Url:  "new cover image url",
					Type: "image",
				},
				Position: 1,
			},
		},
		AttachedImages: []domain.File{
			{
				Name: "new attached image",
				Url:  "new attached image url",
				Type: "image",
			},
		},
		AttachedFiles: []domain.File{
			{
				Name: "new attached file",
				Url:  "new attached file url",
				Type: "file",
			},
		},
	}
	updateStartDate, err := time.Parse(time.RFC3339, dto.StartDate)
	require.NoError(t, err)
	updateEndDate, err := time.Parse(time.RFC3339, dto.EndDate)
	require.NoError(t, err)
	oldEvent := &domain.Event{
		ID:                 "event_id",
		ClubId:             1,
		OwnerId:            1,
		CollaboratorClubs:  nil,
		Organizers:         nil,
		Title:              "old title",
		Description:        "old description",
		Type:               "old type",
		Status:             "old status",
		Tags:               []string{"old tag1", "old tag2"},
		MaxParticipants:    50,
		ParticipantsCount:  0,
		LocationLink:       "old location link",
		LocationUniversity: "old location university",
		StartDate:          time.Time{},
		EndDate:            time.Time{},
		CoverImages: []domain.CoverImage{
			{
				File: domain.File{
					Name: "old cover image",
					Url:  "old cover image url",
					Type: "image",
				},
				Position: 1,
			},
		},
		AttachedImages: []domain.File{
			{
				Name: "old attached image",
				Url:  "old attached image url",
				Type: "image",
			},
		},
		AttachedFiles: []domain.File{
			{
				Name: "old attached file",
				Url:  "old attached file url",
				Type: "file",
			},
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: time.Time{},
	}

	tests := []struct {
		name          string
		updatePath    []string
		oldEvent      *domain.Event
		expectedEvent *domain.Event
	}{
		{
			name: "update all",
			updatePath: []string{
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
			},
			oldEvent: &domain.Event{
				ID:                 "event_id",
				OwnerId:            1,
				Title:              "old title",
				Description:        "old description",
				Type:               "old type",
				Tags:               []string{"old tag1", "old tag2"},
				MaxParticipants:    50,
				LocationLink:       "old location link",
				LocationUniversity: "old location university",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				CoverImages: []domain.CoverImage{
					{
						File: domain.File{
							Name: "old cover image",
							Url:  "old cover image url",
							Type: "image",
						},
						Position: 1,
					},
				},
				AttachedImages: []domain.File{
					{
						Name: "old attached image",
						Url:  "old attached image url",
						Type: "image",
					},
				},
				AttachedFiles: []domain.File{
					{
						Name: "old attached file",
						Url:  "old attached file url",
						Type: "file",
					},
				},
			},
			expectedEvent: &domain.Event{
				ID:                 oldEvent.ID,
				OwnerId:            oldEvent.OwnerId,
				Title:              dto.Title,
				Description:        dto.Description,
				Type:               dto.Type,
				Tags:               dto.Tags,
				MaxParticipants:    dto.MaxParticipants,
				LocationLink:       dto.LocationLink,
				LocationUniversity: dto.LocationUniversity,
				StartDate:          updateStartDate,
				EndDate:            updateEndDate,
				CoverImages:        dto.CoverImages,
				AttachedImages:     dto.AttachedImages,
				AttachedFiles:      dto.AttachedFiles,
			},
		},
		{
			name:       "update title",
			updatePath: []string{"title"},
			oldEvent: &domain.Event{
				ID:      "event_id",
				OwnerId: 1,
				Title:   "old title",
			},
			expectedEvent: &domain.Event{
				ID:      oldEvent.ID,
				OwnerId: oldEvent.OwnerId,
				Title:   dto.Title,
			},
		},
		{
			name:       "update description",
			updatePath: []string{"description"},
			oldEvent: &domain.Event{
				ID:          "event_id",
				OwnerId:     1,
				Description: "old description",
			},
			expectedEvent: &domain.Event{
				ID:          oldEvent.ID,
				OwnerId:     oldEvent.OwnerId,
				Description: dto.Description,
			},
		},
		{
			name:       "update type",
			updatePath: []string{"type"},
			oldEvent: &domain.Event{
				ID:      "event_id",
				OwnerId: 1,
				Type:    "old type",
			},
			expectedEvent: &domain.Event{
				ID:      oldEvent.ID,
				OwnerId: oldEvent.OwnerId,
				Type:    dto.Type,
			},
		},
		{
			name:       "update tags",
			updatePath: []string{"tags"},
			oldEvent: &domain.Event{
				ID:      "event_id",
				OwnerId: 1,
				Tags:    []string{"old tag1", "old tag2"},
			},
			expectedEvent: &domain.Event{
				ID:      oldEvent.ID,
				OwnerId: oldEvent.OwnerId,
				Tags:    dto.Tags,
			},
		},
		{
			name:       "update max_participants",
			updatePath: []string{"max_participants"},
			oldEvent: &domain.Event{
				ID:              "event_id",
				OwnerId:         1,
				MaxParticipants: 50,
			},
			expectedEvent: &domain.Event{
				ID:              oldEvent.ID,
				OwnerId:         oldEvent.OwnerId,
				MaxParticipants: dto.MaxParticipants,
			},
		},
		{
			name:       "update location_link",
			updatePath: []string{"location_link"},
			oldEvent: &domain.Event{
				ID:           "event_id",
				OwnerId:      1,
				LocationLink: "old location link",
			},
			expectedEvent: &domain.Event{
				ID:           oldEvent.ID,
				OwnerId:      oldEvent.OwnerId,
				LocationLink: dto.LocationLink,
			},
		},
		{
			name:       "update location_university",
			updatePath: []string{"location_university"},
			oldEvent: &domain.Event{
				ID:                 "event_id",
				OwnerId:            1,
				LocationUniversity: "old location university",
			},
			expectedEvent: &domain.Event{
				ID:                 oldEvent.ID,
				OwnerId:            oldEvent.OwnerId,
				LocationUniversity: dto.LocationUniversity,
			},
		},
		{
			name:       "update start_date",
			updatePath: []string{"start_date"},
			oldEvent: &domain.Event{
				ID:        "event_id",
				OwnerId:   1,
				StartDate: time.Time{}, // Assuming oldEvent.StartDate is an empty string
			},
			expectedEvent: &domain.Event{
				ID:        oldEvent.ID,
				OwnerId:   oldEvent.OwnerId,
				StartDate: updateStartDate, // Assuming dto.StartDate is an empty string
			},
		},
		{
			name:       "update end_date",
			updatePath: []string{"end_date"},
			oldEvent: &domain.Event{
				ID:      "event_id",
				OwnerId: 1,
				EndDate: time.Time{}, // Assuming oldEvent.EndDate is an empty string
			},
			expectedEvent: &domain.Event{
				ID:      oldEvent.ID,
				OwnerId: oldEvent.OwnerId,
				EndDate: updateEndDate, // Assuming dto.EndDate is an empty string
			},
		},
		{
			name:       "update cover_images",
			updatePath: []string{"cover_images"},
			oldEvent: &domain.Event{
				ID:          "event_id",
				OwnerId:     1,
				CoverImages: oldEvent.CoverImages, // Assuming oldCoverImages is an empty slice
			},
			expectedEvent: &domain.Event{
				ID:          oldEvent.ID,
				OwnerId:     oldEvent.OwnerId,
				CoverImages: dto.CoverImages,
			},
		},
		{
			name:       "update attached_images",
			updatePath: []string{"attached_images"},
			oldEvent: &domain.Event{
				ID:             "event_id",
				OwnerId:        1,
				AttachedImages: oldEvent.AttachedImages, // Assuming oldAttachedImages is an empty slice
			},
			expectedEvent: &domain.Event{
				ID:             oldEvent.ID,
				OwnerId:        oldEvent.OwnerId,
				AttachedImages: dto.AttachedImages,
			},
		},
		{
			name:       "update attached_files",
			updatePath: []string{"attached_files"},
			oldEvent: &domain.Event{
				ID:            "event_id",
				OwnerId:       1,
				AttachedFiles: oldEvent.AttachedFiles, // Assuming oldAttachedFiles is an empty slice
			},
			expectedEvent: &domain.Event{
				ID:            oldEvent.ID,
				OwnerId:       oldEvent.OwnerId,
				AttachedFiles: dto.AttachedFiles,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)

			ctx := context.Background()
			dto := dto
			dto.Paths = tt.updatePath

			suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(tt.oldEvent, nil)
			suite.mockStorage.On("UpdateEvent", mock.Anything, tt.expectedEvent).Return(tt.expectedEvent, nil)

			event, err := suite.ManagementService.UpdateEvent(ctx, dto)
			require.NoError(t, err)
			assert.NotNil(t, event)
			assert.ObjectsAreEqual(tt.expectedEvent, event)

			suite.mockStorage.AssertExpectations(t)
		})
	}

}
