package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"math"
	"time"
)

type SortOrder string
type SortBy string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"

	SortByDate         SortBy = "date"
	SortByParticipants SortBy = "participants"
	SortByType         SortBy = "type"
)

type Filters struct {
	Page      int32
	PageSize  int32
	Query     string
	SortBy    SortBy
	SortOrder SortOrder
	ClubId    int64
	UserId    int64
	Tags      []string
	FromDate  time.Time
	ToDate    time.Time
	Status    EventStatus
	Paths     []string
}

func (f Filters) Limit() int32 {
	return f.PageSize
}
func (f Filters) Offset() int32 {
	return (f.Page - 1) * f.PageSize
}

func ProtoToFilers(req *eventv1.ListEventsRequest) Filters {
	fromDate, _ := time.Parse(timeLayout, req.GetFilter().GetFromDate())
	tillDate, _ := time.Parse(timeLayout, req.GetFilter().GetTillDate())

	return Filters{
		Page:      req.GetPageNumber(),
		PageSize:  req.GetPageSize(),
		Query:     req.GetQuery(),
		SortBy:    SortBy(req.GetSortBy()),
		SortOrder: SortOrder(req.GetSortOrder()),
		ClubId:    req.GetFilter().GetClubId(),
		UserId:    req.GetFilter().GetUserId(),
		Tags:      req.GetFilter().GetTags(),
		FromDate:  fromDate,
		ToDate:    tillDate,
		Status:    EventStatus(req.GetFilter().GetStatus()),
		Paths:     req.GetFilterMask().GetPaths(),
	}
}

type PaginationMetadata struct {
	CurrentPage  int32
	PageSize     int32
	FirstPage    int32
	LastPage     int32
	TotalRecords int32
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
