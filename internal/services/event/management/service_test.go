package eventmanagement

import (
	"context"
	"errors"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain/dto"
	eventservice "github.com/arumandesu/uniclubs-posts-service/internal/services/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/services/event/management/mocks"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
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
		StartDate:          time.Now(),
		EndDate:            time.Now().AddDate(0, 0, 20),
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
		Paths: map[string]bool{
			"title":               true,
			"description":         true,
			"type":                true,
			"tags":                true,
			"max_participants":    true,
			"location_link":       true,
			"location_university": true,
			"start_date":          true,
			"end_date":            true,
			"cover_images":        true,
			"attached_images":     true,
			"attached_files":      true,
		},
	}

	oldEvent := &domain.Event{
		ID:                 "event_id",
		OwnerId:            1,
		Title:              "old title",
		Status:             domain.EventStatusDraft,
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
	}

	expectedEvent := &domain.Event{
		ID:                 "event_id",
		OwnerId:            1,
		Status:             domain.EventStatusDraft,
		Title:              dto.Title,
		Description:        dto.Description,
		Type:               dto.Type,
		Tags:               dto.Tags,
		MaxParticipants:    dto.MaxParticipants,
		LocationLink:       dto.LocationLink,
		LocationUniversity: dto.LocationUniversity,
		StartDate:          dto.StartDate,
		EndDate:            dto.EndDate,
		CoverImages:        dto.CoverImages,
		AttachedImages:     dto.AttachedImages,
		AttachedFiles:      dto.AttachedFiles,
	}

	suite := newSuite(t)

	ctx := context.Background()

	suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)
	suite.mockStorage.On("UpdateEvent",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*domain.Event")).Return(
		func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
			return event, nil
		},
	)

	event, err := suite.ManagementService.UpdateEvent(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, expectedEvent, event)

	suite.mockStorage.AssertExpectations(t)

}

func TestService_UpdateEvent_GetEventError(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	dto := &dtos.UpdateEvent{
		EventId: "event_id",
		UserId:  1,
	}

	expectedErr := errors.New("get event error")
	suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(nil, expectedErr)

	_, err := suite.ManagementService.UpdateEvent(ctx, dto)
	require.ErrorIs(t, err, expectedErr)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_UpdateEvent_UserIsNotEventOwner(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	dto := &dtos.UpdateEvent{
		EventId: "event_id",
		UserId:  1,
	}

	oldEvent := &domain.Event{
		ID:      dto.EventId,
		OwnerId: dto.UserId + 1, // Different user
	}

	suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)

	_, err := suite.ManagementService.UpdateEvent(ctx, dto)
	require.ErrorIs(t, err, eventservice.ErrUserIsNotEventOwner)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_UpdateEvent_UnknownStatus(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	dto := &dtos.UpdateEvent{
		EventId: "event_id",
		UserId:  1,
	}

	oldEvent := &domain.Event{
		ID:      dto.EventId,
		OwnerId: dto.UserId,
	}

	expectedErr := eventservice.ErrUnknownStatus
	suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)

	_, err := suite.ManagementService.UpdateEvent(ctx, dto)
	require.ErrorIs(t, err, expectedErr)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_UpdateEvent_EventIsNotEditable(t *testing.T) {

	statuses := []domain.EventStatus{domain.EventStatusFinished, domain.EventStatusArchived, domain.EventStatusCanceled}

	for _, s := range statuses {
		t.Run(fmt.Sprintf("Status: %s", s), func(t *testing.T) {
			suite := newSuite(t)
			ctx := context.Background()
			dto := &dtos.UpdateEvent{
				EventId: "event_id",
				UserId:  1,
			}

			oldEvent := &domain.Event{
				ID:      dto.EventId,
				OwnerId: dto.UserId,
				Status:  s,
			}

			expectedErr := eventservice.ErrEventIsNotEditable
			suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)

			_, err := suite.ManagementService.UpdateEvent(ctx, dto)
			require.ErrorIs(t, err, expectedErr)

			suite.mockStorage.AssertExpectations(t)
		})
	}

}

