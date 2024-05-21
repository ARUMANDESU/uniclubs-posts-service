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

func TestManagement_RejectEvent_HappyPath(t *testing.T) {
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

	// Reject the created event
	rejectReq := &eventv1.RejectEventRequest{
		EventId: sendToReviewResp.GetId(),
		User: &eventv1.UserObject{
			Id:        gofakeit.Int64(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Barcode:   gofakeit.UUID(),
		},
		Reason: gofakeit.Sentence(10),
	}

	rejectResp, err := st.EventClient.RejectEvent(ctx, rejectReq)
	require.NoError(t, err, fmt.Sprintf("RejectEvent failed, error: %v", err))
	assert.Equal(t, domain.EventStatusRejected.String(), rejectResp.GetStatus())
	assert.Equal(t, rejectReq.GetReason(), rejectResp.GetRejectMetadata().Reason)
	assert.NotEmpty(t, rejectResp.GetRejectMetadata())
}

func TestManagement_RejectEvent_AlreadyApproved(t *testing.T) {
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

	// Approve the  event
	approveReq := &eventv1.ApproveEventRequest{
		EventId: sendToReviewResp.GetId(),
		User: &eventv1.UserObject{
			Id:        gofakeit.Int64(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Barcode:   gofakeit.UUID(),
		},
	}

	_, err = st.EventClient.ApproveEvent(ctx, approveReq)
	require.NoError(t, err, "ApproveEvent failed")

	// Reject the approved event
	rejectReq := &eventv1.RejectEventRequest{
		EventId: sendToReviewResp.GetId(),
		User: &eventv1.UserObject{
			Id:        gofakeit.Int64(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Barcode:   gofakeit.UUID(),
		},
		Reason: gofakeit.Sentence(10),
	}

	_, err = st.EventClient.RejectEvent(ctx, rejectReq)
	require.Error(t, err, "RejectEvent should fail")
}
