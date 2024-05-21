package validate

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestCreateEvent_ValidRequest(t *testing.T) {
	req := &eventv1.CreateEventRequest{
		Club: newClub(t),
		User: newUser(t),
	}
	err := CreateEvent(req)
	assert.Nil(t, err)
}

func TestCreateEvent_InvalidType(t *testing.T) {
	req := "Invalid Type"
	err := CreateEvent(req)
	assert.NotNil(t, err)
}
func TestCreateEvent_Invalid(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.CreateEventRequest
	}{
		{
			name: "Empty Club",
			req: &eventv1.CreateEventRequest{
				User: newUser(t),
			},
		},
		{
			name: "Empty User",
			req: &eventv1.CreateEventRequest{
				Club: newClub(t),
			},
		},
		{
			name: "Empty Club and User",
			req:  &eventv1.CreateEventRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateEvent(tt.req)
			assert.NotNil(t, err)
		})
	}
}

func TestGetEvent_ValidRequest(t *testing.T) {
	req := &eventv1.GetEventRequest{
		EventId: "Test Event",
		UserId:  1,
	}
	err := GetEvent(req)
	assert.Nil(t, err)
}

func TestGetEvent_InvalidType(t *testing.T) {
	req := "Invalid Type"
	err := GetEvent(req)
	assert.NotNil(t, err)
}

func TestGetEvent_Invalid(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.GetEventRequest
	}{
		{
			name: "Empty EventId",
			req: &eventv1.GetEventRequest{
				UserId: 1,
			},
		},
		{
			name: "Empty EventId and UserId",
			req:  &eventv1.GetEventRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetEvent(tt.req)
			assert.NotNil(t, err)
		})
	}
}

func TestUpdateEventWithValidRequest(t *testing.T) {
	req := &eventv1.UpdateEventRequest{
		EventId:            gofakeit.UUID(),
		UserId:             1,
		Title:              "Test Title",
		Description:        "Test Description",
		Type:               domain.EventTypeUniversity.String(),
		Tags:               []string{"tag1", "tag2"},
		MaxParticipants:    100,
		StartDate:          "2022-01-01T00:00:00Z",
		EndDate:            "2022-12-31T23:59:59Z",
		LocationLink:       gofakeit.URL(),
		LocationUniversity: "Test University",
		CoverImages: []*eventv1.CoverImage{
			{
				Url:      gofakeit.URL(),
				Name:     gofakeit.ProductName(),
				Type:     "image/jpeg",
				Position: 1,
			},
			{
				Url:      gofakeit.URL(),
				Name:     gofakeit.ProductName(),
				Type:     "image/jpeg",
				Position: 2,
			},
		},
		AttachedImages: []*eventv1.FileObject{
			{
				Url:  gofakeit.URL(),
				Name: gofakeit.ProductName(),
				Type: "image/jpeg",
			},
			{
				Url:  gofakeit.URL(),
				Name: gofakeit.ProductName(),
				Type: "image/jpeg",
			},
		},
		AttachedFiles: []*eventv1.FileObject{
			{
				Url:  gofakeit.URL(),
				Name: gofakeit.ProductName(),
				Type: "file/pdf",
			},
			{
				Url:  gofakeit.URL(),
				Name: gofakeit.ProductName(),
				Type: "file/lol",
			},
		},
	}
	err := UpdateEvent(req)
	assert.Nil(t, err)
}

func TestUpdateEvent_InvalidType(t *testing.T) {
	req := "Invalid Type"
	err := UpdateEvent(req)
	assert.NotNil(t, err)
}

