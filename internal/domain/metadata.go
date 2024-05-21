package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"math"
	"time"
)

type PaginationMetadata struct {
	CurrentPage  int32
	PageSize     int32
	FirstPage    int32
	LastPage     int32
	TotalRecords int32
}

type ApproveMetadata struct {
	ApprovedBy User
	ApprovedAt time.Time
}

type RejectMetadata struct {
	RejectedBy User
	RejectedAt time.Time
	Reason     string
}

func (m ApproveMetadata) ToProto() *eventv1.ApproveMetadata {
	return &eventv1.ApproveMetadata{
		ApprovedBy: m.ApprovedBy.ToProto(),
		ApprovedAt: m.ApprovedAt.Format(TimeLayout),
	}
}

func (m RejectMetadata) ToProto() *eventv1.RejectMetadata {
	return &eventv1.RejectMetadata{
		RejectedBy: m.RejectedBy.ToProto(),
		RejectedAt: m.RejectedAt.Format(TimeLayout),
		Reason:     m.Reason,
	}
}

func CalculatePaginationMetadata(totalRecords, page, pageSize int32) PaginationMetadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return PaginationMetadata{}
	}
	return PaginationMetadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int32(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func (m PaginationMetadata) ToProto() *eventv1.PaginationMetadata {
	return &eventv1.PaginationMetadata{
		CurrentPage:  m.CurrentPage,
		PageSize:     m.PageSize,
		FirstPage:    m.FirstPage,
		LastPage:     m.LastPage,
		TotalRecords: m.TotalRecords,
	}
}
