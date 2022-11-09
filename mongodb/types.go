package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Uri        string
	Client     *mongo.Client
	Option     *options.ClientOptions
	Database   string
	Collection string
	Document   string
}

type MongoCopy struct {
	Uri            string
	Client         *mongo.Client
	Option         *options.ClientOptions
	FromDatabase   string
	FromCollection string
	ToDatabase     string
	ToCollection   string
	Documents      []string
	Num            bool
}
