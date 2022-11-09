package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/paulmatencio/oura-go/db"
	"github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	buildBlkCmd = &cobra.Command{
		Use:   "buildBlock",
		Short: "Build a given bloc",
		Long:  ``,
		Run:   buildBlock,
	}
	buildBlksCmd = &cobra.Command{
		Use:   "buildBlocks",
		Short: "Build many blocks",
		Long:  ``,
		Run:   buildBlocks,
	}
)

func init() {
	rootCmd.AddCommand(buildBlkCmd)
	rootCmd.AddCommand(buildBlksCmd)
	initBuildBlk(buildBlkCmd)
	initBuildBlks(buildBlksCmd)
	if database == "" {
		database = mongoDatabase
	}
	if database == "" {
		log.Fatal().Msg("mongodb database is missing")
	}
}

func initBuildBlk(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().BoolVarP(&printIt, "print", "p", true, "print the result")
}

func initBuildBlks(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "", "local database directory. If no a full path, it with be prefix by home directory")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().StringVarP(&fromId, "from-id", "", "", "from object id")
	cmd.Flags().Int64VarP(&limit, "limit", "", 10, "max returned")
	cmd.Flags().BoolVarP(&printIt, "print", "p", false, "print the result")
	cmd.Flags().BoolVarP(&concurrent, "concurrent", "C", false, "Concurrent upload")
	cmd.Flags().IntVarP(&bulk, "bulk", "n", 10, "max bulkload")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger full file path")
}

type Responses struct {
	Block types.BlckN
	Err   error
}

func buildBlock(cmd *cobra.Command, args []string) {

	flags := SetFlags()
	if req, Filter, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
		defer req.DisConnect()
		BuildBlk(flags, req, Filter)
	} else {
		log.Error().Err(err).Msg("InitBuild")
	}
}