func TestUpdateEvent_Invalid(t *testing.T) {
	tenYearsFromNow := time.Now().AddDate(10, 1, 0).Format(time.RFC3339)

	tests := []struct {
		name string
		req  *eventv1.UpdateEventRequest
	}{
		{
			name: "Empty EventId",
			req: &eventv1.UpdateEventRequest{
				UserId: 1,
			},
		},
		{
			name: "Empty UserId",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
			},
		},
		{
			name: "Empty EventId and UserId",
			req:  &eventv1.UpdateEventRequest{},
		},
		{
			name: "Title Length > Max",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  1,
				Title:   gofakeit.Paragraph(1, 5, 100, " "),
			},
		},
		/*{ // This test is commented out because this crashes the test, your computer, and the universe
			name: "Description Length > Max",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  1,
				Title:   gofakeit.Paragraph(40, 1000, 35000, " "),
			},
		},*/
		{
			name: "Invalid Type",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  1,
				Type:    "Invalid Type",
			},
		},
		{
			name: "Tags Length < Min",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  1,
				Tags:    []string{"t"},
			},
		},
		{
			name: "Tags Length > Max",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  1,
				Tags:    []string{gofakeit.Paragraph(1, 5, 100, " ")},
			},
		},
		{
			name: "MaxParticipants < 0",
			req: &eventv1.UpdateEventRequest{
				EventId:         gofakeit.UUID(),
				UserId:          1,
				MaxParticipants: -1,
			},
		},
		{
			name: "MaxParticipants > Max",
			req: &eventv1.UpdateEventRequest{
				EventId:         gofakeit.UUID(),
				UserId:          1,
				MaxParticipants: 100001,
			},
		},
		{
			name: "StartDate > 10 years",
			req: &eventv1.UpdateEventRequest{
				EventId:   gofakeit.UUID(),
				UserId:    1,
				StartDate: tenYearsFromNow,
			},
		},
		{
			name: "StartDate < 6 years",
			req: &eventv1.UpdateEventRequest{
				EventId:   gofakeit.UUID(),
				UserId:    1,
				StartDate: "2015-01-01T00:00:00Z",
			},
		},
		{
			name: "EndDate > 10 years",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  1,
				EndDate: tenYearsFromNow,
			},
		},
		{
			name: "EndDate < 6 years",
			req: &eventv1.UpdateEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  1,
				EndDate: "2015-12-31T23:59:59Z",
			},
		},
		{
			name: "LocationLink Length > Max",
			req: &eventv1.UpdateEventRequest{
				EventId:      gofakeit.UUID(),
				UserId:       1,
				LocationLink: gofakeit.Paragraph(1, 5, 2501, " "),
			},
		},
		{
			name: "LocationUniversity Length > Max",
			req: &eventv1.UpdateEventRequest{
				EventId:            gofakeit.UUID(),
				UserId:             1,
				LocationUniversity: gofakeit.Paragraph(1, 5, 251, " "),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateEvent(tt.req)
			assert.NotNil(t, err)
		})
	}
}

func TestListEvents_ValidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.ListEventsRequest
	}{
		{
			name: "Empty Request, except page number and page size",
			req: &eventv1.ListEventsRequest{
				PageNumber: 1,
				PageSize:   10,
			},
		},
		{
			name: "Full",
			req: &eventv1.ListEventsRequest{
				Query:      "test",
				SortBy:     "date",
				SortOrder:  "asc",
				PageNumber: 1,
				PageSize:   10,
				Filter:     &eventv1.EventFilter{UserId: 1},
			},
		},
		{
			name: "No Filter",
			req: &eventv1.ListEventsRequest{
				Query:      "test",
				SortBy:     "date",
				SortOrder:  "asc",
				PageNumber: 1,
				PageSize:   10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ListEvents(tt.req)
			assert.Nil(t, err)
		})

	}
}

func TestListEvents_InvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.ListEventsRequest
	}{
		{
			name: "Invalid Query Length",
			req: &eventv1.ListEventsRequest{
				Query:      gofakeit.Paragraph(1, 1, 1001, ""),
				PageNumber: 1,
				PageSize:   10,
			},
		},
		{
			name: "Invalid SortBy Value",
			req: &eventv1.ListEventsRequest{
				SortBy:     "invalid",
				SortOrder:  "asc",
				PageNumber: 1,
				PageSize:   10,
			},
		},
		{
			name: "Invalid SortOrder Value",
			req: &eventv1.ListEventsRequest{
				SortBy:     "date",
				SortOrder:  "invalid",
				PageNumber: 1,
				PageSize:   10,
			},
		},
		{
			name: "Missing PageNumber",
			req: &eventv1.ListEventsRequest{
				SortBy:    "date",
				SortOrder: "asc",
				PageSize:  10,
			},
		},
		{
			name: "Missing PageSize",
			req: &eventv1.ListEventsRequest{
				PageNumber: 1,
			},
		},
		{
			name: "Sort order is empty when sort by is set",
			req: &eventv1.ListEventsRequest{
				SortBy:     "date",
				PageNumber: 1,
				PageSize:   10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ListEvents(tt.req)
			assert.Error(t, err)
		})
	}
}

