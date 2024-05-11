package domain

type User struct {
	ID        int64  `json:"id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Barcode   string `json:"barcode" bson:"barcode"`
	AvatarURL string `json:"avatar_url" bson:"avatar_url"`
}