func buildBlocks(cmd *cobra.Command, args []string) {

	var (
		ns       = []byte("buildblock")
		key      = []byte("lastid")
		badgerDB *db.BadgerDB
		opts     options.FindOptions
		start    = time.Now()
		objectID primitive.ObjectID
	)
	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}
	log.Info().Msgf("Upload %v", upload)
	fmt.Printf("Upload %v", upload)
	flags := SetFlags()
	// req is used for read
	req, Filter, err := lib.InitReq(flags, mongoUrl, filter)
	if err != nil {
		log.Error().Stack().Msgf(" Error InitBuild req %v", err)
		return
	}
	defer req.DisConnect()

	// req1 will be used for mongodb write
	req1 := mongodb.MongoDB{
		Option:   req.Option,
		Uri:      req.Uri,
		Database: req.Database,
	}
	if req1.Client, err = req1.Connect(); err != nil {
		log.Error().Msgf("%v", err)
		return
	}
	defer req1.DisConnect()

	if badgerDB, err = OpenBdb(dataDir); err != nil {
		log.Error().Err(err).Msg("opening badger db")
		return
	}
	if filter == "" {
		if fromId != "" {
			objectID, err = primitive.ObjectIDFromHex(fromId)
		} else {
			value, err := badgerDB.Get(ns, key)
			if err == nil {
				objectID, err = primitive.ObjectIDFromHex(string(value))
				objectID, err = primitive.ObjectIDFromHex(string(value))
			}
		}
		if err == nil {
			Filter = bson.D{
				{"_id", bson.D{
					{"$gt", objectID},
				}}}
		} else {
			log.Error().Err(err).Msg("Hex to primitive ObjectID")
			Filter = bson.D{{}}
		}
	}

	if concurrent {
		utils.SetCPU("50%")
	}
	/* loop on tx collection */
	req.Collection = "blck"
	var result []types.Blck

	var total, Total, Terror = 0, 0, 0
	if limit > 0 {
		opts.SetLimit(limit)
	}
	coll := req.Client.Database(req.Database).Collection(req.Collection)
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cur, err := coll.Find(ctx, Filter, &opts)
	var k = 0
	if cur != nil {
		for cur.Next(ctx) {
			var res types.Blck
			if err = cur.Decode(&res); err == nil {
				result = append(result, res)
				k++
			}
		}
	}
	if k > 0 {
		objectId := result[k-1].ID
		key = []byte("_lastid_")
		badgerDB.Set(ns, key, []byte((objectId).Hex()))
		log.Info().Msgf("ObjectId %s before processing", objectId.String())
	}

	var result1 []types.Blck
	for _, v := range result {
		result1 = append(result1, v)
		total++
		if total%bulk == 0 {

			if concurrent {
				tbuild, terror := BuildConBlocks(flags, result1, req, &req1)
				Total += tbuild
				Terror += terror
			} else {
				tbuild, terror := BuildSeqBlocks(flags, result1, req)
				Total += tbuild
				Terror += terror
			}
			// reset []result1
			result1 = []types.Blck{}
		}

		if total == len(result) {
			if concurrent {
				tbuild, terror := BuildConBlocks(flags, result1, req, &req1)
				Total += tbuild
				Terror += terror
			} else {
				tbuild, terror := BuildSeqBlocks(flags, result1, req)
				Total += tbuild
				Terror += terror
			}
			// reset []result1
			result1 = []types.Blck{}
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

func BuildBlk(flags types.Options, req *mongodb.MongoDB, Filter interface{}) {
	var (
		blk  types.Blck
		req2 = &mongodb.MongoDB{
			Option:   req.Option,
			Uri:      req.Uri,
			Database: req.Database,
		}
		opt1 *options.FindOneOptions
	)
	if req2.Client, err = req2.Connect(); err != nil {
		log.Error().Msgf("%v", err)
		return
	}
	defer req2.DisConnect()
	req.Collection = "blck"
	if r, err := lib.FindOne(opt1, blk, Filter, req); err == nil {
		if block, err := BuildBlock1(flags, r, req, req2); err == nil {
			if b, err := json.Marshal(block); err == nil {
				utils.PrintJson(string(b))
			} else {
				log.Error().Stack().Err(err).Msgf("tx %v", err)
			}
		}
	} else {
		log.Error().Err(err).Msgf("blck %v", Filter)
	}

	return
}

func BuildSeqBlocks(flags types.Options, result []types.Blck, req *mongodb.MongoDB) (total int, terror int) {

	var (
		blockN types.BlckN
		err    error
	)

	for _, r := range result {
		if blockN, err = BuildBlock1(flags, r, req, nil); err == nil {
			if _, err := req.InsertOne(blockN); err != nil {
				log.Error().Err(err).Msg("Insert one")
				terror += 1
			} else {
				total += 1
			}
		}

	}

	return
}

func BuildConBlocks(flags types.Options, result []types.Blck, req *mongodb.MongoDB, req1 *mongodb.MongoDB) (total int, terror int) {

	var (
		ch                 = make(chan *Responses)
		Blocks             []interface{}
		nRequest, nReceive = 0, 0
		req2               = &mongodb.MongoDB{
			Option:   req.Option,
			Uri:      req.Uri,
			Database: req.Database,
		}
	)

	if req2.Client, err = req2.Connect(); err != nil {
		log.Error().Msgf("%v", err)
		return
	}
	defer req2.DisConnect()
	nRequest = len(result)

	/*
		BuildBlk1b  is called concurrently  for every block of the result array
	*/
	for _, r := range result {
		go func(r types.Blck, req *mongodb.MongoDB, req2 *mongodb.MongoDB) {
			var resp Responses
			resp.Block, resp.Err = BuildBlock1(flags, r, req, req2)
			ch <- &resp
		}(r, req, req2)
	}

	nReceive = 0
	for {
		if len(result) == 0 {
			return 0, 0
		}
		select {
		case rec := <-ch:
			nReceive++
			if rec.Err == nil {
				Blocks = append(Blocks, rec.Block)
			} else {
				terror += 1
				log.Error().Err(err).Msgf("building block %d", rec.Block.Block.Number)
			}
			if nReceive == nRequest {
				total += len(Blocks)
				log.Info().Msgf("Concurrent bulk uploading %d block documents", total)
				req1.Collection = "block"
				if _, err := req1.InsertMany(Blocks); err != nil {

					log.Error().Err(err).Msg("Insert Many")
					terror += 1
				}
				return total, terror
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Printf(".")
		}
	}
}

func BuildBlock1(flags types.Options, blk types.Blck, req *mongodb.MongoDB, req2 *mongodb.MongoDB) (block types.BlckN, err error) {

	var (
		trans                                                                                  types.Trans
		fees, totalOutput, mintCount                                                           int64
		inputCount, metaCount, outputCount, datumCount, rdmrCount, plutusWCount, nativeWCount  int
		cip25Count, poolRegCount, stakeDeleCount, poolRetiCount, stakeRegCount, stakeDereCount int
		cip20Count                                                                             int
		findOptions                                                                            *options.FindOptions
		opt1                                                                                   *options.FindOneOptions
	)

	req.Collection = "trans"
	Filter2 := bson.M{
		"context.block_number": blk.Block.Number,
		"context.block_hash":   blk.Block.Hash,
	}
	block.CopyFrom(&blk)
	findOptions.SetLimit(int64(blk.Block.TxCount))
	if _, rr, err := lib.Find(findOptions, trans, Filter2, req); err == nil {
		for _, v := range rr {
			block.Block.Transactions = append(block.Block.Transactions, v.Transaction)
			block.Block.TxMeta = append(block.Block.TxMeta, v.TxMeta)
			fees += v.Transaction.Fee
			totalOutput += v.Transaction.TotalOutput
			inputCount += v.TxMeta.InputCount
			outputCount += v.TxMeta.OutputCount
			mintCount += v.TxMeta.MintCount
			metaCount += v.TxMeta.MetaCount
			datumCount += v.TxMeta.PlutusDatumCount
			rdmrCount += v.TxMeta.PlutusRdmrCount
			plutusWCount += v.TxMeta.PlutusWitnessesCount
			nativeWCount += v.TxMeta.NativeWitnessesCount
			cip25Count += v.TxMeta.Cip25AssetCount
			/*
				poolRegCount += v.TxMeta.PoolRegistrationCount
				poolRetiCount += v.TxMeta.PoolRetirementCount
				stakeDeleCount += v.TxMeta.StakeDelegationCount
				stakeRegCount += v.TxMeta.StakeRegistrationCount
				stakeDereCount += v.TxMeta.StakeDeregistrationCount
			*/
			Filter3 := bson.M{"context.tx_hash": v.Context.TxHash}

			/*
			   Continue with pool
			*/

			req.Collection = "pool"
			var pn types.Pool
			if p, err := lib.FindOne(opt1, pn, Filter3, req); err == nil {
				if p.PoolRegistration.PoolId != "" {
					poolRegCount += 1
				}
			}
			req.Collection = "reti"
			var ret types.Reti
			if reti, err := lib.FindOne(opt1, ret, Filter3, req); err == nil {
				if reti.PoolRetirement.Pool != "" {
					poolRetiCount += 1
				}
			}

			req.Collection = "dele"
			var del types.Dele
			if dele, err := lib.FindOne(opt1, del, Filter3, req); err == nil {
				if dele.StakeDelegation.PoolHash != "" {
					stakeDeleCount += 1
				}
			}

			req.Collection = "skre"
			var skr types.Skre
			if skre, err := lib.FindOne(opt1, skr, Filter3, req); err == nil {
				if skre.StakeRegistration.Credential.AddrKeyhash != "" {
					stakeRegCount += 1
				}
			}

			req.Collection = "skde"
			var skd types.Skde
			if skde, err := lib.FindOne(opt1, skd, Filter3, req); err == nil {
				if skde.StakeDeregistration.Credential.AddrKeyhash != "" {
					stakeDereCount += 1
				}
			}

			for _, m := range v.Metadata {
				metadata, _ := PrimitiveDtoMap(m)
				if metadata["label"] == "674" {
					cip20Count += 1
				}
			}
		}
		if len(rr) > 0 {
			block.Block.Fees = fees
			block.Block.TotalOutput = totalOutput
			block.Block.InputCount = inputCount
			block.Block.OutputCount = outputCount
			block.Block.MintCount = mintCount
			block.Block.MetaCount = metaCount
			block.Block.NativeWitnessesCount = nativeWCount
			block.Block.PlutusDatumCount = datumCount
			block.Block.PlutusRdmrCount = rdmrCount
			block.Block.PlutusWitnessesCount = plutusWCount
			block.Block.Cip25AssetCount = cip25Count
			block.Block.Cip20Count = cip20Count
			block.Block.PoolRegistrationCount = poolRegCount
			block.Block.PoolRetirementCount = poolRetiCount
			block.Block.StakeDelegationCount = stakeDeleCount
			block.Block.StakeRegistrationCount = stakeRegCount
			block.Block.StakeDeregistrationCount = stakeDereCount
		}

	} else {
		log.Error().Err(err).Msgf("Find tx %v", Filter2)
	}

	Filter4 := bson.D{{"height", blk.Block.Number}}
	var (
		blockF types.BlockfB
		r      types.BlockfB
		err1   error
	)
	if req2 != nil {
		req2.Collection = "blockf"
		r, err1 = lib.FindOne(opt1, blockF, Filter4, req2)
	} else {
		req.Collection = "blockf"
		r, err1 = lib.FindOne(opt1, blockF, Filter4, req)
	}
	if err1 == nil {
		block.Block.SlotLeader = r.SlotLeader
		block.Block.NextHash = r.NextBlock
		block.Block.Confirmations = r.Confirmations
	} else {
		log.Error().Err(err1).Msgf("Filter: %v", Filter2)
	}
	return
}