func TestPublishEvent_ValidRequest(t *testing.T) {
	t.Run("University Scope Event", func(t *testing.T) {
		event := &domain.Event{
			Status:    domain.EventStatusApproved,
			Title:     "Test Title",
			Type:      domain.EventTypeUniversity,
			StartDate: time.Now(),
			EndDate:   time.Now(),
			CoverImages: []domain.CoverImage{
				{
					File: domain.File{
						Url:  "https://example.com/image.jpg",
						Name: "Test Image",
						Type: "image/jpeg",
					},
					Position: 1,
				},
			},
		}

		err := PublishEvent(event)
		assert.Nil(t, err)
	})

	t.Run("Club Scope Event", func(t *testing.T) {
		event := &domain.Event{
			Status:    domain.EventStatusDraft,
			Title:     "Test Title",
			Type:      domain.EventTypeIntraClub,
			StartDate: time.Now(),
			EndDate:   time.Now(),
			CoverImages: []domain.CoverImage{
				{
					File: domain.File{
						Url:  "https://example.com/image.jpg",
						Name: "Test Image",
						Type: "image/jpeg",
					},
					Position: 1,
				},
			},
		}

		err := PublishEvent(event)
		assert.Nil(t, err)
	})

}

func TestPublishEvent_InvalidType(t *testing.T) {
	event := "Invalid Type"
	err := PublishEvent(event)
	assert.NotNil(t, err)
}

func TestPublishEvent_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		event *domain.Event
	}{
		{
			name: "Invalid Status",
			event: &domain.Event{
				Status: "Invalid Status",
			},
		},
		{
			name: "Invalid Title",
			event: &domain.Event{
				Status: domain.EventStatusDraft,
				Title:  "",
			},
		},
		{
			name: "Invalid Type",
			event: &domain.Event{
				Status: domain.EventStatusDraft,
				Title:  "Test Title",
				Type:   "Invalid Type",
			},
		},
		{
			name: "Invalid StartDate",
			event: &domain.Event{
				Status: domain.EventStatusDraft,
				Title:  "Test Title",
				Type:   domain.EventTypeUniversity,
			},
		},
		{
			name: "Invalid EndDate",
			event: &domain.Event{
				Status:    domain.EventStatusDraft,
				Title:     "Test Title",
				Type:      domain.EventTypeUniversity,
				StartDate: time.Now(),
			},
		},
		{
			name: "Invalid CoverImages",
			event: &domain.Event{
				Status:      domain.EventStatusDraft,
				Title:       "Test Title",
				Type:        domain.EventTypeUniversity,
				StartDate:   time.Now(),
				EndDate:     time.Now(),
				CoverImages: []domain.CoverImage{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PublishEvent(tt.event)
			assert.NotNil(t, err)
		})
	}
}

func TestPublishEvent_UniversityEvent_InvalidStatus(t *testing.T) {
	event := &domain.Event{
		Title:     "Test Title",
		Type:      domain.EventTypeUniversity,
		StartDate: time.Now(),
		EndDate:   time.Now(),
		CoverImages: []domain.CoverImage{
			{
				File: domain.File{
					Url:  "https://example.com/image.jpg",
					Name: "Test Image",
					Type: "image/jpeg",
				},
				Position: 1,
			},
		},
	}

	tests := []struct {
		name   string
		status domain.EventStatus
	}{
		{
			name:   "University Event with Draft Status",
			status: domain.EventStatusDraft,
		},
		{
			name:   "University Event with Pending Status",
			status: domain.EventStatusPending,
		},
		{
			name:   "University Event with In Progress Status",
			status: domain.EventStatusInProgress,
		},
		{
			name:   "University Event with Finished Status",
			status: domain.EventStatusFinished,
		},
		{
			name:   "University Event with Canceled Status",
			status: domain.EventStatusCanceled,
		},
		{
			name:   "University Event with Rejected Status",
			status: domain.EventStatusRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := event
			event.Status = tt.status
			err := PublishEvent(event)
			assert.NotNil(t, err)
		})
	}
}