func TestService_UpdateEvent_StatusApproved(t *testing.T) {
	t.Run("HasUnchangeableFields, updated status is pending", func(t *testing.T) {
		suite := newSuite(t)
		ctx := context.Background()
		dto := &dtos.UpdateEvent{
			EventId: "event_id",
			UserId:  1,
			Title:   "updated title",
			Paths:   map[string]bool{"title": true},
		}

		oldEvent := &domain.Event{
			ID:      dto.EventId,
			OwnerId: dto.UserId,
			Status:  domain.EventStatusApproved,
		}

		expectedEvent := &domain.Event{
			ID:      dto.EventId,
			OwnerId: dto.UserId,
			Status:  domain.EventStatusPending,
			Title:   "updated title",
		}

		suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)
		suite.mockStorage.On("UpdateEvent",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("*domain.Event")).Return(
			func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
				return event, nil
			},
		)

		event, err := suite.ManagementService.UpdateEvent(ctx, dto)
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, expectedEvent, event)

		suite.mockStorage.AssertExpectations(t)
	})

	t.Run("Does not have unchangeable fields, updated status stays approved", func(t *testing.T) {
		startDate := time.Now()
		endDate := time.Now().AddDate(0, 0, 20)

		suite := newSuite(t)
		ctx := context.Background()
		dto := &dtos.UpdateEvent{
			EventId:               "event_id",
			UserId:                1,
			Tags:                  []string{"updated tag1", "updated tag2"},
			StartDate:             startDate,
			EndDate:               endDate,
			LocationUniversity:    "updated location university",
			LocationLink:          "updated location link",
			IsHiddenForNonMembers: true,
			Paths: map[string]bool{
				"tags":                      true,
				"start_date":                true,
				"end_date":                  true,
				"location_university":       true,
				"location_link":             true,
				"is_hidden_for_non_members": true,
			},
		}

		oldEvent := &domain.Event{
			ID:      dto.EventId,
			OwnerId: dto.UserId,
			Type:    domain.EventTypeIntraClub,
			Status:  domain.EventStatusApproved,
		}

		expectedEvent := &domain.Event{
			ID:                    dto.EventId,
			OwnerId:               dto.UserId,
			Tags:                  []string{"updated tag1", "updated tag2"},
			StartDate:             startDate,
			EndDate:               endDate,
			Type:                  domain.EventTypeIntraClub,
			LocationUniversity:    "updated location university",
			LocationLink:          "updated location link",
			Status:                domain.EventStatusApproved,
			IsHiddenForNonMembers: true,
		}

		suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)
		suite.mockStorage.On("UpdateEvent",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("*domain.Event")).Return(
			func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
				return event, nil
			},
		)

		event, err := suite.ManagementService.UpdateEvent(ctx, dto)
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, expectedEvent, event)

		suite.mockStorage.AssertExpectations(t)

	})
}

