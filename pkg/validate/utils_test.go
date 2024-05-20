package validate

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_ValidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.UserObject
	}{
		{
			name: "Valid Request",
			req: &eventv1.UserObject{
				Id:        1,
				FirstName: "FirstName",
				LastName:  "LastName",
				Barcode:   "Barcode",
				AvatarUrl: "AvatarUrl",
			},
		},
		{
			name: "Also valid request",
			req: &eventv1.UserObject{
				Id:        1,
				FirstName: "FirstName",
				LastName:  "LastName",
				Barcode:   "Barcode",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user(tt.req)
			assert.NoError(t, err)
		})
	}
}

func TestUser_InvalidType(t *testing.T) {
	tests := []struct {
		name string
		req  any
	}{
		{
			name: "Invalid Type: String",
			req:  "Invalid Type",
		},
		{
			name: "Invalid Type: Int",
			req:  1,
		},
		{
			name: "Invalid Type: Float",
			req:  1.0,
		},
		{
			name: "Invalid Type: Bool",
			req:  true,
		},
		{
			name: "Invalid Type: Array",
			req:  []string{"Invalid Type"},
		},
		{
			name: "Invalid Type: Map",
			req:  map[string]string{"Invalid Type": "Invalid Type"},
		},
		{
			name: "Invalid Type: Struct",
			req: struct {
				InvalidType string
			}{
				InvalidType: "Invalid Type",
			},
		},
		{
			name: "Invalid Type: Nil",
			req:  nil,
		},
		{
			name: "Invalid Type: Empty",
		},
		{
			name: "Invalid Type: Function",
			req: func() string {
				return "Invalid Type"
			},
		},
		{
			name: "Invalid Type: Interface",
			req:  interface{}("Invalid Type"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user(tt.req)
			assert.Error(t, err)
		})

	}
}

func TestUser_Invalid(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.UserObject
	}{
		{
			name: "Empty FirstName",
			req: &eventv1.UserObject{
				Id:        1,
				FirstName: "",
				LastName:  "LastName",
				Barcode:   "Barcode",
			},
		},
		{
			name: "Empty LastName",
			req: &eventv1.UserObject{
				Id:        1,
				FirstName: "FirstName",
				LastName:  "",
				Barcode:   "Barcode",
			},
		},
		{
			name: "Empty Barcode",
			req: &eventv1.UserObject{
				Id:        1,
				FirstName: "FirstName",
				LastName:  "LastName",
				Barcode:   "",
			},
		},
		{
			name: "Invalid ID",
			req: &eventv1.UserObject{
				Id:        -1,
				FirstName: "FirstName",
				LastName:  "LastName",
				Barcode:   "Barcode",
			},
		},
		{
			name: "Empty ID",
			req: &eventv1.UserObject{
				FirstName: "FirstName",
				LastName:  "LastName",
				Barcode:   "Barcode",
			},
		},
		{
			name: "Empty request",
			req:  &eventv1.UserObject{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user(tt.req)
			assert.Error(t, err)
		})
	}
}

func TestClub_ValidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.ClubObject
	}{
		{
			name: "Valid Request without LogoUrl",
			req: &eventv1.ClubObject{
				Id:   1,
				Name: "Test Club",
			},
		},
		{
			name: "Valid request",
			req: &eventv1.ClubObject{
				Id:      1,
				Name:    "Test Club",
				LogoUrl: "https://test.com/logo.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := club(tt.req)
			assert.NoError(t, err)
		})
	}
}

func TestClub_InvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.ClubObject
	}{
		{
			name: "Empty Name",
			req: &eventv1.ClubObject{
				Id: 1,
			},
		},
		{
			name: "Invalid ID",
			req: &eventv1.ClubObject{
				Id:   -1,
				Name: "Test Club",
			},
		},
		{
			name: "Empty ID",
			req: &eventv1.ClubObject{
				Name: "Test Club",
			},
		},
		{
			name: "Empty Name and ID",
			req:  &eventv1.ClubObject{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := club(tt.req)
			assert.Error(t, err)
		})
	}
}

func TestAttachedFiles_ValidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  []*eventv1.FileObject
	}{
		{
			name: "With two files",
			req: []*eventv1.FileObject{
				{
					Url:  "https://test.com/file1.pdf",
					Name: "file1.pdf",
					Type: "pdf",
				},
				{
					Url:  "https://test.com/file2.pdf",
					Name: "file2.pdf",
					Type: "pdf",
				},
			},
		},
		{
			name: "With one file",
			req: []*eventv1.FileObject{
				{
					Url:  gofakeit.URL(),
					Name: "file1.pdf",
					Type: "pdf",
				},
			},
		},
		{
			name: "no files",
			req:  []*eventv1.FileObject{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := attachedFiles(tt.req)
			assert.NoError(t, err)
		})
	}

}

func TestAttachedFiles_InvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  any
	}{
		{
			name: "With invalid type",
			req: []*eventv1.CoverImage{
				{
					Url:  gofakeit.URL(),
					Name: gofakeit.ProductName(),
					Type: "pdf",
				},
				{
					Url:  gofakeit.URL(),
					Name: gofakeit.ProductName(),
					Type: "lol",
				},
			},
		},
		{
			name: "With invalid URL",
			req: []*eventv1.FileObject{
				{
					Url:  "invalid",
					Name: gofakeit.ProductName(),
					Type: "pdf",
				},
			},
		},
		{
			name: "Name is empty",
			req: []*eventv1.FileObject{
				{
					Url:  gofakeit.URL(),
					Name: "",
					Type: "pdf",
				},
			},
		},
		{
			name: "type is empty",
			req: []*eventv1.FileObject{
				{
					Url:  gofakeit.URL(),
					Name: gofakeit.ProductName(),
					Type: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := attachedFiles(tt.req)
			assert.Error(t, err)
		})
	}
}

