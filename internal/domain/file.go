package domain

import posts "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts"

type File struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
	Type string `json:"type" bson:"type"`
}

type CoverImage struct {
	File
	Position uint32 `json:"position" bson:"position"`
}

func (f File) ToPb() *posts.FileObject {
	return &posts.FileObject{
		Name: f.Name,
		Url:  f.Url,
		Type: f.Type,
	}
}

func (c CoverImage) ToPb() *posts.CoverImage {
	return &posts.CoverImage{
		Name:     c.Name,
		Url:      c.Url,
		Type:     c.Type,
		Position: int32(c.Position),
	}
}

func CoverImagesToPb(images []CoverImage) []*posts.CoverImage {
	convertedImages := make([]*posts.CoverImage, len(images))
	for i, image := range images {
		convertedImages[i] = image.ToPb()
	}
	return convertedImages
}

func FilesToPb(files []File) []*posts.FileObject {
	convertedFiles := make([]*posts.FileObject, len(files))
	for i, file := range files {
		convertedFiles[i] = file.ToPb()
	}
	return convertedFiles
}

func PbToFile(file *posts.FileObject) File {
	return File{
		Name: file.Name,
		Url:  file.Url,
		Type: file.Type,
	}
}

func PbToCoverImage(image *posts.CoverImage) CoverImage {
	return CoverImage{
		File: File{
			Name: image.Name,
			Url:  image.Url,
			Type: image.Type,
		},
		Position: uint32(image.Position),
	}
}

func PbToCoverImages(images []*posts.CoverImage) []CoverImage {
	convertedImages := make([]CoverImage, len(images))
	for i, image := range images {
		convertedImages[i] = PbToCoverImage(image)
	}
	return convertedImages
}

func PbToFiles(files []*posts.FileObject) []File {
	convertedFiles := make([]File, len(files))
	for i, file := range files {
		convertedFiles[i] = PbToFile(file)
	}
	return convertedFiles
}