func TestService_UpdateEvent_StatusInProgress_Pending(t *testing.T) {
	statuses := []domain.EventStatus{domain.EventStatusInProgress, domain.EventStatusPending}

	for _, s := range statuses {
		t.Run(fmt.Sprintf("Status: %s", s), func(t *testing.T) {
			t.Run("HasUnchangeableFields, return error", func(t *testing.T) {
				suite := newSuite(t)
				ctx := context.Background()
				dto := &dtos.UpdateEvent{
					EventId: "event_id",
					UserId:  1,
					Paths:   map[string]bool{"title": true},
				}

				oldEvent := &domain.Event{
					ID:      dto.EventId,
					OwnerId: dto.UserId,
					Type:    domain.EventTypeIntraClub,
					Status:  s,
				}

				expectedErr := eventservice.ErrContainsUnchangeable
				suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)

				_, err := suite.ManagementService.UpdateEvent(ctx, dto)
				require.ErrorIs(t, err, expectedErr)

				suite.mockStorage.AssertExpectations(t)
			})

			t.Run("Does not have unchangeable fields, update event", func(t *testing.T) {
				startDate := time.Now()
				endDate := time.Now().AddDate(0, 0, 20)

				suite := newSuite(t)
				ctx := context.Background()
				dto := &dtos.UpdateEvent{
					EventId:               "event_id",
					UserId:                1,
					Tags:                  []string{"updated tag1", "updated tag2"},
					StartDate:             startDate,
					EndDate:               endDate,
					LocationUniversity:    "updated location university",
					LocationLink:          "updated location link",
					IsHiddenForNonMembers: true,
					Paths: map[string]bool{
						"tags":                      true,
						"start_date":                true,
						"end_date":                  true,
						"location_university":       true,
						"location_link":             true,
						"is_hidden_for_non_members": true,
					},
				}

				oldEvent := &domain.Event{
					ID:      dto.EventId,
					OwnerId: dto.UserId,
					Type:    domain.EventTypeIntraClub,
					Status:  s,
				}

				expectedEvent := &domain.Event{
					ID:                    dto.EventId,
					OwnerId:               dto.UserId,
					Type:                  domain.EventTypeIntraClub,
					Tags:                  []string{"updated tag1", "updated tag2"},
					StartDate:             startDate,
					IsHiddenForNonMembers: true,
					EndDate:               endDate,
					LocationUniversity:    "updated location university",
					LocationLink:          "updated location link",
					Status:                s,
				}

				suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(oldEvent, nil)
				suite.mockStorage.On("UpdateEvent",
					mock.AnythingOfType("*context.timerCtx"),
					mock.AnythingOfType("*domain.Event")).Return(
					func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
						return event, nil
					},
				)

				event, err := suite.ManagementService.UpdateEvent(ctx, dto)
				require.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, expectedEvent, event)

				suite.mockStorage.AssertExpectations(t)
			})
		})
	}

}

func TestUpdateEvent_FailPath(t *testing.T) {
	oldEvent := &domain.Event{
		ID:      "test-event-id",
		OwnerId: 1,
		Status:  domain.EventStatusDraft,
		Title:   "old title",
		Type:    domain.EventTypeUniversity,
	}

	tests := []struct {
		name string
		dto  *dtos.UpdateEvent
	}{
		{
			name: "UniversityEvent_HiddenForNonMembers",
			dto: &dtos.UpdateEvent{
				EventId:               "test-event-id",
				UserId:                1,
				Type:                  domain.EventTypeUniversity,
				IsHiddenForNonMembers: true,
				Paths:                 map[string]bool{"is_hidden_for_non_members": true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)
			ctx := context.Background()

			suite.mockStorage.On("GetEvent", mock.Anything, tt.dto.EventId).Return(oldEvent, nil)

			_, err := suite.ManagementService.UpdateEvent(ctx, tt.dto)
			require.ErrorIs(t, err, eventservice.ErrEventInvalidFields)

			suite.mockStorage.AssertExpectations(t)
		})
	}
}

func TestService_DeleteEvent_HappyPath(t *testing.T) {

	tests := []struct {
		name string
		dto  *dtos.DeleteEvent
	}{
		{
			name: "User is event owner, but not admin",
			dto: &dtos.DeleteEvent{
				EventId: "event_id",
				UserId:  1,
				IsAdmin: false,
			},
		},
		{
			name: "User is admin, but not event owner",
			dto: &dtos.DeleteEvent{
				EventId: "event_id",
				UserId:  2,
				IsAdmin: true,
			},
		},
		{
			name: "User is admin and also an event owner",
			dto: &dtos.DeleteEvent{
				EventId: "event_id",
				UserId:  1,
				IsAdmin: true,
			},
		},
	}
	onGetEvent := &domain.Event{
		ID:      "event_id",
		OwnerId: 1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)
			ctx := context.Background()

			suite.mockStorage.On("GetEvent", mock.Anything, tt.dto.EventId).Return(onGetEvent, nil)
			suite.mockStorage.On("DeleteEventById", mock.Anything, tt.dto.EventId).Return(nil)

			event, err := suite.ManagementService.DeleteEvent(ctx, tt.dto)
			require.NoError(t, err)
			assert.NotNil(t, event)

			suite.mockStorage.AssertExpectations(t)
		})

	}

}

