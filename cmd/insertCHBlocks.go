/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
//

package cmd

import (
	"context"
	"github.com/paulmatencio/oura-go/lib"
	"strings"
	//ctypes "github.com/paulmatencio/clickhouse/types"
	ctypes "github.com/paulmatencio/oura-go/ch/types"
	// cutils "github.com/paulmatencio/clickhouse/utils"
	cutils "github.com/paulmatencio/oura-go/ch/utils"
	"github.com/paulmatencio/oura-go/db"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// insertBlocksCmd represents the insertBlocks command
var (
	insertBlocksCmd = &cobra.Command{
		Use:   "insertChBlocks",
		Short: "insert blocks to click house",
		Long:  ``,
		Run:   insertChBlks,
	}
	Urls   string
	chAddr []string
	create bool
)

func init() {
	rootCmd.AddCommand(insertBlocksCmd)
	initChblock(insertBlocksCmd)

}

func initChblock(cmd *cobra.Command) {

	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "", "badger database directory. If no a full path, it with be prefix by home directory")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().StringVarP(&fromId, "from-id", "", "", "from object id")
	cmd.Flags().Int64VarP(&limit, "limit", "", 100, "max returned")
	cmd.Flags().StringVarP(&Urls, "ch-urls", "C", "", "clickhouse urls separated by comma")
	cmd.Flags().StringVarP(&cDatabase, "ch-db", "", "cardano", "clickhouse database name")
	cmd.Flags().IntVarP(&bulk, "bulk", "n", 10, "max bulk load number")
	cmd.Flags().BoolVarP(&create, "create-table", "", false, "create destination table. It will delete the existing one")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger Filename ")
}

func insertChBlks(cmd *cobra.Command, args []string) {

	var (
		ns       = []byte("chblock")
		key      = []byte("lastid")
		badgerDB *db.BadgerDB
		opts     options.FindOptions
		start    = time.Now()
		objectID primitive.ObjectID
		chblock  ctypes.ChBlock
	)

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}

	if limit > 0 {
		opts.SetLimit(limit)
	}

	if Urls == "" {
		if chAddrs != "" {
			Urls = chAddrs
		} else {
			log.Fatal().Msgf(" clickhouse addresses are missing %s ", chAddrs)
		}
	}
	chAddr = strings.Split(Urls, ",")

	if database == "" {
		if mongoDatabase != "" {
			database = mongoDatabase
		} else {
			log.Fatal().Msg("mongodb database is missing")
		}
	}

	if cDatabase == "" {
		if chDatabase != "" {
			cDatabase = chDatabase
		} else {
			log.Fatal().Msg("clickhouse database is missing")
		}
	}

	flags := SetFlags()
	req, Filter, err := lib.InitReq(flags, mongoUrl, filter)
	if err != nil {
		log.Error().Stack().Msgf(" Error InitBuild req %v", err)
		return
	}

	// open badger db connection
	if badgerDB, err = OpenBdb(dataDir); err != nil {
		log.Error().Msgf("Error opening badger db %v", err)
		return
	}

	//  open Clickhouse connection
	conn, err := cutils.Connect(chAddr, chDatabase, chUser, chPassword)
	if err != nil {
		log.Error().Msgf("connection failed %v", err)
		return
	}

	if filter == "" {
		if fromId != "" {
			objectID, err = primitive.ObjectIDFromHex(fromId)
		} else {
			value, err := badgerDB.Get(ns, key)
			if err == nil {
				objectID, err = primitive.ObjectIDFromHex(string(value))
			}
		}
		if err == nil {
			Filter = bson.D{{"_id", bson.D{{"$gt", objectID}}}}
		} else {
			log.Error().Msgf("Hex to primitive ObjectID %v", err)
			Filter = bson.D{{}}
		}
	}

	chblock.Table = "chblock"
	if create {
		err = chblock.Drop(conn)
		if err != nil {
			log.Error().Msgf("drop table failed %v", err)
		}

		err = chblock.CreateTable(conn)
		if err != nil {
			log.Error().Msgf("create table failed %v", err)
			return
		}
	}

	req.Collection = "block"
	var result []types.BlckN
	var total, tInsert, tError = 0, 0, 0

	coll := req.Client.Database(req.Database).Collection(req.Collection)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cur, err := coll.Find(ctx, Filter, &opts)
	var k = 0
	if cur != nil {
		for cur.Next(ctx) {
			var res types.BlckN
			if err = cur.Decode(&res); err == nil {
				result = append(result, res)
				k++
			}
		}
	}
	if k > 0 {
		objectId := result[0].ID
		key = []byte("_lastid_")
		badgerDB.Set(ns, key, []byte((objectId).Hex()))
		log.Info().Msgf("ObjectId %s before processing", objectId.String())
	}

	var (
		result1 []ctypes.Block
	)

	for _, v := range result {
		var block ctypes.Block
		CopierBlock(&block, &v.Block)
		block.TimeStamp = time.Unix(v.Context.Timestamp, 0)
		result1 = append(result1, block)
		total++
		if total%bulk == 0 {
			err = chblock.PrepareBatch(conn)
			if err != nil {
				log.Error().Msgf("prepare failed %v", err)
				return
			} else {

				log.Trace().Msgf("prepare batch %v", chblock.Batch)
			}
			if len(result1) > 0 {
				chblock.Blocks = result1
				err = chblock.BulkInsertStruct()
				if err != nil {
					tError++
					log.Error().Msgf("Bulk Insert structure  failed %v", err)
				} else {
					tInsert++
					log.Trace().Msgf("Bulk Insert of %d row done", len(result1))
				}
			}
			result1 = []ctypes.Block{}
		}

		if total == len(result) {
			err = chblock.PrepareBatch(conn)
			if err != nil {
				log.Error().Msgf("prepare failed %v", err)
				return
			} else {
				log.Trace().Msgf("batch %v", chblock.Batch)
			}
			if len(result1) > 0 {
				chblock.Blocks = result1
				err = chblock.BulkInsertStruct()
				if err != nil {
					tError++
					log.Error().Msgf("Bulk Insert failed %v", err)
				} else {
					tInsert++
					log.Trace().Msgf("Bulk Insert of %d row done", len(result1))
				}
			}
			result1 = []ctypes.Block{}
		}

	}
	if k > 0 {
		objectId := result[k-1].ID
		key = []byte("lastid")
		badgerDB.Set(ns, key, []byte((objectId).Hex()))
		log.Info().Msgf("ObjectId %s after processing", objectId.String())
		/*
			fmt.Println("Print result")
			var r []ctypes.Block
			if err = conn.Select(ctx, &r, "SELECT *  FROM chblock"); err != nil {
				log.Error().Msgf("%v", err)
			}
			for _, v := range r {
				fmt.Println(v)
			}
		*/
	}
	req.DisConnect()
	log.Info().Stack().Msgf("# transactions: %d - # bulk-inserted: %d - # error: %d - Total Elapsed time: %v", len(result), tInsert, tError, time.Since(start))
}

