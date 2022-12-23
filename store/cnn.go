package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	localColl *mongo.Collection
}

func Connect() *Store {
	clientOptions := options.Client().ApplyURI("mongo")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("data")

	return &Store{
		localColl: db.Collection("person"),
	}
}