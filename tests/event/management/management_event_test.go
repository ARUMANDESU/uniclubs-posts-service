package management

import (
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestManagement_CreateEvent(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
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

func TestManagement_CreateEvent_Invalid(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, st := suite.New(t)

	tests := []struct {
		name string
		req  *eventv1.CreateEventRequest
		code codes.Code
	}{
		{
			name: "empty club",
			req: &eventv1.CreateEventRequest{
				Club: nil,
				User: &eventv1.UserObject{
					Id:        gofakeit.Int64(),
					FirstName: gofakeit.FirstName(),
					LastName:  gofakeit.LastName(),
					Barcode:   gofakeit.UUID(),
				},
			},
			code: codes.InvalidArgument,
		},
		{
			name: "empty user",
			req: &eventv1.CreateEventRequest{
				Club: &eventv1.ClubObject{
					Id:   gofakeit.Int64(),
					Name: gofakeit.AppName(),
				},
				User: nil,
			},
			code: codes.InvalidArgument,
		},
		{
			name: "empty club and user",
			req: &eventv1.CreateEventRequest{
				Club: nil,
				User: nil,
			},
			code: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.EventClient.CreateEvent(ctx, tt.req)
			require.Error(t, err)
			assert.Equal(t, tt.code, status.Code(err), fmt.Sprintf("expected code %v, got %v", tt.code, status.Code(err)))
		})

	}
}

func TestManagement_Create_UpdateEvent(t *testing.T) {
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

func TestManagement_DeleteEvent(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, st := suite.New(t)

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

	req := &eventv1.DeleteEventRequest{
		EventId: createResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	resp, err := st.EventClient.DeleteEvent(ctx, req)
	require.NoError(t, err, fmt.Sprintf("DeleteEvent failed, error: %v", err))
	assert.NotEqual(t, "", resp.GetId(), "DeleteEvent failed: event ID is zero")
}

func TestManagement_DeleteEvent_Invalid(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx, st := suite.New(t)

	tests := []struct {
		name string
		req  *eventv1.DeleteEventRequest
		code codes.Code
	}{
		{
			name: "empty event ID",
			req: &eventv1.DeleteEventRequest{
				EventId: "",
				UserId:  gofakeit.Int64(),
			},
			code: codes.InvalidArgument,
		},
		{
			name: "empty user ID",
			req: &eventv1.DeleteEventRequest{
				EventId: gofakeit.UUID(),
				UserId:  0,
			},
			code: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.EventClient.DeleteEvent(ctx, tt.req)
			require.Error(t, err)
			assert.Equal(t, tt.code, status.Code(err), fmt.Sprintf("expected code %v, got %v", tt.code, status.Code(err)))
		})

	}
}
