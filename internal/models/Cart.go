package models

import "go.mongodb.org/mongo-driver/v2/bson"

type CartItem struct {
	ProductId bson.ObjectID `json:"productId" bson:"productId"`
	Quantity  int           `json:"quantity" bson:"quantity"`
	Price     float64       `json:"price" bson:"price"`
}
type Cart struct {
	Id     bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId string        `json:"userId" bson:"userId"`
	Items  []CartItem    `json:"items" bson:"items"`
}
