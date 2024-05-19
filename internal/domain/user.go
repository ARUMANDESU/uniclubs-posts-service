package domain

import eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"

type User struct {
	ID        int64  `json:"id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Barcode   string `json:"barcode" bson:"barcode"`
	AvatarURL string `json:"avatar_url" bson:"avatar_url"`
}

func (u User) ToProto() *eventv1.UserObject {
	return &eventv1.UserObject{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Barcode:   u.Barcode,
		AvatarUrl: u.AvatarURL,
	}
}

func UserFromProto(user *eventv1.UserObject) User {
	return User{
		ID:        user.GetId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Barcode:   user.GetBarcode(),
		AvatarURL: user.GetAvatarUrl(),
	}
}
