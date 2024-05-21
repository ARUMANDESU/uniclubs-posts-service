package info

import (
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestInfo_GetEvent(t *testing.T) {
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

	req := &eventv1.GetEventRequest{
		EventId: createResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	resp, err := st.EventClient.GetEvent(ctx, req)
	require.NoError(t, err, fmt.Sprintf("GetEvent failed, error: %v", err))
	assert.NotEqual(t, "", resp.GetId(), "GetEvent failed: event ID is zero")
}

func TestInfo_GetEvent_Invalid(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, st := suite.New(t)

	t.Run("Invalid event id", func(t *testing.T) {
		expectedErr := codes.InvalidArgument
		req := &eventv1.GetEventRequest{
			EventId: "123",
			UserId:  123,
		}

		_, err := st.EventClient.GetEvent(ctx, req)
		require.Error(t, err)
		assert.Equal(t, expectedErr, status.Code(err), fmt.Sprintf("expected code %v, got %v", expectedErr, status.Code(err)))
	})

	t.Run("Not found", func(t *testing.T) {
		expectedErr := codes.NotFound
		req := &eventv1.GetEventRequest{
			EventId: "664b638aa22861e6811b3c10",
			UserId:  123,
		}

		_, err := st.EventClient.GetEvent(ctx, req)
		require.Error(t, err)
		assert.Equal(t, expectedErr, status.Code(err), fmt.Sprintf("expected code %v, got %v", expectedErr, status.Code(err)))
	})

	t.Run("User not organizer", func(t *testing.T) {
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

		req := &eventv1.GetEventRequest{
			EventId: createResp.GetId(),
			UserId:  createResp.GetOwnerId() + 1, // different user
		}

		_, err = st.EventClient.GetEvent(ctx, req)
		require.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err), fmt.Sprintf("expected code %v, got %v", codes.NotFound, status.Code(err)))

	})
}

func TestInfo_ListEvents(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, st := suite.New(t)

	// create events to list
	for i := 0; i < 20; i++ {
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

		_, err := st.EventClient.CreateEvent(ctx, createReq)
		require.NoError(t, err, fmt.Sprintf("CreateEvent failed, error: %v", err))
	}

}
