package types

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Options struct {
	Concurrent      bool
	Limit           int64
	Upload          bool
	RetryRead       bool
	RetryWrite      bool
	MaxRetry        int
	PrintIt         bool
	CheckDup        bool
	Reg             *bsoncodec.Registry
	MongoOptions    *options.ClientOptions
	MongoDatabase   string
	MongoCollection string
	MongoSSL        string
	MarketPlace     MktPlace
}

func (o *Options) SetMarketPlace(mkpl MktPlace) {
	o.MarketPlace = mkpl
}
