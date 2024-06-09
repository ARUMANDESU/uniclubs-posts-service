package dao

import "github.com/arumandesu/uniclubs-posts-service/internal/domain"

type File struct {
	URL  string `bson:"url"`
	Name string `bson:"name"`
	Type string `bson:"type"`
}

type CoverImage struct {
	File
	Position uint32 `bson:"position"`
}

// Into dao

func ToCoverImages(coverImages []domain.CoverImage) []CoverImage {
	coverImagesMongo := make([]CoverImage, len(coverImages))
	for i, coverImage := range coverImages {
		coverImagesMongo[i] = CoverImage{
			File: File{
				URL:  coverImage.Url,
				Name: coverImage.Name,
				Type: coverImage.Type,
			},
			Position: coverImage.Position,
		}
	}
	return coverImagesMongo
}

func ToFiles(files []domain.File) []File {
	filesMongo := make([]File, len(files))
	for i, file := range files {
		filesMongo[i] = File{
			URL:  file.Url,
			Name: file.Name,
			Type: file.Type,
		}
	}
	return filesMongo
}

// From dao to domain

func ToDomainFiles(filesMongo []File) []domain.File {
	files := make([]domain.File, len(filesMongo))
	for i, file := range filesMongo {
		files[i] = domain.File{
			Url:  file.URL,
			Name: file.Name,
			Type: file.Type,
		}
	}
	return files
}

func ToDomainCoverImages(coverImagesMongo []CoverImage) []domain.CoverImage {
	coverImages := make([]domain.CoverImage, len(coverImagesMongo))
	for i, coverImage := range coverImagesMongo {
		coverImages[i] = domain.CoverImage{
			File: domain.File{
				Url:  coverImage.URL,
				Name: coverImage.Name,
				Type: coverImage.Type,
			},
			Position: coverImage.Position,
		}
	}
	return coverImages
}