func TestService_DeleteEvent_EventNotFound(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	dto := &dtos.DeleteEvent{
		EventId: "event_id",
		UserId:  1,
		IsAdmin: true,
	}

	suite.mockStorage.On("GetEvent", mock.Anything, dto.EventId).Return(nil, storage.ErrEventNotFound)

	_, err := suite.ManagementService.DeleteEvent(ctx, dto)
	require.ErrorIs(t, err, eventservice.ErrEventNotFound)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_DeleteEvent_UserIsNotEventOwner(t *testing.T) {
	tests := []struct {
		name        string
		dto         *dtos.DeleteEvent
		expectedErr error
	}{
		{
			name: "User is not event owner and not an admin",
			dto: &dtos.DeleteEvent{
				EventId: "event_id",
				UserId:  1,
				IsAdmin: false,
			},
			expectedErr: eventservice.ErrPermissionsDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)
			ctx := context.Background()

			onGetEvent := &domain.Event{
				ID:      tt.dto.EventId,
				OwnerId: tt.dto.UserId + 1, // Different user
			}

			suite.mockStorage.On("GetEvent", mock.Anything, tt.dto.EventId).Return(onGetEvent, nil)

			_, err := suite.ManagementService.DeleteEvent(ctx, tt.dto)
			require.ErrorIs(t, err, tt.expectedErr)

			suite.mockStorage.AssertExpectations(t)
		})
	}
}

func TestService_PublishEvent_HappyPath(t *testing.T) {
	t.Run("University Scope event", func(t *testing.T) {
		suite := newSuite(t)
		ctx := context.Background()
		eventId := "event_id"
		userId := int64(1)

		onGetEvent := &domain.Event{
			ID:        eventId,
			OwnerId:   userId,
			Type:      domain.EventTypeUniversity,
			Status:    domain.EventStatusApproved,
			Title:     "old title",
			StartDate: time.Now(),
			EndDate:   time.Now(),
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
		}

		suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
		suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
			func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
				return event, nil
			},
		)

		event, err := suite.ManagementService.PublishEvent(ctx, eventId, userId)
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, domain.EventStatusInProgress, event.Status)

		suite.mockStorage.AssertExpectations(t)
	})

	t.Run("Club Scope event when status is draft", func(t *testing.T) {
		suite := newSuite(t)
		ctx := context.Background()
		eventId := "event_id"
		userId := int64(1)

		onGetEvent := &domain.Event{
			ID:        eventId,
			OwnerId:   userId,
			Type:      domain.EventTypeIntraClub,
			Status:    domain.EventStatusDraft,
			Title:     "old title",
			StartDate: time.Now(),
			EndDate:   time.Now(),
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
		}

		suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
		suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
			func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
				return event, nil
			},
		)

		event, err := suite.ManagementService.PublishEvent(ctx, eventId, userId)
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, domain.EventStatusInProgress, event.Status)

		suite.mockStorage.AssertExpectations(t)
	})
	t.Run("Club Scope event when status is approved", func(t *testing.T) {
		ctx := context.Background()
		eventId := "event_id"
		userId := int64(1)

		status := []domain.EventStatus{
			domain.EventStatusApproved,
			domain.EventStatusDraft,
		}

		for _, s := range status {
			t.Run(fmt.Sprintf("Status: %s", s), func(t *testing.T) {
				suite := newSuite(t)
				onGetEvent := &domain.Event{
					ID:        eventId,
					OwnerId:   userId,
					Type:      domain.EventTypeIntraClub,
					Status:    s,
					Title:     "old title",
					StartDate: time.Now(),
					EndDate:   time.Now(),
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
				}

				suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
				suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
					func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
						return event, nil
					},
				)

				event, err := suite.ManagementService.PublishEvent(ctx, eventId, userId)
				require.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, domain.EventStatusInProgress, event.Status)

				suite.mockStorage.AssertExpectations(t)
			})
		}

	})
}

