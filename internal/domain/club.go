package domain

type Club struct {
	ID      int64  `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	LogoURL string `json:"logo_url" bson:"logo_url"`
}
