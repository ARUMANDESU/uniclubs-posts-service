package domain

import eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"

type File struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
	Type string `json:"type" bson:"type"`
}

type CoverImage struct {
	File
	Position int32 `json:"position" bson:"position"`
}

func (f File) ToProto() *eventv1.FileObject {
	return &eventv1.FileObject{
		Name: f.Name,
		Url:  f.Url,
		Type: f.Type,
	}
}

func (c CoverImage) ToProto() *eventv1.CoverImage {
	return &eventv1.CoverImage{
		Name:     c.Name,
		Url:      c.Url,
		Type:     c.Type,
		Position: c.Position,
	}
}

func CoverImagesToProto(images []CoverImage) []*eventv1.CoverImage {
	convertedImages := make([]*eventv1.CoverImage, len(images))
	for _, image := range images {
		convertedImages = append(convertedImages, image.ToProto())
	}
	return convertedImages
}

func FilesToProto(files []File) []*eventv1.FileObject {
	convertedFiles := make([]*eventv1.FileObject, len(files))
	for _, file := range files {
		convertedFiles = append(convertedFiles, file.ToProto())
	}
	return convertedFiles
}

func ProtoToFile(file *eventv1.FileObject) *File {
	return &File{
		Name: file.Name,
		Url:  file.Url,
		Type: file.Type,
	}
}

func ProtoToCoverImage(image *eventv1.CoverImage) *CoverImage {
	return &CoverImage{
		File: File{
			Name: image.Name,
			Url:  image.Url,
			Type: image.Type,
		},
		Position: image.Position,
	}
}

func ProtoToCoverImages(images []*eventv1.CoverImage) []*CoverImage {
	convertedImages := make([]*CoverImage, len(images))
	for _, image := range images {
		convertedImages = append(convertedImages, ProtoToCoverImage(image))
	}
	return convertedImages
}

func ProtoToFiles(files []*eventv1.FileObject) []*File {
	convertedFiles := make([]*File, len(files))
	for _, file := range files {
		convertedFiles = append(convertedFiles, ProtoToFile(file))
	}
	return convertedFiles
}
