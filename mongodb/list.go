package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *MongoDB) List(options *options.FindOptions, filter interface{}) (cur *mongo.Cursor, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := r.Client.Database(r.Database).Collection(r.Collection)
	return db.Find(ctx, filter, options)
}

func (r *MongoDB) List1(limit int64, filter interface{}, reg *bsoncodec.Registry) (cur *mongo.Cursor, err error) {

	var (
		db          *mongo.Collection
		findOptions = options.Find().SetLimit(limit)
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if r.Collection == "pool" || r.Collection == "trans" || r.Collection == "pooln" {
		db = r.Client.Database(r.Database).Collection(r.Collection, &options.CollectionOptions{
			Registry: reg,
		})
	} else {
		db = r.Client.Database(r.Database).Collection(r.Collection)
	}
	// return db.Find(ctx, bson.D{{}}, findOptions)
	return db.Find(ctx, filter, findOptions)
}
