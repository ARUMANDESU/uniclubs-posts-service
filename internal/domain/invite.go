package domain

type Invite struct {
	ID      string `json:"id"`
	EventId string `json:"event_id"`
	ClubId  int64  `json:"club_id"`
}

type UserInvite struct {
	Invite
	UserId int64 `json:"user_id"`
}

func (u UserInvite) IsInvited(userId int64) bool {
	return u.UserId == userId
}
