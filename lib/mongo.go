package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
MongoDB findOne
*/

func FindOne[T any](opts *options.FindOneOptions, v T, filter interface{}, req *mongodb.MongoDB) (result T, err error) {
	opts.SetAllowPartialResults(false)
	err = req.FindOne(opts, filter).Decode(&result)
	return result, err
}

func Find[T any](opts *options.FindOptions, v T, filter interface{}, req *mongodb.MongoDB) (cur *mongo.Cursor, result []T, err error) {
	cur, err = req.Find(opts, filter)
	if err == nil {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cur.All(ctx, &result)
	}
	return
}

func FindCtx[T any](opts *options.FindOptions, v T, filter interface{}, req *mongodb.MongoDB, ctx context.Context) (cur *mongo.Cursor, result []T, err error) {
	cur, err = req.Find(opts, filter)
	if err == nil {
		cur.All(ctx, &result)
	}
	return
}

/*
MongoDB find
*/

func FindAll[T any](opts *options.FindOptions, v T, filter interface{}, req *mongodb.MongoDB) (result []T, numErr int) {

	var (
		err error
		cur *mongo.Cursor
		db  *mongo.Collection
	)
	//ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db = req.Client.Database(req.Database).Collection(req.Collection)

	cur, err = db.Find(ctx, filter, opts)
	if cur != nil {
		for cur.Next(ctx) {
			var res T
			if err = cur.Decode(&res); err == nil {
				result = append(result, res)
			} else {
				numErr++
				LogError(err, fmt.Sprintf("find filter:%v - Decode", filter))

			}
		}
	} else {
		if err != mongo.ErrNoDocuments {
			numErr++
			LogError(err, fmt.Sprintf("FindAll filter %v", filter))
		}

	}
	return
}

/*
Init mon odb connection
*/

func InitReq(flags types.Options, mongoUrl string, filter string) (req *mongodb.MongoDB, Filter interface{}, err error) {
	var (
		filters types.Filters
	)

	req = &mongodb.MongoDB{
		Option:   flags.MongoOptions,
		Database: flags.MongoDatabase,
	}

	if filter != "" {
		opVal := filters.ValidOp()
		if filter != "" {
			if filters, err = filters.ParseOp(filter, opVal); err != nil {
				log.Error().Err(err).Msg("parsing filter")
				return req, Filter, err
			}
		}
		Filter = filters.BuildFilter()
	} else {
		Filter = bson.D{{}}
	}

	(*req).Uri = "mongodb://" + mongoUrl + "/?ssl=" + flags.MongoSSL
	log.Info().Msgf("request uri %s", (*req).Uri)
	if (*req).Client, err = (*req).Connect(); err != nil {
		log.Error().Err(err).Msg("mongodb connect")
	}
	return req, Filter, err
}

func List[T any](opts *options.FindOptions, result []T, filter interface{}, req *mongodb.MongoDB) {

	cur, err := req.Find(opts, filter)
	if err == nil {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cur.All(ctx, &result)
		for _, v := range result {
			var b []byte
			if req.Collection == "pool" || req.Collection == "trans" || req.Collection == "pooln" {
				b, err = json.Marshal(v)
			} else {
				b, err = bson.MarshalExtJSON(v, false, true)
			}
			if err == nil {
				j, _ := utils.PrettyJson(string(b))
				fmt.Println(j)
			} else {
				LogError(err, "Marshal Extended Json")
				// log.Error().Err(err).Msg("Marshal Extended Json")
			}
		}
		fmt.Printf("Total number  documents: %d\n", len(result))
	} else {
		log.Error().Err(err).Msg("List")
	}
}
