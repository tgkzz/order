package mongo

import "go.mongodb.org/mongo-driver/v2/mongo"

type Mongo struct {
	db *mongo.Database
}