func TestPublishEvent_IntraClubEvent_InvalidStatus(t *testing.T) {
	event := &domain.Event{
		Title:     "Test Title",
		Type:      domain.EventTypeIntraClub,
		StartDate: time.Now(),
		EndDate:   time.Now(),
		CoverImages: []domain.CoverImage{
			{
				File: domain.File{
					Url:  "https://example.com/image.jpg",
					Name: "Test Image",
					Type: "image/jpeg",
				},
				Position: 1,
			},
		},
	}

	tests := []struct {
		name   string
		status domain.EventStatus
	}{
		{
			name:   "Club scope with Pending Status",
			status: domain.EventStatusPending,
		},
		{
			name:   "Club scope with In Progress Status",
			status: domain.EventStatusInProgress,
		},
		{
			name:   "Club scope with Finished Status",
			status: domain.EventStatusFinished,
		},
		{
			name:   "Club scope with Canceled Status",
			status: domain.EventStatusCanceled,
		},
		{
			name:   "Club scope with Rejected Status",
			status: domain.EventStatusRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := event
			event.Status = tt.status
			err := PublishEvent(event)
			assert.NotNil(t, err)
		})
	}
}

func TestSendToReview_HappyPath(t *testing.T) {
	tests := []struct {
		name    string
		event   *domain.Event
		wantErr bool
	}{
		{
			name: "ValidEvent",
			event: &domain.Event{
				Title:       "Test Title",
				Type:        domain.EventTypeUniversity,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(24 * time.Hour),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: false,
		},
		{
			name: "InvalidType",
			event: &domain.Event{
				Title:       "Test Title",
				Type:        "Invalid Type",
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(24 * time.Hour),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: true,
		},
		{
			name: "NoCoverImages",
			event: &domain.Event{
				Title:     "Test Title",
				Type:      domain.EventTypeUniversity,
				StartDate: time.Now(),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "TitleTooShort",
			event: &domain.Event{
				Title:       "T",
				Type:        domain.EventTypeUniversity,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(24 * time.Hour),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: true,
		},
		{
			name: "TitleTooLong",
			event: &domain.Event{
				Title:       strings.Repeat("a", 501),
				Type:        domain.EventTypeUniversity,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(24 * time.Hour),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendToReview(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendToReview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendToReview_FailPath(t *testing.T) {
	tests := []struct {
		name    string
		event   *domain.Event
		wantErr bool
	}{
		{
			name: "MissingTitle",
			event: &domain.Event{
				Type:        domain.EventTypeUniversity,
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(24 * time.Hour),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: true,
		},
		{
			name: "MissingType",
			event: &domain.Event{
				Title:       "Test Title",
				StartDate:   time.Now(),
				EndDate:     time.Now().Add(24 * time.Hour),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: true,
		},
		{
			name: "MissingStartDate",
			event: &domain.Event{
				Title:       "Test Title",
				Type:        domain.EventTypeUniversity,
				EndDate:     time.Now().Add(24 * time.Hour),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: true,
		},
		{
			name: "MissingEndDate",
			event: &domain.Event{
				Title:       "Test Title",
				Type:        domain.EventTypeUniversity,
				StartDate:   time.Now(),
				CoverImages: []domain.CoverImage{{}},
			},
			wantErr: true,
		},
		{
			name: "MissingCoverImages",
			event: &domain.Event{
				Title:     "Test Title",
				Type:      domain.EventTypeUniversity,
				StartDate: time.Now(),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendToReview(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendToReview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