func TestService_PublishEvent_FailPath(t *testing.T) {

	ctx := context.Background()
	eventId := "event_id"
	userId := int64(1)

	tests := []struct {
		name          string
		onGetEvent    *domain.Event
		onGetEventErr error
		expectedError error
	}{
		{
			name:          "Event Not Found",
			onGetEvent:    &domain.Event{},
			onGetEventErr: storage.ErrEventNotFound,
			expectedError: eventservice.ErrEventNotFound,
		},
		{
			name:          "User Is Not Event Owner",
			onGetEvent:    &domain.Event{ID: eventId, OwnerId: userId + 1},
			onGetEventErr: nil,
			expectedError: eventservice.ErrUserIsNotEventOwner,
		},
		{
			name:          "Event have no cover image, title, start date, end date",
			onGetEvent:    &domain.Event{ID: eventId, OwnerId: userId, Status: domain.EventStatusPending, Type: domain.EventTypeUniversity},
			onGetEventErr: nil,
			expectedError: eventservice.ErrEventInvalidFields,
		},
		{
			name: "Event is not approved",
			onGetEvent: &domain.Event{
				ID:        eventId,
				OwnerId:   userId,
				Status:    domain.EventStatusPending,
				Type:      domain.EventTypeUniversity,
				Title:     "title",
				StartDate: time.Now(),
				EndDate:   time.Now(),
				CoverImages: []domain.CoverImage{
					{File: domain.File{Name: "cover image", Url: "cover image url", Type: "image"}, Position: 1},
				},
			},
			onGetEventErr: nil,
			expectedError: eventservice.ErrEventInvalidFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)
			suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(tt.onGetEvent, tt.onGetEventErr)

			_, err := suite.ManagementService.PublishEvent(ctx, eventId, userId)
			require.ErrorIs(t, err, tt.expectedError)

			suite.mockStorage.AssertExpectations(t)
		})
	}

}

func TestService_SendToReview_HappyPath(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	eventId := "event_id"
	userId := int64(1)

	onGetEvent := &domain.Event{
		ID:        eventId,
		OwnerId:   userId,
		Type:      domain.EventTypeUniversity,
		Status:    domain.EventStatusDraft,
		Title:     "old title",
		StartDate: time.Now(),
		EndDate:   time.Now(),
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
	}

	suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
	suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
		func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
			return event, nil
		},
	)

	event, err := suite.ManagementService.SendToReview(ctx, eventId, userId)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, domain.EventStatusPending, event.Status)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_SendToReview_FailPath(t *testing.T) {
	eventId := "event_id"
	userId := int64(1)

	tests := []struct {
		name          string
		onGetEvent    *domain.Event
		onGetEventErr error
		expectedError error
	}{
		{
			name:          "Event Not Found",
			onGetEvent:    &domain.Event{},
			onGetEventErr: storage.ErrEventNotFound,
			expectedError: eventservice.ErrEventNotFound,
		},
		{
			name:          "User Is Not Event Owner",
			onGetEvent:    &domain.Event{ID: eventId, OwnerId: userId + 1},
			onGetEventErr: nil,
			expectedError: eventservice.ErrUserIsNotEventOwner,
		},
		{
			name:          "Event Is Not Draft",
			onGetEvent:    &domain.Event{ID: eventId, OwnerId: userId, Status: domain.EventStatusPending, Type: domain.EventTypeUniversity},
			onGetEventErr: nil,
			expectedError: eventservice.ErrEventInvalidFields,
		},
		{
			name: "Event Is already sent to review",
			onGetEvent: &domain.Event{
				ID:        eventId,
				OwnerId:   userId,
				Status:    domain.EventStatusPending,
				Type:      domain.EventTypeUniversity,
				Title:     "title",
				StartDate: time.Now(),
				EndDate:   time.Now(),
				CoverImages: []domain.CoverImage{
					{File: domain.File{Name: "cover image", Url: "cover image url", Type: "image"}, Position: 1},
				},
			},
			onGetEventErr: nil,
			expectedError: eventservice.ErrInvalidEventStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			suite := newSuite(t)
			suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(tt.onGetEvent, tt.onGetEventErr)

			_, err := suite.ManagementService.SendToReview(ctx, eventId, userId)
			require.ErrorIs(t, err, tt.expectedError)

			suite.mockStorage.AssertExpectations(t)
		})
	}

}

