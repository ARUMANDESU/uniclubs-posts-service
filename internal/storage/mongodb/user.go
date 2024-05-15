package mongodb

type User struct {
	ID        int64  `json:"id,omitempty" bson:"_id"`
	FirstName string `json:"first_name,omitempty" bson:"first_name"`
	LastName  string `json:"last_name,omitempty" bson:"last_name"`
	Barcode   string `json:"barcode,omitempty" bson:"barcode"`
	AvatarURL string `json:"avatar_url,omitempty" bson:"avatar_url"`
}

type Organizer struct {
	User
	ClubId int64 `json:"club_id" bson:"club_id"`
}
