package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	GetCollection(name string) *mongo.Collection
}

type database struct {
	client *mongo.Client
}

func MongoDB() Database {
	client := connect(getEnv("DB_CONNECTION"))
	return &database{client}
}

// MONGODB
func (d *database) GetCollection(name string) *mongo.Collection {
	return d.client.Database(getEnv("DB_NAME")).Collection(name)
}

func connect(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo client read error: ", err.Error())
	}

	fmt.Println(client.Database(getEnv("DB_NAME")).ListCollectionNames(
		ctx,
		bson.M{},
	))
	return client
}

func getEnv(value string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(value)
}