func TestService_SendToReview_InvalidStatus(t *testing.T) {
	eventId := "event_id"
	userId := int64(1)

	status := []domain.EventStatus{
		domain.EventStatusPending,
		domain.EventStatusApproved,
		domain.EventStatusInProgress,
		domain.EventStatusFinished,
		domain.EventStatusCanceled,
		domain.EventStatusArchived,
	}

	for _, s := range status {
		t.Run(fmt.Sprintf("Status: %s", s), func(t *testing.T) {
			suite := newSuite(t)
			ctx := context.Background()
			onGetEvent := &domain.Event{
				ID:        eventId,
				OwnerId:   userId,
				Status:    s,
				Type:      domain.EventTypeUniversity,
				Title:     "title",
				StartDate: time.Now(),
				EndDate:   time.Now(),
				CoverImages: []domain.CoverImage{
					{File: domain.File{Name: "cover image", Url: "cover image url", Type: "image"}, Position: 1},
				},
			}
			suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)

			_, err := suite.ManagementService.SendToReview(ctx, eventId, userId)
			require.ErrorIs(t, err, eventservice.ErrInvalidEventStatus)

			suite.mockStorage.AssertExpectations(t)
		})

	}

}

func TestService_RevokeReview_HappyPath(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	eventId := "event_id"
	userId := int64(1)

	onGetEvent := &domain.Event{
		ID:      eventId,
		OwnerId: userId,
		Status:  domain.EventStatusPending,
	}

	suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
	suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
		func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
			return event, nil
		},
	)

	event, err := suite.ManagementService.RevokeReview(ctx, eventId, userId)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, domain.EventStatusDraft, event.Status)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_RevokeReview_FailPath(t *testing.T) {
	var userId int64 = 1
	tests := []struct {
		name       string
		eventId    string
		onGetEvent *domain.Event
		onGetErr   error
		wantErr    error
	}{
		{
			name:    "RevokeReview returns error when event is not found",
			eventId: "nonexistent",
			onGetEvent: &domain.Event{
				ID:      "nonexistent",
				OwnerId: 1,
				Status:  domain.EventStatusPending,
			},
			onGetErr: storage.ErrEventNotFound,
			wantErr:  eventservice.ErrEventNotFound,
		},
		{
			name:    "RevokeReview returns error when user is not the owner",
			eventId: "event1",
			onGetEvent: &domain.Event{
				ID:      "event1",
				OwnerId: 2,
				Status:  domain.EventStatusPending,
			},
			onGetErr: nil,
			wantErr:  eventservice.ErrUserIsNotEventOwner,
		},
		{
			name:    "RevokeReview returns error when event status is not pending",
			eventId: "event1",
			onGetEvent: &domain.Event{
				ID:      "event1",
				OwnerId: 1,
				Status:  domain.EventStatusApproved,
			},
			onGetErr: nil,
			wantErr:  eventservice.ErrInvalidEventStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)

			suite.mockStorage.On("GetEvent", mock.Anything, tt.eventId).Return(tt.onGetEvent, tt.onGetErr)

			_, err := suite.ManagementService.RevokeReview(context.Background(), tt.eventId, userId)
			require.ErrorIs(t, err, tt.wantErr)

			suite.mockStorage.AssertExpectations(t)
		})
	}
}

