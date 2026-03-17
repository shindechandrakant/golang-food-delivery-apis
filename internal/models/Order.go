package models

import "go.mongodb.org/mongo-driver/v2/bson"

type OrderItem struct {
	ProductId string `json:"productId" bson:"productId"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}

type Order struct {
	Id    bson.ObjectID `json:"-" bson:"_id,omitempty"`
	UUID  string        `json:"id" bson:"uuid"`
	Items []OrderItem   `json:"items" bson:"items"`
}
