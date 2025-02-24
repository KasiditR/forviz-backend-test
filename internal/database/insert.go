package database

import (
	"context"
	"github.com/KasiditR/forviz-backend-api-test/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"reflect"
	"time"
)

func InsertOne(modelPtr interface{}) (*mongo.InsertOneResult, error) {
	collection := Get().Database.Collection(utils.GetName(modelPtr))
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	val := reflect.ValueOf(modelPtr).Elem()
	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		newID := bson.NewObjectID()
		idField.Set(reflect.ValueOf(newID))
	}

	res, err := collection.InsertOne(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	if insertedID, ok := res.InsertedID.(bson.ObjectID); ok {
		idField.Set(reflect.ValueOf(insertedID))
	}

	return res, nil
}