func TestService_ApproveEvent_HappyPath(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	eventId := "event_id"
	user := domain.User{
		ID:        1,
		FirstName: "first_name",
		LastName:  "last_name",
		Barcode:   "barcode",
		AvatarURL: "url",
	}

	onGetEvent := &domain.Event{
		ID:      eventId,
		OwnerId: 1,
		Status:  domain.EventStatusPending,
	}

	suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
	suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
		func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
			return event, nil
		},
	)

	event, err := suite.ManagementService.ApproveEvent(ctx, eventId, user)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.NotNil(t, event.ApproveMetadata)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_ApproveEvent_FailPath(t *testing.T) {
	var userId int64 = 1
	tests := []struct {
		name       string
		eventId    string
		onGetEvent *domain.Event
		onGetErr   error
		wantErr    error
	}{
		{
			name:    "ApproveEvent returns error when event is not found",
			eventId: "nonexistent",
			onGetEvent: &domain.Event{
				ID:      "nonexistent",
				OwnerId: 1,
				Status:  domain.EventStatusPending,
			},
			onGetErr: storage.ErrEventNotFound,
			wantErr:  eventservice.ErrEventNotFound,
		},
		{
			name:    "ApproveEvent returns error when event status is not pending",
			eventId: "event1",
			onGetEvent: &domain.Event{
				ID:      "event1",
				OwnerId: 1,
				Status:  domain.EventStatusApproved,
			},
			onGetErr: nil,
			wantErr:  eventservice.ErrInvalidEventStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)

			suite.mockStorage.On("GetEvent", mock.Anything, tt.eventId).Return(tt.onGetEvent, tt.onGetErr)

			_, err := suite.ManagementService.ApproveEvent(context.Background(), tt.eventId, domain.User{ID: userId})
			require.ErrorIs(t, err, tt.wantErr)

			suite.mockStorage.AssertExpectations(t)
		})
	}
}

func TestService_RejectEvent_HappyPath(t *testing.T) {
	suite := newSuite(t)
	ctx := context.Background()
	eventId := "event_id"
	dto := dtos.RejectEvent{
		EventId: eventId,
		User: domain.User{
			ID:        1,
			FirstName: "first_name",
			LastName:  "last_name",
			Barcode:   "barcode",
			AvatarURL: "url",
		},
		Reason: "reason",
	}

	onGetEvent := &domain.Event{
		ID:      eventId,
		OwnerId: 1,
		Status:  domain.EventStatusPending,
	}

	suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
	suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
		func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
			return event, nil
		},
	)

	event, err := suite.ManagementService.RejectEvent(ctx, &dto)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.NotNil(t, event.RejectMetadata)
	assert.Equal(t, domain.EventStatusRejected, event.Status)

	suite.mockStorage.AssertExpectations(t)
}

func TestService_RejectEvent_FailPath(t *testing.T) {
	var userId int64 = 1
	tests := []struct {
		name       string
		eventId    string
		onGetEvent *domain.Event
		onGetErr   error
		wantErr    error
	}{
		{
			name:    "RejectEvent returns error when event is not found",
			eventId: "nonexistent",
			onGetEvent: &domain.Event{
				ID:      "nonexistent",
				OwnerId: 1,
				Status:  domain.EventStatusPending,
			},
			onGetErr: storage.ErrEventNotFound,
			wantErr:  eventservice.ErrEventNotFound,
		},
		{
			name:    "RejectEvent returns error when event status is not pending",
			eventId: "event1",
			onGetEvent: &domain.Event{
				ID:      "event1",
				OwnerId: 1,
				Status:  domain.EventStatusApproved,
			},
			onGetErr: nil,
			wantErr:  eventservice.ErrInvalidEventStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)

			suite.mockStorage.On("GetEvent", mock.Anything, tt.eventId).Return(tt.onGetEvent, tt.onGetErr)

			_, err := suite.ManagementService.RejectEvent(context.Background(), &dtos.RejectEvent{EventId: tt.eventId, User: domain.User{ID: userId}})
			require.ErrorIs(t, err, tt.wantErr)

			suite.mockStorage.AssertExpectations(t)
		})
	}
}

