package event

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"testing"
)

func TestManagement_CreateEvent(t *testing.T) {
	ctx, st := suite.New(t)

	req := &eventv1.CreateEventRequest{
		Club: &eventv1.ClubObject{
			Id:   gofakeit.Int64(),
			Name: gofakeit.AppName(),
		},
		User: &eventv1.UserObject{
			Id:        gofakeit.Int64(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Barcode:   gofakeit.UUID(),
		},
	}

	resp, err := st.EventClient.CreateEvent(ctx, req)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	// Check if the event was created
	if resp.GetId() == "" {
		t.Fatalf("CreateEvent failed: event ID is zero")
	}
}
