package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID       bson.ObjectID `json:"user_id" bson:"_id,omitempty"`
	Username *string       `json:"username" bson:"username" validate:"required"`
	Password *string       `json:"password" bson:"password" validate:"required"`
}
