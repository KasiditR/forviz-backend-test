package database

import (
	"context"
	"log"
	"time"

	"github.com/KasiditR/forviz-backend-api-test/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	client   *mongo.Client
	Database *mongo.Database
	Err      error
}

var (
	_mongo         Mongo
	ShortWaitTime  time.Duration = 2
	MediumWaitTime time.Duration = 5
	LongWaitTime   time.Duration = 10
)

func ConnectDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), LongWaitTime*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(config.LoadConfig().MongoURI)
	_mongo.client, _mongo.Err = mongo.Connect(clientOptions)
	if _mongo.Err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", _mongo.Err)
	}

	if err := _mongo.client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("MongoDB is not reachable: %v", err)
	}

	_mongo.Database = _mongo.client.Database(config.LoadConfig().MongoDatabase)

	migrateCollections("books")
	migrateCollections("borrow_records")
	migrateCollections("users")

	log.Println("Connected to MongoDB ")
}

func Get() Mongo {
	if _mongo.client == nil {
		ConnectDatabase()
	}

	return _mongo
}

func migrateCollections(collectionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := _mongo.Database.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		log.Fatal("Error listing collections:", err)
	}

	if len(collections) == 0 {
		err = _mongo.Database.CreateCollection(ctx, collectionName)
		if err != nil {
			log.Fatal("Failed to create collection:", err)
		} else {
			log.Println("Collection created successfully")
		}
	}
}
