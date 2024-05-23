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

func TestPublishEvent_ClubScope_HappyPath(t *testing.T) {
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
		Type:      domain.EventTypeIntraClub.String(),
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

	// Publish the created event
	publishReq := &eventv1.EventActionRequest{
		EventId: updateResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	publishResp, err := st.EventClient.PublishEvent(ctx, publishReq)
	require.NoError(t, err, "PublishEvent failed")

	// Assert that the event is published
	assert.Equal(t, domain.EventStatusInProgress.String(), publishResp.GetStatus())
}

func TestPublishEvent_UniversityScope_HappyPath(t *testing.T) {
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

	// Approve the created event
	approveReq := &eventv1.ApproveEventRequest{
		EventId: sendToReviewResp.GetId(),
		User: &eventv1.UserObject{
			Id:        gofakeit.Int64(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Barcode:   gofakeit.UUID(),
		},
	}

	approveResp, err := st.EventClient.ApproveEvent(ctx, approveReq)
	require.NoError(t, err, "ApproveEvent failed")
	require.Equal(t, domain.EventStatusApproved.String(), approveResp.GetStatus())
	require.Equal(t, approveReq.GetUser().GetId(), approveResp.GetApproveMetadata().GetApprovedBy().GetId())

	// Publish the created event
	publishReq := &eventv1.EventActionRequest{
		EventId: approveResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	publishResp, err := st.EventClient.PublishEvent(ctx, publishReq)
	require.NoError(t, err, "PublishEvent failed")

	// Assert that the event is published
	assert.Equal(t, domain.EventStatusInProgress.String(), publishResp.GetStatus())
}

func TestPublishEvent_UniversityScope_Unauthorized(t *testing.T) {
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

	// Approve the created event
	approveReq := &eventv1.ApproveEventRequest{
		EventId: sendToReviewResp.GetId(),
		User: &eventv1.UserObject{
			Id:        gofakeit.Int64(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Barcode:   gofakeit.UUID(),
		},
	}

	approveResp, err := st.EventClient.ApproveEvent(ctx, approveReq)
	require.NoError(t, err, "ApproveEvent failed")
	require.Equal(t, domain.EventStatusApproved.String(), approveResp.GetStatus())
	require.Equal(t, approveReq.GetUser().GetId(), approveResp.GetApproveMetadata().GetApprovedBy().GetId())

	// Publish the created event
	publishReq := &eventv1.EventActionRequest{
		EventId: approveResp.GetId(),
		UserId:  gofakeit.Int64(),
	}

	_, err = st.EventClient.PublishEvent(ctx, publishReq)
	require.Error(t, err, "PublishEvent should fail")
}

func TestUnPublishEvent_UniversityScope_HappyPath(t *testing.T) {
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

	// Approve the created event
	approveReq := &eventv1.ApproveEventRequest{
		EventId: sendToReviewResp.GetId(),
		User: &eventv1.UserObject{
			Id:        gofakeit.Int64(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Barcode:   gofakeit.UUID(),
		},
	}

	approveResp, err := st.EventClient.ApproveEvent(ctx, approveReq)
	require.NoError(t, err, "ApproveEvent failed")
	require.Equal(t, domain.EventStatusApproved.String(), approveResp.GetStatus())
	require.Equal(t, approveReq.GetUser().GetId(), approveResp.GetApproveMetadata().GetApprovedBy().GetId())

	// Publish the created event
	publishReq := &eventv1.EventActionRequest{
		EventId: approveResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	publishResp, err := st.EventClient.PublishEvent(ctx, publishReq)
	require.NoError(t, err, "PublishEvent failed")

	// Assert that the event is published
	assert.Equal(t, domain.EventStatusInProgress.String(), publishResp.GetStatus())

	// Unpublish the created event
	unpublishReq := &eventv1.EventActionRequest{
		EventId: approveResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	unpublishResp, err := st.EventClient.UnpublishEvent(ctx, unpublishReq)
	require.NoError(t, err, "UnpublishEvent failed")

	// Assert that the event is unpublished
	assert.Equal(t, domain.EventStatusApproved.String(), unpublishResp.GetStatus())
}

func TestUnPublishEvent_ClubScope_HappyPath(t *testing.T) {
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
		Type:      domain.EventTypeIntraClub.String(),
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

	// Publish the created event
	publishReq := &eventv1.EventActionRequest{
		EventId: updateResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	publishResp, err := st.EventClient.PublishEvent(ctx, publishReq)
	require.NoError(t, err, "PublishEvent failed")

	// Assert that the event is published
	assert.Equal(t, domain.EventStatusInProgress.String(), publishResp.GetStatus())

	// Unpublish the created event
	unpublishReq := &eventv1.EventActionRequest{
		EventId: updateResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	unpublishResp, err := st.EventClient.UnpublishEvent(ctx, unpublishReq)
	require.NoError(t, err, "UnpublishEvent failed")

	// Assert that the event is unpublished
	assert.Equal(t, domain.EventStatusApproved.String(), unpublishResp.GetStatus())
}

func TestUnpublishEvent_ClubScope_Unauthorized(t *testing.T) {
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
		Type:      domain.EventTypeIntraClub.String(),
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

	// Publish the created event
	publishReq := &eventv1.EventActionRequest{
		EventId: updateResp.GetId(),
		UserId:  createResp.GetOwnerId(),
	}

	publishResp, err := st.EventClient.PublishEvent(ctx, publishReq)
	require.NoError(t, err, "PublishEvent failed")

	// Assert that the event is published
	assert.Equal(t, domain.EventStatusInProgress.String(), publishResp.GetStatus())

	// Unpublish the created event with unauthorized user
	unpublishReq := &eventv1.EventActionRequest{
		EventId: updateResp.GetId(),
		UserId:  gofakeit.Int64(), // Unauthorized user
	}

	_, err = st.EventClient.UnpublishEvent(ctx, unpublishReq)
	require.Error(t, err, "UnpublishEvent should fail")
}
