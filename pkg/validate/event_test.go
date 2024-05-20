package validate

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateEvent_ValidRequest(t *testing.T) {
	req := &eventv1.CreateEventRequest{
		Club: &eventv1.ClubObject{
			Id:      1,
			Name:    "Test Club",
			LogoUrl: "",
		},
		User: &eventv1.UserObject{
			Id:        1,
			FirstName: "FirstName",
			LastName:  "LastName",
			Barcode:   "Barcode",
			AvatarUrl: "AvatarUrl",
		},
	}
	err := CreateEvent(req)
	assert.Nil(t, err)
}

func TestCreateEvent_InvalidType(t *testing.T) {
	req := "Invalid Type"
	err := CreateEvent(req)
	assert.NotNil(t, err)
}

func TestGetEventWithValidRequest(t *testing.T) {
	req := &eventv1.GetEventRequest{
		EventId: "Test Event",
		UserId:  1,
	}
	err := GetEvent(req)
	assert.Nil(t, err)
}

func TestGetEventWithInvalidType(t *testing.T) {
	req := "Invalid Type"
	err := GetEvent(req)
	assert.NotNil(t, err)
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

func TestUpdateEventWithInvalidType(t *testing.T) {
	req := "Invalid Type"
	err := UpdateEvent(req)
	assert.NotNil(t, err)
}
