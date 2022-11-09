package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *MongoDB) FindOne(opts *options.FindOneOptions, filter interface{}) (result *mongo.SingleResult) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := r.Client.Database(r.Database).Collection(r.Collection)
	result = db.FindOne(ctx, filter, opts)
	return
}

func (r *MongoDB) Find(opts *options.FindOptions, filter interface{}) (result *mongo.Cursor, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := r.Client.Database(r.Database).Collection(r.Collection)
	result, err = db.Find(ctx, filter, opts)
	return
}
