package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
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
	fromDate, _ := time.Parse(TimeLayout, req.GetFilter().GetFromDate())
	tillDate, _ := time.Parse(TimeLayout, req.GetFilter().GetTillDate())

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
