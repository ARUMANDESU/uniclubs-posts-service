package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"time"
)

type BanRecord struct {
	EventId  string    `json:"event_id"`
	User     User      `json:"user"`
	BannedAt time.Time `json:"banned_at"`
	Reason   string    `json:"reason"`
	ByWhoId  int64     `json:"by_who_id"`
}

func ToProtoBanRecord(banRecord BanRecord) *eventv1.BanRecord {
	return &eventv1.BanRecord{
		User:     banRecord.User.ToProto(),
		Reason:   banRecord.Reason,
		BannedAt: banRecord.BannedAt.String(),
		BannedBy: banRecord.ByWhoId,
	}
}

func ToProtoBanRecords(banRecords []BanRecord) []*eventv1.BanRecord {
	result := make([]*eventv1.BanRecord, 0, len(banRecords))
	for _, banRecord := range banRecords {
		result = append(result, ToProtoBanRecord(banRecord))
	}
	return result
}
