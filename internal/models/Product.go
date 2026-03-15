package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Image struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}
type Product struct {
	Id          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name" validate:"required"`
	Cuisines    []string      `json:"cuisines,omitempty" required:"omitempty"`
	Category    string        `json:"category" validate:"required"`
	Price       float64       `json:"price" validate:"required,gte=0"`
	Description string        `json:"description,omitempty" bson:"omitempty"`
	Rating      float32       `json:"rating" validate:"gte=0"`
	Image       []Image       `json:"image"`
}
