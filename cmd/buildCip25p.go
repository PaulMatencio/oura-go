/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"fmt"
	"github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Response struct {
	Total  int
	Terror int
}

// buildMetasCmd represents the buildMetas command
var (
	buildMetasCmd = &cobra.Command{
		Use:   "buildCip25p",
		Short: "build cip25p events",
		Long:  `.`,
		Run:   buildCip25,
	}
	checkDup bool
)

func init() {
	rootCmd.AddCommand(buildMetasCmd)
	initMetasCmd(buildMetasCmd)

}

func initMetasCmd(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "badger", "local database directory. If no a full path, it with be prefix by home directory")
	cmd.Flags().Int64VarP(&blockNumber, "block-number", "N", 0, "from block number")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().Int64VarP(&limit, "limit", "", 100, "max returned")
	cmd.Flags().StringVarP(&fromId, "from-id", "", "", "from tx object id")
	// cmd.Flags().Int64VarP(&fromtime, "from-time", "f", 0, "max number of loop")
	cmd.Flags().IntVarP(&loop, "loop", "", 1, "max number of loop")
	cmd.Flags().BoolVarP(&concurrent, "concurrent", "C", false, "concurrent upload")
	cmd.Flags().BoolVarP(&checkDup, "checkDup", "", false, "check duplicated fingerprint")
	cmd.Flags().BoolVarP(&upload, "upload", "U", false, "upload to mongodb ")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger Filename ")
	cmd.Flags().BoolVarP(&printIt, "print", "p", true, "print the result")

}

func buildCip25(cmd *cobra.Command, args []string) {

	var (
		ns       = []byte("meta")
		key      = []byte("lastid")
		filters  types.Filters
		Filter   interface{}
		objectID primitive.ObjectID
	)

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}

	if dataDir == "" {
		log.Warn().Msgf("Badger DB directory %s is missing", dataDir)
		return
	}

	if badgerDB, err = OpenBdb(dataDir); err != nil {
		LogError(err, fmt.Sprintf("Opening Badger DB:%s", dataDir))
		return
	} else {
		log.Info().Msgf("Badger name space:%s", string(ns))
	}

	if filter != "" {
		opVal := filters.ValidOp()
		if filter != "" {
			if filters, err = filters.ParseOp(filter, opVal); err != nil {
				LogError(err, PeOperator)
				return
			}
		}
		Filter = filters.BuildFilter()
	} else {
		if fromId != "" {
			objectID, err = primitive.ObjectIDFromHex(fromId)
		} else {
			var value []byte
			value, err = badgerDB.Get(ns, key)
			if err == nil {
				objectID, err = primitive.ObjectIDFromHex(string(value))
			}
		}
		// if objectId is valid
		if err == nil {
			Filter = bson.M{
				"metadata.label": "721",
				"_id":            bson.M{"$gt": objectID},
			}
		} else {
			LogError(err, "Hex to primitive ObjectID. ObjectId is bypassed")
			Filter = bson.M{
				"metadata.label": "721",
			}
		}
	}

	log.Info().Msgf("Upload %v", upload)
	log.Info().Msgf("Name space:%s - Key:%s - Filter:%v\n", string(ns), string(key), Filter)
	flags := SetFlags()

	if req, _, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
		var (
			Nloop                    = 0
			NMetas, NCip25s, NErrors = 0, 0, 0
			start                    = time.Now()
			req1                     = mongodb.MongoDB{
				Uri:        req.Uri,
				Option:     req.Option,
				Database:   req.Database,
				Collection: "cip25p",
			}
			req2 = mongodb.MongoDB{
				Uri:        req.Uri,
				Option:     req.Option,
				Database:   req.Database,
				Collection: "mint",
			}
		)
		defer req.DisConnect()
		if req1.Client, err = req1.Connect(); err != nil {
			log.Error().Err(err).Msg(MeConnect)
		}
		defer req1.DisConnect()

		if req2.Client, err = req2.Connect(); err != nil {
			log.Error().Err(err).Msg(MeConnect)
		}
		defer req2.DisConnect()
		req.Collection = "meta"
		for {
			var (
				start1                   = time.Now()
				nMetas, nCip25s, nErrors int
				objectID                 primitive.ObjectID
			)
			if !concurrent {
				// nMetas, nCip25s, nErrors, objectID = BuildC25Seq(req, Filter)
				nMetas, nCip25s, nErrors, objectID = lib.BuildC25Seq(flags, req, &req1, &req2, Filter)
			} else {
				// nMetas, nCip25s, nErrors, objectID = BuildC25Con(req, Filter)
				nMetas, nCip25s, nErrors, objectID = lib.BuildC25Con(flags, req, &req1, &req2, Filter)
			}
			Nloop++

			if upload && objectID != primitive.NilObjectID {
				if objectID != primitive.NilObjectID {
					key = []byte("lastid")
					badgerDB.Set(ns, key, []byte((objectID).Hex()))
				}
			}
			log.Info().Msgf("number of fetched metas:%d - number of cip25 uploaded:%d - number of errors:%d - elapsed time: %v ms - next objectId:%v ", nMetas, nCip25s, nErrors, time.Since(start1).Milliseconds(), objectID)
			NMetas += nMetas
			NCip25s += nCip25s
			NErrors += nErrors
			log.Info().Msgf("Total fetched metas  %d - Total cip25 uploaded:%d - Total errors:%d - Total elapsed time: %v ms", NMetas, NCip25s, NErrors, time.Since(start).Milliseconds())
			Filter = bson.M{
				"metadata.label": "721",
				"_id":            bson.M{"$gt": objectID},
			}
			if (Nloop < loop || loop == 0) && nCip25s > 0 {
				continue
			} else {
				break
			}
		}

	} else {
		LogError(err, MeInitReq)
	}
}
