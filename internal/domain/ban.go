package domain

import "time"

type BanRecord struct {
	EventId  string    `json:"event_id"`
	UserId   int64     `json:"user_id"`
	BannedAt time.Time `json:"banned_at"`
	Reason   string    `json:"reason"`
	ByWhoId  int64     `json:"by_who_id"`
}
