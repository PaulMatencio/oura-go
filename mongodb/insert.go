package mongodb

import (
	"context"
	"fmt"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

func (r *MongoDB) InsertOne(document interface{}) (result *mongo.InsertOneResult, err error) {

	db := r.Client.Database(r.Database).Collection(r.Collection)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result, err = db.InsertOne(ctx, document)
	return
}

func (r *MongoDB) InsertMany(docs []interface{}) (result *mongo.InsertManyResult, err error) {

	db := r.Client.Database(r.Database).Collection(r.Collection)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if result, err = db.InsertMany(ctx, docs); err != nil {
		result, err = r.CheckInsertMany(docs, err)
	}
	return
}

/*
	When bulk insert get a duplicate key error, it stops inserting the remaining documents in the list

	if insert many returned with duplicate key  then
        Check what missing
        insert missing key

*/

func (r *MongoDB) CheckInsertMany(docs []interface{}, err error) (results *mongo.InsertManyResult, err1 error) {

	if mongo.IsDuplicateKeyError(err) {
		log.Warn().Msgf("InsertMany %d docs into collection %s -  err %v - trying insertOne", len(docs), r.Collection, err)
		var (
			r1                    MongoDB
			opts                  *options.FindOneOptions
			result                types.Unknown
			results               mongo.InsertManyResult
			find, found, inserted int
		)
		r1.Uri = r.Uri
		r1.Option = r.Option
		r1.Database = r.Database
		r1.Collection = r.Collection
		if r1.Client, err1 = r1.Connect(); err1 == nil {
			defer r1.DisConnect()
			for _, doc := range docs {
				var fingerprint string
				v := reflect.ValueOf(doc)
				if v.Kind() == reflect.Ptr || v.Kind() == reflect.Struct {
					if map1, err := utils.ToMap(doc, "json"); err == nil {
						fingerprint = map1["fingerprint"].(string)
					}
				} else {
					if v.Kind() == reflect.Map {
						fingerprint = utils.GetFingerPrint(doc)
					}
				}
				if fingerprint != "" {
					// fingerprint = map1["fingerprint"].(string)
					filter := bson.M{"fingerprint": fingerprint}
					log.Trace().Msgf(fmt.Sprintf("find %v in collection %s", filter, r1.Collection))
					find++
					if err2 := r1.FindOne(opts, filter).Decode(&result); err2 != nil {
						if err2 == mongo.ErrNoDocuments {
							log.Trace().Msgf(fmt.Sprintf("insert one filter %v  into collection %s", filter, r1.Collection))
							if res2, err3 := r1.InsertOne(doc); err3 != nil {
								log.Error().Err(err).Msgf("insert fingerprint %v to collection %s", fingerprint, r1.Collection)
							} else {
								inserted++
								log.Info().Msgf(fmt.Sprintf("dcocument %v has been inserted into collection %s", res2.InsertedID, r1.Collection))
								results.InsertedIDs = append(results.InsertedIDs, res2.InsertedID.(primitive.ObjectID))
							}
						} else {
							log.Error().Err(err).Msgf("find One filter %v", filter)
						}
					} else {
						found++
						log.Trace().Msgf("found %v in collection %s", filter, r1.Collection)
						results.InsertedIDs = append(results.InsertedIDs, result.ID)
					}
				} else {
					log.Error().Err(err).Msgf("doc does not a fingerprint")
					err1 = err
				}
			}
			log.Info().Msgf("Number of doc#:%d - Find doc#:%d  Found doc#:%d - Inserted doc#:%d", len(docs), find, found, inserted)
			// r1.DisConnect()
		}

	} else {
		err1 = err
	}
	return
}
