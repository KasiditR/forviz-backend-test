package database

import (
	"context"
	"time"

	"github.com/KasiditR/forviz-backend-api-test/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CountDocument(filter bson.M, b interface{}) (c int64, err error) {
	collection := Get().Database.Collection(utils.GetName(b))
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, err
}

func FindOne(filter bson.M, b interface{}) (err error) {
	collection := Get().Database.Collection(utils.GetName(b))
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	res := collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return res.Err()
	}

	err = res.Decode(b)
	if err != nil {
		return err
	}

	return nil
}

func FindByID(id string, b interface{}) (err error) {
	collection := Get().Database.Collection(utils.GetName(b))
	ctx, cancel := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)
	defer cancel()

	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res := collection.FindOne(ctx, bson.M{"_id": userID})
	if res.Err() != nil {
		return res.Err()
	}
	err = res.Decode(b)
	if err != nil {
		return err
	}

	return nil
}

func FindByObjectID(objectID bson.ObjectID, bPtr interface{}) (err error) {
	collection := Get().Database.Collection(utils.GetName(bPtr))
	ctx, cancel := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)
	defer cancel()

	res := collection.FindOne(ctx, bson.M{"_id": objectID})
	if res.Err() != nil {
		return res.Err()
	}

	err = res.Decode(bPtr)
	if err != nil {
		return err
	}

	return nil
}

func FindAll(filter bson.M, modelsOutArrayPtr interface{}) error {
	collection := Get().Database.Collection(utils.GetName(modelsOutArrayPtr))
	ctx, cancel := context.WithTimeout(context.Background(), LongWaitTime*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	err = cur.All(ctx, modelsOutArrayPtr)
	if err != nil {
		return err
	}
	return nil
}

func FindAllWithPagination(filter bson.M, start int64, count int64, modelsOutArrayPtr interface{}) error {
	collection := Get().Database.Collection(utils.GetName(modelsOutArrayPtr))
	ctx, cancel := context.WithTimeout(context.Background(), LongWaitTime*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSkip(start)
	findOptions.SetLimit(count)

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return err
	}
	err = cur.All(ctx, modelsOutArrayPtr)
	if err != nil {
		return err
	}
	return nil
}

func FindOneAndUpdate(collectionName string, filter bson.M, update bson.M) (err error) {
	collection := Get().Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	cursor := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if cursor.Err() != nil {
		return cursor.Err()
	}

	return nil
}

func FindByIDAndUpdate(collectionName string, objectID bson.ObjectID, update bson.M) (err error) {
	collection := Get().Database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	cursor := collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts)
	if cursor.Err() != nil {
		return cursor.Err()
	}

	return nil
}
