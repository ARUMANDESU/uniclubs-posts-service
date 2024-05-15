package mongodb

type Club struct {
	ID      int64  `json:"id,omitempty" bson:"_id"`
	Name    string `json:"name,omitempty" bson:"name"`
	LogoURL string `json:"logo_url,omitempty" bson:"logo_url"`
}
