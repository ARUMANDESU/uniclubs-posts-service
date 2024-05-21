package dao

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"time"
)

type ApproveMetadata struct {
	ApprovedBy User      `bson:"user"`
	ApprovedAt time.Time `bson:"approved_at"`
}

type RejectMetadata struct {
	RejectedBy User      `bson:"user"`
	RejectedAt time.Time `bson:"rejected_at"`
	Reason     string    `bson:"reason,omitempty"`
}

func (m ApproveMetadata) ToDomain() domain.ApproveMetadata {
	return domain.ApproveMetadata{
		ApprovedBy: ToDomainUser(m.ApprovedBy),
		ApprovedAt: m.ApprovedAt,
	}
}

func (m RejectMetadata) ToDomain() domain.RejectMetadata {
	return domain.RejectMetadata{
		RejectedBy: ToDomainUser(m.RejectedBy),
		RejectedAt: m.RejectedAt,
		Reason:     m.Reason,
	}
}

func ToApproveMetadata(m domain.ApproveMetadata) ApproveMetadata {
	return ApproveMetadata{
		ApprovedBy: UserFromDomainUser(m.ApprovedBy),
		ApprovedAt: m.ApprovedAt,
	}
}

func ToRejectMetadata(m domain.RejectMetadata) RejectMetadata {
	return RejectMetadata{
		RejectedBy: UserFromDomainUser(m.RejectedBy),
		RejectedAt: m.RejectedAt,
		Reason:     m.Reason,
	}
}
