package model

type Candidate struct {
	MongoId string `json:"_id,omitempty" bson:"_id,omitempty"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Teaser  string `json:"teaser"`
	Info    string `json:"info"`
	Image   string `json:"image"`
}
