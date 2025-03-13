package model

type Campaign struct {
	MongoId string `json:"_id,omitempty" bson:"_id,omitempty"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}
