package data

import "go.mongodb.org/mongo-driver/mongo"

// DB holds the mongo db driver collection
type DB struct {
	collection *mongo.Collection
}