func TestCoverImages_ValidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  []*eventv1.CoverImage
	}{
		{
			name: "With two images",
			req: []*eventv1.CoverImage{
				{
					Url:      gofakeit.URL(),
					Name:     gofakeit.ProductName(),
					Type:     "image/jpeg",
					Position: 1,
				},
				{
					Url:      gofakeit.URL(),
					Name:     gofakeit.ProductName(),
					Type:     "image/jpeg",
					Position: 2,
				},
			},
		},
		{
			name: "With one image",
			req: []*eventv1.CoverImage{
				{
					Url:      gofakeit.URL(),
					Name:     gofakeit.ProductName(),
					Type:     "image/jpeg",
					Position: 1,
				},
			},
		},
		{
			name: "no images",
			req:  []*eventv1.CoverImage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := coverImages(tt.req)
			assert.NoError(t, err)
		})
	}
}

func TestCoverImages_InvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  any
	}{
		{
			name: "With invalid type",
			req: []*eventv1.FileObject{
				{
					Url:  gofakeit.URL(),
					Name: gofakeit.ProductName(),
					Type: "pdf",
				},
				{
					Url:  gofakeit.URL(),
					Name: gofakeit.ProductName(),
					Type: "lol",
				},
			},
		},
		{
			name: "With invalid URL",
			req: []*eventv1.CoverImage{
				{
					Url:      "invalid",
					Name:     gofakeit.ProductName(),
					Type:     "image/jpeg",
					Position: 1,
				},
			},
		},
		{
			name: "Name is empty",
			req: []*eventv1.CoverImage{
				{
					Url:      gofakeit.URL(),
					Type:     "image/jpeg",
					Position: 1,
				},
			},
		},
		{
			name: "type is empty",
			req: []*eventv1.CoverImage{
				{
					Url:      gofakeit.URL(),
					Name:     gofakeit.ProductName(),
					Position: 1,
				},
			},
		},
		{
			name: "Position is empty",
			req: []*eventv1.CoverImage{
				{
					Url:  gofakeit.URL(),
					Name: gofakeit.ProductName(),
					Type: "image/jpeg",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := coverImages(tt.req)
			assert.Error(t, err)
		})
	}
}

func TestEventFilter_ValidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.EventFilter
	}{
		{
			name: "Empty Request",
			req:  &eventv1.EventFilter{},
		},
		{
			name: "Valid Request",
			req: &eventv1.EventFilter{
				UserId:   1,
				ClubId:   1,
				Tags:     []string{"tag1", "tag2"},
				FromDate: time.Now().Format(time.RFC3339),
				TillDate: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				Status:   domain.EventStatusPublished,
			},
		},
		{
			name: "Valid Request with empty tags",
			req: &eventv1.EventFilter{
				UserId:   1,
				ClubId:   1,
				FromDate: time.Now().Format(time.RFC3339),
				TillDate: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				Status:   domain.EventStatusPublished,
			},
		},
		{
			name: "Valid Request with empty dates",
			req: &eventv1.EventFilter{
				UserId: 1,
				ClubId: 1,
				Tags:   []string{"tag1", "tag2"},
				Status: domain.EventStatusPublished,
			},
		},
		{
			name: "Valid Request with empty status",
			req: &eventv1.EventFilter{
				UserId:   1,
				ClubId:   1,
				Tags:     []string{"tag1", "tag2"},
				FromDate: time.Now().Format(time.RFC3339),
				TillDate: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := eventFilter(tt.req)
			assert.NoError(t, err)
		})
	}
}

func TestEventFilter_InvalidRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *eventv1.EventFilter
	}{
		{
			name: "Invalid UserId",
			req: &eventv1.EventFilter{
				UserId: -1,
				ClubId: 1,
			},
		},
		{
			name: "Invalid ClubId",
			req: &eventv1.EventFilter{
				UserId: 1,
				ClubId: -1,
			},
		},
		{
			name: "Invalid Status",
			req: &eventv1.EventFilter{
				UserId: 1,
				ClubId: 1,
				Status: "INVALID_STATUS",
			},
		},
		{
			name: "Invalid FromDate",
			req: &eventv1.EventFilter{
				UserId:   1,
				ClubId:   1,
				FromDate: "INVALID_DATE",
			},
		},
		{
			name: "Invalid TillDate",
			req: &eventv1.EventFilter{
				UserId:   1,
				ClubId:   1,
				TillDate: "INVALID_DATE",
			},
		},
		{
			name: "Invalid Tags",
			req: &eventv1.EventFilter{
				UserId: 1,
				ClubId: 1,
				Tags: []string{
					gofakeit.Paragraph(1, 1, 76, " "),
					"valid tag",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := eventFilter(tt.req)
			assert.Error(t, err)
		})
	}
}

func newUser(t *testing.T) *eventv1.UserObject {
	t.Helper()
	return &eventv1.UserObject{
		Id:        1,
		FirstName: "FirstName",
		LastName:  "LastName",
		Barcode:   "Barcode",
		AvatarUrl: "AvatarUrl",
	}
}

func newClub(t *testing.T) *eventv1.ClubObject {
	t.Helper()
	return &eventv1.ClubObject{
		Id:      1,
		Name:    "Test Club",
		LogoUrl: "https://test.com/logo.jpg",
	}
}