func CopierBlock(to *ctypes.Block, from *types.BlockN) {
	to.BlockNumber = uint64(from.Number)
	to.Slot = uint64(from.Slot)
	to.BodySize = uint32(from.BodySize)
	to.Epoch = uint32(from.Epoch)
	to.EpochSlot = uint32(from.EpochSlot)
	to.Era = from.Era
	to.Hash = from.Hash
	to.IssuerVkey = from.IssuerVkey
	to.SlotLeader = from.SlotLeader
	to.TxCount = uint32(from.TxCount)
	to.Fees = uint64(from.Fees)
	to.TotalOutput = uint64(from.TotalOutput)
	to.InputCount = uint32(from.InputCount)
	to.OutputCount = uint32(from.OutputCount)
	to.MintCount = uint64(from.MintCount)
	to.MetaCount = uint32(from.MetaCount)
	to.NativeWitnessesCount = uint32(from.NativeWitnessesCount)
	to.PlutusDatumCount = uint32(from.PlutusDatumCount)
	to.PlutusRdmrCount = uint32(from.PlutusRdmrCount)
	to.PlutusWitnessesCount = uint32(from.PlutusWitnessesCount)
	to.Cip25AssetCount = uint32(from.Cip25AssetCount)
	to.Cip20Count = uint32(from.Cip20Count)
	to.PoolRegistrationCount = uint32(from.PoolRegistrationCount)
	to.PoolRetirementCount = uint32(from.PoolRetirementCount)
	to.StakeDelegationCount = uint32(from.StakeDelegationCount)
	to.StakeRegistrationCount = uint32(from.StakeRegistrationCount)
	to.StakeDeregistrationCount = uint32(from.StakeDeregistrationCount)
	to.Confirmations = uint32(from.Confirmations)
}
