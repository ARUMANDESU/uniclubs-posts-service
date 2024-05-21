package management

import (
	"fmt"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/protobuf/field_mask"
	"testing"
	"time"
)

func TestManagement_SendToReview_HappyPath(t *testing.T) {
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
	updateReq := &eventv1.UpdateEventRequest{
		EventId:   createResp.GetId(),
		UserId:    createResp.GetOwnerId(),
		Type:      domain.EventTypeUniversity.String(),
		Title:     gofakeit.AppName(),
		StartDate: time.Now().Format(domain.TimeLayout),
		EndDate:   time.Now().AddDate(0, 0, 10).Format(domain.TimeLayout),
		CoverImages: []*eventv1.CoverImage{
			{
				Url:      gofakeit.URL(),
				Name:     gofakeit.AppName(),
				Type:     "image/png",
				Position: 1,
			},
		},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"title", "start_date", "end_date", "cover_images", "type"}},
	}

	updateResp, err := st.EventClient.UpdateEvent(ctx, updateReq)
	require.NoError(t, err, fmt.Sprintf("UpdateEvent failed, error: %v", err))

	//Send to review
	sendToReviewReq := &eventv1.EventActionRequest{
		EventId: updateResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	sendToReviewResp, err := st.EventClient.SendToReview(ctx, sendToReviewReq)
	require.NoError(t, err, "SendToReview failed")
	assert.NotEqual(t, "", sendToReviewResp.GetId(), "SendToReview failed: event ID is zero")
	assert.Equal(t, domain.EventStatusPending.String(), sendToReviewResp.GetStatus(), "SendToReview failed: event status is not pending")
}

func TestManagement_SendToReview_InvalidEventId(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx, st := suite.New(t)

	// Send to review with invalid event ID
	sendToReviewReq := &eventv1.EventActionRequest{
		EventId: gofakeit.UUID(),
		UserId:  gofakeit.Int64(),
	}

	_, err := st.EventClient.SendToReview(ctx, sendToReviewReq)
	require.Error(t, err, "SendToReview should fail with invalid event ID")
}

func TestManagement_RevokeEvent_HappyPah(t *testing.T) {
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
	updateReq := &eventv1.UpdateEventRequest{
		EventId:   createResp.GetId(),
		UserId:    createResp.GetOwnerId(),
		Type:      domain.EventTypeUniversity.String(),
		Title:     gofakeit.AppName(),
		StartDate: time.Now().Format(domain.TimeLayout),
		EndDate:   time.Now().AddDate(0, 0, 10).Format(domain.TimeLayout),
		CoverImages: []*eventv1.CoverImage{
			{
				Url:      gofakeit.URL(),
				Name:     gofakeit.AppName(),
				Type:     "image/png",
				Position: 1,
			},
		},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"title", "start_date", "end_date", "cover_images", "type"}},
	}

	updateResp, err := st.EventClient.UpdateEvent(ctx, updateReq)
	require.NoError(t, err, fmt.Sprintf("UpdateEvent failed, error: %v", err))

	//Send to review
	sendToReviewReq := &eventv1.EventActionRequest{
		EventId: updateResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	sendToReviewResp, err := st.EventClient.SendToReview(ctx, sendToReviewReq)
	require.NoError(t, err, "SendToReview failed")

	// Revoke the created event
	revokeReq := &eventv1.EventActionRequest{
		EventId: sendToReviewResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	revokeResp, err := st.EventClient.RevokeReview(ctx, revokeReq)
	require.NoError(t, err, "RevokeEvent failed")
	assert.NotEqual(t, "", revokeResp.GetId(), "RevokeEvent failed: event ID is zero")
	assert.Equal(t, domain.EventStatusDraft.String(), revokeResp.GetStatus(), "RevokeEvent failed: event status is not rejected")
	assert.NotNil(t, revokeResp.GetRejectMetadata().GetRejectedAt(), "RevokeEvent failed: rejected at is nil")
	assert.NotNil(t, revokeResp.GetRejectMetadata().GetRejectedBy(), "RevokeEvent failed: rejected by is nil")
	assert.NotNil(t, revokeResp.GetRejectMetadata().GetReason(), "RevokeEvent failed: reason is nil")
}
