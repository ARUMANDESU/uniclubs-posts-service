package dao

import "github.com/arumandesu/uniclubs-posts-service/internal/domain"

type FileMongo struct {
	URL  string `bson:"url"`
	Name string `bson:"name"`
	Type string `bson:"type"`
}

type CoverImageMongo struct {
	FileMongo
	Position uint32 `bson:"position"`
}

// Into dao

func ToCoverImages(coverImages []domain.CoverImage) []CoverImageMongo {
	coverImagesMongo := make([]CoverImageMongo, len(coverImages))
	for i, coverImage := range coverImages {
		coverImagesMongo[i] = CoverImageMongo{
			FileMongo: FileMongo{
				URL:  coverImage.Url,
				Name: coverImage.Name,
				Type: coverImage.Type,
			},
			Position: coverImage.Position,
		}
	}
	return coverImagesMongo
}

func ToFiles(files []domain.File) []FileMongo {
	filesMongo := make([]FileMongo, len(files))
	for i, file := range files {
		filesMongo[i] = FileMongo{
			URL:  file.Url,
			Name: file.Name,
			Type: file.Type,
		}
	}
	return filesMongo
}

// From dao to domain

func ToDomainFiles(filesMongo []FileMongo) []domain.File {
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

func ToDomainCoverImages(coverImagesMongo []CoverImageMongo) []domain.CoverImage {
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
