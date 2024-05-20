package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile_ToProto(t *testing.T) {
	file := File{
		Name: "testFile",
		Url:  "http://example.com/testFile",
		Type: "jpg",
	}

	protoFile := file.ToProto()

	assert.Equal(t, file.Name, protoFile.GetName())
	assert.Equal(t, file.Url, protoFile.GetUrl())
	assert.Equal(t, file.Type, protoFile.GetType())
}

func TestCoverImage_ToProto(t *testing.T) {
	coverImage := CoverImage{
		File: File{
			Name: "testCoverImage",
			Url:  "http://example.com/testCoverImage",
			Type: "jpg",
		},
		Position: 1,
	}

	protoCoverImage := coverImage.ToProto()

	assert.Equal(t, coverImage.Name, protoCoverImage.GetName())
	assert.Equal(t, coverImage.Url, protoCoverImage.GetUrl())
	assert.Equal(t, coverImage.Type, protoCoverImage.GetType())
	assert.Equal(t, int32(coverImage.Position), protoCoverImage.GetPosition())
}

func TestProtoToFile(t *testing.T) {
	protoFile := &eventv1.FileObject{
		Name: "testProtoFile",
		Url:  "http://example.com/testProtoFile",
		Type: "jpg",
	}

	file := ProtoToFile(protoFile)

	assert.Equal(t, protoFile.GetName(), file.Name)
	assert.Equal(t, protoFile.GetUrl(), file.Url)
	assert.Equal(t, protoFile.GetType(), file.Type)
}

func TestProtoToCoverImage(t *testing.T) {
	protoCoverImage := &eventv1.CoverImage{
		Name:     "testProtoCoverImage",
		Url:      "http://example.com/testProtoCoverImage",
		Type:     "jpg",
		Position: 1,
	}

	coverImage := ProtoToCoverImage(protoCoverImage)

	assert.Equal(t, protoCoverImage.GetName(), coverImage.Name)
	assert.Equal(t, protoCoverImage.GetUrl(), coverImage.Url)
	assert.Equal(t, protoCoverImage.GetType(), coverImage.Type)
	assert.Equal(t, uint32(protoCoverImage.GetPosition()), coverImage.Position)
}