func TestService_UnpublishEvent_HappyPath(t *testing.T) {
	t.Run("University Scope event", func(t *testing.T) {
		suite := newSuite(t)
		ctx := context.Background()
		eventId := "event_id"
		userId := int64(1)

		onGetEvent := &domain.Event{
			ID:      eventId,
			OwnerId: userId,
			Status:  domain.EventStatusInProgress,
			Type:    domain.EventTypeUniversity,
		}

		suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
		suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
			func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
				return event, nil
			},
		)

		event, err := suite.ManagementService.UnpublishEvent(ctx, eventId, userId)
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, domain.EventStatusApproved, event.Status)

		suite.mockStorage.AssertExpectations(t)
	})

	t.Run("Club Scope event", func(t *testing.T) {
		suite := newSuite(t)
		ctx := context.Background()
		eventId := "event_id"
		userId := int64(1)

		onGetEvent := &domain.Event{
			ID:      eventId,
			OwnerId: userId,
			Status:  domain.EventStatusInProgress,
			Type:    domain.EventTypeIntraClub,
		}

		suite.mockStorage.On("GetEvent", mock.Anything, eventId).Return(onGetEvent, nil)
		suite.mockStorage.On("UpdateEvent", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Event")).Return(
			func(ctx context.Context, event *domain.Event) (*domain.Event, error) {
				return event, nil
			},
		)

		event, err := suite.ManagementService.UnpublishEvent(ctx, eventId, userId)
		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, domain.EventStatusDraft, event.Status)

		suite.mockStorage.AssertExpectations(t)
	})

}

func TestService_UnpublishEvent_FailPath(t *testing.T) {
	var userId int64 = 1
	tests := []struct {
		name       string
		eventId    string
		onGetEvent *domain.Event
		onGetErr   error
		wantErr    error
	}{
		{
			name:    "UnpublishEvent returns error when event is not found",
			eventId: "nonexistent",
			onGetEvent: &domain.Event{
				ID:      "event_id",
				OwnerId: userId,
				Status:  domain.EventStatusInProgress,
			},
			onGetErr: storage.ErrEventNotFound,
			wantErr:  eventservice.ErrEventNotFound,
		},
		{
			name:    "UnpublishEvent returns error when user is not the owner",
			eventId: "event_id",
			onGetEvent: &domain.Event{
				ID:      "event_id",
				OwnerId: userId + 1, // Different user
				Status:  domain.EventStatusInProgress,
			},
			onGetErr: nil,
			wantErr:  eventservice.ErrUserIsNotEventOwner,
		},
		{
			name:    "UnpublishEvent returns error when event status is not in progress",
			eventId: "event_id",
			onGetEvent: &domain.Event{
				ID:      "event_id",
				OwnerId: userId,
				Status:  domain.EventStatusApproved, // Not in progress
			},
			onGetErr: nil,
			wantErr:  domain.ErrEventIsNotPublished,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := newSuite(t)

			suite.mockStorage.On("GetEvent", mock.Anything, tt.eventId).Return(tt.onGetEvent, tt.onGetErr)

			_, err := suite.ManagementService.UnpublishEvent(context.Background(), tt.eventId, userId)
			require.ErrorIs(t, err, tt.wantErr)

			suite.mockStorage.AssertExpectations(t)
		})
	}
}
