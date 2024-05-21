package management

import (
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/protobuf/field_mask"
	"testing"
)

func TestManagement_UpdateEvent_HappyPath(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, st := suite.New(t)

	// Create an event
	createReq := &eventv1.CreateEventRequest{
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

	createResp, err := st.EventClient.CreateEvent(ctx, createReq)
	require.NoError(t, err, fmt.Sprintf("CreateEvent failed, error: %v", err))

	// Update the created event
	req := &eventv1.UpdateEventRequest{
		EventId:    createResp.GetId(),
		UserId:     createResp.GetOwnerId(),
		Title:      gofakeit.AppName(),
		UpdateMask: &field_mask.FieldMask{Paths: []string{"title"}},
	}

	resp, err := st.EventClient.UpdateEvent(ctx, req)
	require.NoError(t, err, fmt.Sprintf("UpdateEvent failed, error: %v", err))
	assert.NotEqual(t, "", resp.GetId(), "UpdateEvent failed: event ID is zero")

	// Check if the event was updated
	assert.Equal(t, req.GetTitle(), resp.GetTitle())
}
