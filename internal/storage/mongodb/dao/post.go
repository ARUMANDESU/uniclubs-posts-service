package dao

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID            primitive.ObjectID `bson:"_id"`
	Club          Club               `bson:"club"`
	Title         string             `bson:"title"`
	Description   string             `bson:"description"`
	Tags          []string           `bson:"tags"`
	CoverImages   []CoverImage       `bson:"cover_images"`
	AttachedFiles []File             `bson:"attached_files"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func PostFromDomain(p *domain.Post) *Post {
	objectID, _ := primitive.ObjectIDFromHex(p.ID)

	return &Post{
		ID:            objectID,
		Club:          ClubFromDomain(p.Club),
		Title:         p.Title,
		Description:   p.Description,
		Tags:          p.Tags,
		CoverImages:   ToCoverImages(p.CoverImages),
		AttachedFiles: ToFiles(p.AttachedFiles),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func PostToDomain(p *Post) *domain.Post {
	return &domain.Post{
		ID:            p.ID.Hex(),
		Club:          ToDomainClub(p.Club),
		Title:         p.Title,
		Description:   p.Description,
		Tags:          p.Tags,
		CoverImages:   ToDomainCoverImages(p.CoverImages),
		AttachedFiles: ToDomainFiles(p.AttachedFiles),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
