package mongodb

import (
	"context"
	"fmt"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"

	//	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CopyMsg struct {
	Result *mongo.InsertOneResult
	Err    error
}

func (req *MongoCopy) Copy(filter string) []*mongo.InsertOneResult {
	var (
		ch      = make(chan CopyMsg)
		fromDB  = req.Client.Database(req.FromDatabase).Collection(req.FromCollection)
		toDB    = req.Client.Database(req.ToDatabase).Collection(req.ToCollection)
		receive = 0
		results []*mongo.InsertOneResult
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var docs = req.Documents
	go func(docs []string) {
		for _, v := range docs {
			var (
				Filter bson.D
				msg    CopyMsg
			)
			if req.Num {
				if v1, err := strconv.Atoi(v); err == nil {
					Filter = bson.D{{filter, v1}}
				}
			} else {
				Filter = bson.D{{filter, v}}
			}
			if req.FromCollection == "block" {
				var data types.BlckN
				msg = CopyDocs(data, ctx, fromDB, toDB, Filter)
			}
			if req.FromCollection == "trans" {
				var data types.Trans
				msg = CopyDocs(data, ctx, fromDB, toDB, Filter)
			}
			ch <- msg
		}
	}(docs)
	for {
		select {
		case msg := <-ch:
			{
				receive++
				if msg.Err == nil {
					results = append(results, msg.Result)
				} else {
					log.Error().Msgf("Copy error %v", msg.Err)
				}
				if receive == len(docs) {
					return results
				}
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Printf(".")
		}
	}
}

func CopyDocs[T types.Trans | types.BlckN](data T, ctx context.Context, from *mongo.Collection, to *mongo.Collection, Filter bson.D) CopyMsg {
	var msg CopyMsg
	msg.Err = from.FindOne(ctx, Filter).Decode(&data)
	if msg.Err == nil {
		msg.Result, msg.Err = to.InsertOne(ctx, data)
	} else {
		msg.Result = nil
	}
	return msg
}
