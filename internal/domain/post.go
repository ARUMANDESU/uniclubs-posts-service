package domain

import (
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Post struct {
	ID            string
	Club          Club
	Title         string
	Description   string
	Tags          []string
	CoverImages   []CoverImage
	AttachedFiles []File
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func PostToPb(post *Post) *postv1.PostObject {
	if post == nil {
		return nil
	}

	return &postv1.PostObject{
		Id:            post.ID,
		Club:          post.Club.ToProto(),
		Title:         post.Title,
		Description:   post.Description,
		Tags:          post.Tags,
		CoverImages:   CoverImagesToPb(post.CoverImages),
		AttachedFiles: FilesToPb(post.AttachedFiles),
		CreatedAt:     timestamppb.New(post.CreatedAt),
		UpdatedAt:     timestamppb.New(post.UpdatedAt),
	}
}

func PostsToPb(posts []Post) []*postv1.PostObject {
	if posts == nil {
		return nil
	}

	pbPosts := make([]*postv1.PostObject, 0, len(posts))
	for _, post := range posts {
		pbPosts = append(pbPosts, PostToPb(&post))
	}

	return pbPosts
}
