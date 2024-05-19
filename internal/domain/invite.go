package domain

type Invite struct {
	ID      string `json:"id"`
	EventId string `json:"event_id"`
	ClubId  int64  `json:"club_id"`
}

type UserInvite struct {
	ID      string `json:"id"`
	EventId string `json:"event_id"`
	ClubId  int64  `json:"club_id"`
	ByWhoId int64  `json:"by_who_id"`
	User    User   `json:"user"`
}

func (u UserInvite) IsInvited(userId int64) bool {
	return u.User.ID == userId
}

func (u UserInvite) IsByWho(userId int64) bool {
	return u.ByWhoId == userId
}
