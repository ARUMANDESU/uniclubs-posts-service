package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
	"testing"
	"time"
)

func TestCalculateMetadata(t *testing.T) {
	tests := []struct {
		totalRecords   int32
		page           int32
		pageSize       int32
		expectedResult PaginationMetadata
	}{
		{
			totalRecords: 5,
			page:         1,
			pageSize:     10,
			expectedResult: PaginationMetadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 5,
			},
		},
		{
			totalRecords:   0,
			page:           1,
			pageSize:       10,
			expectedResult: PaginationMetadata{},
		},
		{
			totalRecords: 20,
			page:         2,
			pageSize:     10,
			expectedResult: PaginationMetadata{
				CurrentPage:  2,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     2,
				TotalRecords: 20,
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			metadata := CalculatePaginationMetadata(tt.totalRecords, tt.page, tt.pageSize)

			assert.Equal(t, tt.expectedResult, metadata)
		})
	}
}

func TestFilters_Limit(t *testing.T) {
	tests := []struct {
		filters       Filters
		expectedLimit int32
	}{
		{
			filters: Filters{
				Page:     1,
				PageSize: 10,
			},
			expectedLimit: 10,
		},
		{
			filters: Filters{
				Page:     1,
				PageSize: 1,
			},
			expectedLimit: 1,
		},
		{
			filters: Filters{
				Page:     1,
				PageSize: -1,
			},
			expectedLimit: -1,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.filters.Limit()

			assert.Equal(t, tt.expectedLimit, result)
		})
	}

}

func TestFilters_Offset(t *testing.T) {
	tests := []struct {
		filters        Filters
		expectedOffset int32
	}{
		{
			filters: Filters{
				Page:     1,
				PageSize: 10,
			},
			expectedOffset: 0,
		},
		{
			filters: Filters{
				Page:     2,
				PageSize: 20,
			},
			expectedOffset: 20,
		},
		{
			filters: Filters{
				Page:     5,
				PageSize: 30,
			},
			expectedOffset: 120,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := tt.filters.Offset()

			assert.Equal(t, tt.expectedOffset, result)
		})
	}
}

func TestProtoToFilers_ValidRequest(t *testing.T) {
	now := time.Now()

	req := &eventv1.ListEventsRequest{
		PageNumber: 1,
		PageSize:   10,
		Query:      "test",
		SortBy:     "date",
		SortOrder:  "asc",
		Filter: &eventv1.EventFilter{
			UserId:   1,
			ClubId:   1,
			Tags:     []string{"tag1", "tag2"},
			FromDate: now.Format(timeLayout),
			TillDate: now.Add(24 * time.Hour).Format(timeLayout),
			Status:   "PUBLISHED",
		},
		FilterMask: &field_mask.FieldMask{
			Paths: []string{"userId", "clubId", "tags", "fromDate", "tillDate", "status"},
		},
	}

	fromDate, _ := time.Parse(timeLayout, req.GetFilter().GetFromDate())
	tillDate, _ := time.Parse(timeLayout, req.GetFilter().GetTillDate())
	expectedFilters := Filters{
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

	filters := ProtoToFilers(req)

	assert.Equal(t, expectedFilters, filters)
}

func TestPaginationMetadata_ToProto(t *testing.T) {
	metadata := PaginationMetadata{
		CurrentPage:  1,
		PageSize:     10,
		FirstPage:    1,
		LastPage:     2,
		TotalRecords: 20,
	}

	expectedProto := &eventv1.PaginationMetadata{
		CurrentPage:  metadata.CurrentPage,
		PageSize:     metadata.PageSize,
		FirstPage:    metadata.FirstPage,
		LastPage:     metadata.LastPage,
		TotalRecords: metadata.TotalRecords,
	}

	proto := metadata.ToProto()

	assert.Equal(t, expectedProto, proto)
}
