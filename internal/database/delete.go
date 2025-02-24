package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func DeleteOne(filter bson.M, collectionName string) (result *mongo.DeleteResult, err error) {
	collection := Get().Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	result, err = collection.DeleteOne(ctx, collection)
	if err != nil {
		return nil, err
	}

	return result, nil
}
