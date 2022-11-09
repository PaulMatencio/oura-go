package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/paulmatencio/oura-go/db"
	lib "github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "strconv"
	"time"
)

var (
	buildTranCmd = &cobra.Command{
		Use:   "buildTran",
		Short: "Build a given transaction  form a given source",
		Long:  ``,
		Run:   BuildTran,
	}
	buildTransCmd = &cobra.Command{
		Use:   "buildTrans",
		Short: "Build all transactions ( filter) from mongodb collections",
		Long:  ``,
		Run:   BuildTrans,
	}
)

func init() {
	rootCmd.AddCommand(buildTranCmd)
	rootCmd.AddCommand(buildTransCmd)
	initBuildTran(buildTranCmd)
	initBuildTrans(buildTransCmd)
	if database == "" {
		database = mongoDatabase
	}
	if database == "" {
		log.Fatal().Msg("mongodb database is missing")
	}

}

func initBuildTran(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().BoolVarP(&printIt, "print", "p", true, "print the result")
}

func initBuildTrans(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "", "local database directory. If no a full path, it with be prefix by home directory")
	cmd.Flags().StringVarP(&fromId, "from-id", "", "", "from tx object id")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().Int64VarP(&limit, "limit", "", 10, "max returned")
	cmd.Flags().BoolVarP(&printIt, "print", "p", true, "print the result")
	cmd.Flags().BoolVarP(&concurrent, "concurrent", "C", false, "concurrent upload")
	cmd.Flags().BoolVarP(&upload, "upload", "U", false, "upload to mongodb ")
	cmd.Flags().IntVarP(&bulk, "bulk", "n", 10, "max bulkload")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger Filename ")
}

func BuildTran(cmd *cobra.Command, args []string) {

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Error().Err(err).Msgf("logging to file %s", logger, err)
	}
	flags := SetFlags()
	if req, Filter, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
		var (
			start = time.Now()
			trans types.Trans
		)
		defer req.DisConnect()
		if trans, err = lib.BuildTran(flags, req, Filter); err == nil {
			if b, err := json.Marshal(&trans); err == nil {
				if !printIt {
					utils.PrettyJson(string(b))
				} else {
					utils.PrintJson(string(b))
				}
			} else {
				log.Error().Err(err).Msg("marshal transaction")
			}
		} else {
			log.Error().Err(err).Msg("build transaction")
		}

		fmt.Printf("Elapsed time %v\n", time.Since(start))
	} else {
		log.Error().Err(err).Msg("Init build")
	}
}

/*
	Aggregate cardano events tx, utxo,stxi, asset, meta, datum, mint,witp,witn etc ...
    and output trans ( cardano transaction)
*/

func BuildTrans(cmd *cobra.Command, args []string) {

	var (
		err     error
		filters types.Filters
		Filter  interface{}
		req     = mongodb.MongoDB{
			Option:     &clientOption,
			Database:   database,
			Collection: collection,
		}
		req1 = mongodb.MongoDB{
			Option:     &clientOption,
			Database:   database,
			Collection: collection,
		}
		findOptions options.FindOptions
	)

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Error().Err(err).Msgf("logging to file %s", logger, err)
	}

	/*
		mongodb connection
	*/
	req.Uri = "mongodb://" + mongoUrl + "/?ssl=" + mongoSSL
	if req.Client, err = req.Connect(); err != nil {
		log.Error().Err(err).Msgf("Connect")
		return
	}
	defer req.DisConnect()

	req1.Uri = req.Uri
	if req1.Client, err = req1.Connect(); err != nil {
		log.Error().Err(err).Msgf("Connect")
		return
	}
	defer req1.DisConnect()

	var (
		ns       = []byte("buildtrans")
		key      = []byte("lastid")
		badgerDB *db.BadgerDB
		start    = time.Now()
		objectID primitive.ObjectID
	)
	if badgerDB, err = OpenBdb(dataDir); err != nil {
		log.Error().Err(err).Msgf("Opening badger Db %s", dataDir)
		return
	}

	if filter != "" {
		opVal := filters.ValidOp()
		if filter != "" {
			if filters, err = filters.ParseOp(filter, opVal); err != nil {
				LogError(err, "parse operator")
				return
			}
		}
		Filter = filters.BuildFilter()
	} else {
		if fromId != "" {
			objectID, err = primitive.ObjectIDFromHex(fromId)
		} else {
			value, err := badgerDB.Get(ns, key)
			if err == nil {
				objectID, err = primitive.ObjectIDFromHex(string(value))
			}
		}
		//  if objectID is valid
		if err == nil {
			Filter = bson.D{{"_id", bson.D{{"$gt", objectID}}}}
		} else {
			log.Error().Err(err).Msgf("Hex to primitive ObjectID")
			Filter = bson.D{{}}
		}

	}

	if concurrent {
		utils.SetCPU("50%")
	}
	/* loop on tx collection */
	req.Collection = "tx"
	var result []types.Tx

	/*
			 Ex:  collection.Find(bson.M{"brandId": body.BrandID, "category": body.Category, "$and": AndQuery})
		       bson.D{{"fingerprint",bson.D{{"$gte","55813901.tx"}}}}
	*/

	var total, Total, Terror = 0, 0, 0
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	coll := req.Client.Database(req.Database).Collection(req.Collection)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cur, err := coll.Find(ctx, Filter, &findOptions)
	var k = 0
	if cur != nil {
		for cur.Next(ctx) {
			var res types.Tx
			if err = cur.Decode(&res); err == nil {
				result = append(result, res)
				k++
			} else {
				log.Error().Err(err).Msgf("Decoding document ID %v", res.ID)
			}
		}
	}

	if k > 0 {
		objectId := result[k-1].ID
		key = []byte("_lastid_")
		badgerDB.Set(ns, key, []byte((objectId).Hex()))
		log.Info().Msgf("ObjectId %s before processing", objectId.String())
	}

	var (
		result1 []types.Tx
		flags   = SetFlags()
	)
	for _, v := range result {
		result1 = append(result1, v)
		total++
		if total%bulk == 0 {

			if concurrent {
				//tbuild, terror := BuildConTrans(result1, &req, &req1)
				tbuild, terror := lib.BuildTransCon(flags, result1, &req, &req1)
				Total += tbuild
				Terror += terror
			} else {
				// tbuild, terror := BuildSeqTrans(result1, &req, &req1)
				tbuild, terror := lib.BuildTransSeq(flags, result1, &req, &req1)
				Total += tbuild
				Terror += terror
			}

			result1 = []types.Tx{}

		}

		if total == len(result) {
			if concurrent {
				// tbuild, terror := BuildConTrans(result1, &req, &req1)
				tbuild, terror := lib.BuildTransCon(flags, result1, &req, &req1)
				Total += tbuild
				Terror += terror
			} else {
				// tbuild, terror := BuildSeqTrans(result1, &req, &req1)
				tbuild, terror := lib.BuildTransSeq(flags, result1, &req, &req1)
				Total += tbuild
				Terror += terror
			}
			result1 = []types.Tx{}
		}

	}
	if k > 0 {
		objectId := result[k-1].ID
		key = []byte("lastid")
		badgerDB.Set(ns, key, []byte((objectId).Hex()))
		log.Info().Msgf("ObjectId %s after processing", objectId.String())

	}

	log.Info().Stack().Msgf("Number of processed transactions: %d - uploaded: %d - error: %d - Total Elapsed time: %v", len(result), Total, Terror, time.Since(start))
}
