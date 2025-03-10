package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"time"
)

type SortOrder string
type SortBy string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"

	SortByDate         SortBy = "date"
	SortByParticipants SortBy = "participants"
	SortByType         SortBy = "type"
)

func (s SortOrder) String() string {
	return string(s)
}

func (s SortBy) String() string {
	return string(s)
}

type BaseFilter struct {
	Page      int32
	PageSize  int32
	Query     string
	SortBy    SortBy
	SortOrder SortOrder
}

type EventsFilter struct {
	BaseFilter
	ClubId                int64
	UserId                int64
	Tags                  []string
	FromDate              time.Time
	ToDate                time.Time
	Status                []EventStatus
	IsHiddenForNonMembers bool
	Paths                 []string
}

func (f BaseFilter) Limit() int32 {
	return f.PageSize
}
func (f BaseFilter) Offset() int32 {
	return (f.Page - 1) * f.PageSize
}

func (f BaseFilter) Sort() string {
	return f.SortBy.String()
}

func ProtoToFilers(req *eventv1.ListEventsRequest) EventsFilter {
	fromDate, _ := time.Parse(TimeLayout, req.GetFilter().GetFromDate())
	tillDate, _ := time.Parse(TimeLayout, req.GetFilter().GetTillDate())

	var sortOrder SortOrder
	if req.GetSortOrder() == "" {
		sortOrder = SortOrderDesc
	} else {
		sortOrder = SortOrder(req.GetSortOrder())
	}
	filter := req.GetFilter()

	return EventsFilter{
		BaseFilter: BaseFilter{
			Page:      req.GetPageNumber(),
			PageSize:  req.GetPageSize(),
			Query:     req.GetQuery(),
			SortBy:    SortBy(req.GetSortBy()),
			SortOrder: sortOrder,
		},
		ClubId:                filter.GetClubId(),
		UserId:                filter.GetUserId(),
		Tags:                  filter.GetTags(),
		FromDate:              fromDate,
		ToDate:                tillDate,
		Status:                convertToEventStatusSlice(filter.GetStatus()),
		IsHiddenForNonMembers: filter.GetIsHiddenForNonMembers(),
		Paths:                 req.GetFilterMask().GetPaths(),
	}
}

func convertToEventStatusSlice(statuses []string) []EventStatus {
	var eventStatuses []EventStatus
	for _, status := range statuses {
		eventStatuses = append(eventStatuses, EventStatus(status))
	}
	return eventStatuses
}
