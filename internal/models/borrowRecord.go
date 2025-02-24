package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BorrowRecord struct {
	ID       bson.ObjectID  `json:"id" bson:"_id,omitempty"`
	BookID   bson.ObjectID  `json:"book_id" bson:"book_id"`
	UserID   bson.ObjectID  `json:"user_id" bson:"user_id"`
	Borrowed *bson.DateTime `json:"borrowed_at" bson:"borrowed_at"`
	Returned *bson.DateTime `json:"returned_at,omitempty" bson:"returned_at,omitempty"`
}
