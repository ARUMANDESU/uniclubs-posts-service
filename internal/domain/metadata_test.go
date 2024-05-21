package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
