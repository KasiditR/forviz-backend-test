package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Book struct {
	ID       bson.ObjectID `json:"book_id" bson:"_id,omitempty"`
	BookName *string       `json:"book_name" bson:"book_name" validate:"required"`
	Author   *string       `json:"author" bson:"author" validate:"required"`
	Category *string       `json:"category" bson:"category" validate:"required"`
}
