package models

import "go.mongodb.org/mongo-driver/v2/bson"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	Id       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Email    string        `json:"email" bson:"email"`
	Password string        `json:"-" bson:"password"`
	Role     string        `json:"role" bson:"role"`
}
