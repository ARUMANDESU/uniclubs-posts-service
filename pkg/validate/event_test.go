package validate

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
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
		Type:               "university",
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
