package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Product struct {
	Id          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name" validate:"required"`
	Cuisines    []string      `json:"cuisines,omitempty" required:"omitempty"`
	Category    string        `json:"category" validate:"required"`
	Price       float64       `json:"price" validate:"required,gte=0"`
	Description string        `json:"description,omitempty" bson:"omitempty"`
	Rating      float32       `json:"rating" validate:"gte=0"`
}
