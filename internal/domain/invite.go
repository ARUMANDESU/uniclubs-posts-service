package domain

import eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"

type Invite struct {
	ID    string `json:"id"`
	Event Event  `json:"event"`
	Club  Club   `json:"club"`
}

type UserInvite struct {
	ID      string `json:"id"`
	Event   Event  `json:"event"`
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

func ClubInviteToProto(invite Invite) *eventv1.ClubInvite {
	return &eventv1.ClubInvite{
		Id:    invite.ID,
		Event: invite.Event.ToProto(),
		Club:  invite.Club.ToProto(),
	}
}

func ClubInvitesToProto(invites []Invite) []*eventv1.ClubInvite {
	convertedInvites := make([]*eventv1.ClubInvite, len(invites))
	for i, invite := range invites {
		convertedInvites[i] = ClubInviteToProto(invite)
	}
	return convertedInvites
}

func UserInviteToProto(invite UserInvite) *eventv1.OrganizerInvite {
	return &eventv1.OrganizerInvite{
		Id:      invite.ID,
		Event:   invite.Event.ToProto(),
		ClubId:  invite.ClubId,
		ByWhoId: invite.ByWhoId,
		User:    invite.User.ToProto(),
	}
}

func UserInvitesToProto(invites []UserInvite) []*eventv1.OrganizerInvite {
	convertedInvites := make([]*eventv1.OrganizerInvite, len(invites))
	for i, invite := range invites {
		convertedInvites[i] = UserInviteToProto(invite)
	}
	return convertedInvites
}
