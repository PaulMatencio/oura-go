package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IMongoDB interface {
	Connect() (*mongo.Client, error)
	DisConnect() error
}

func (r *MongoDB) Connect() (*mongo.Client, error) {

	clientOptions := r.Option.ApplyURI(r.Uri).SetMaxPoolSize(100)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, clientOptions)
}

func (r *MongoCopy) Connect() (*mongo.Client, error) {
	clientOptions := r.Option.ApplyURI(r.Uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, clientOptions)
}
