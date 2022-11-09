package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (r *MongoDB) DeleteMany(filter interface{}) (*mongo.DeleteResult, error) {
	delOptions := options.DeleteOptions{}
	db := r.Client.Database(r.Database).Collection(r.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.DeleteMany(ctx, filter, &delOptions)
}

func (r *MongoDB) DeleteOne(filter interface{}) (*mongo.DeleteResult, error) {
	delOptions := options.DeleteOptions{}
	db := r.Client.Database(r.Database).Collection(r.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.DeleteOne(ctx, filter, &delOptions)
}

func (r *MongoDB) Drop() error {
	db := r.Client.Database(r.Database).Collection(r.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.Drop(ctx)
}
