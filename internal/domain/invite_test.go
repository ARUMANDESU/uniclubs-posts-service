package domain

import "testing"

func TestUserInvite_IsByWho(t *testing.T) {
	tests := []struct {
		name   string
		u      UserInvite
		userId int64
		want   bool
	}{
		{
			name: "IsByWho",
			u: UserInvite{
				ByWhoId: 1,
			},
			userId: 1,
			want:   true,
		},
		{
			name: "IsNotByWho",
			u: UserInvite{
				ByWhoId: 1,
			},
			userId: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsByWho(tt.userId); got != tt.want {
				t.Errorf("UserInvite.IsByWho() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestUserInvite_IsInvited(t *testing.T) {
	tests := []struct {
		name   string
		u      UserInvite
		userId int64
		want   bool
	}{
		{
			name: "IsInvited",
			u: UserInvite{
				User: User{
					ID: 1,
				},
			},
			userId: 1,
			want:   true,
		},
		{
			name: "IsNotInvited",
			u: UserInvite{
				User: User{
					ID: 1,
				},
			},
			userId: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsInvited(tt.userId); got != tt.want {
				t.Errorf("UserInvite.IsInvited() = %v, want %v", got, tt.want)
			}
		})
	}

}
