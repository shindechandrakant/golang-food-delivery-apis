package config

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	mongoClient   *mongo.Client
	mongoDbClient *mongo.Database
	once          sync.Once
)

func LoadMongoConnection() *mongo.Database {

	once.Do(func() {
		dbUrl := GetEnv("DB_URI")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		client, err := mongo.Connect(options.Client().ApplyURI(dbUrl))
		if err != nil {
			panic(err)
		}

		_ = client.Ping(ctx, readpref.Primary())
		databaseName := GetEnv("DB_NAME")
		mongoDbClient = client.Database(databaseName)
		log.Printf("DB Connected Successfully")

	})

	return mongoDbClient
}

func CloseMongoConnection() {

	if mongoClient == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Println("Mongo disconnect error:", err)
		return
	}

	log.Println("MongoDB disconnected")
}
