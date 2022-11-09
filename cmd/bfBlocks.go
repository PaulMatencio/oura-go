package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/blockfrost/blockfrost-go"
	"github.com/paulmatencio/oura-go/db"
	"github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var (
	retBfBlockCmd = &cobra.Command{
		Use:   "retBfBlock",
		Short: "Retrieve a cardano block from blockfrost.io",
		Long:  ``,
		Run:   retBlock,
	}

	retBfBlocksCmd = &cobra.Command{
		Use:   "retBfBlocks",
		Short: "Retrieve cardano blocks from blockfrost.io",
		Long:  ``,
		Run:   retBlocks,
	}
	saveBfBlocksCmd = &cobra.Command{
		Use:   "saveBfBlocks",
		Short: "save cardano blocks from blockfrost.io to local database ",
		Long:  ``,
		Run:   saveBlocks,
	}

	maxReturned, page int
	pretty            bool
	blockNumber       int64
	loop              int
)

func init() {
	rootCmd.AddCommand(retBfBlockCmd)
	rootCmd.AddCommand(retBfBlocksCmd)
	rootCmd.AddCommand(saveBfBlocksCmd)
	initRetBlock(retBfBlockCmd)
	initRetBlocks(retBfBlocksCmd)
	initSaveBlocks(saveBfBlocksCmd)

}

func initRetBlock(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&blockNumber, "block-number", "N", 0, "block number")
	cmd.Flags().BoolVarP(&printIt, "print", "", false, "print out in pretty json")
}

func initRetBlocks(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&blockNumber, "block-number", "N", 0, "block number")
	cmd.Flags().IntVarP(&maxReturned, "max-returned", "m", 100, "maximum number of return per resource")
	cmd.Flags().IntVarP(&page, "page", "p", 0, fmt.Sprintf("By default, its return %d at a time. You have to use page=2 to list through the results.", maxReturned))
	cmd.Flags().BoolVarP(&printIt, "print", "", false, "print out in pretty json")
	cmd.Flags().IntVarP(&loop, "loop", "L", 1, "loop count")
}

func initSaveBlocks(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "", "local database directory. If no a full path, it with be prefix by home directory")
	cmd.Flags().Int64VarP(&blockNumber, "block-number", "N", 0, "from block number")
	cmd.Flags().IntVarP(&maxReturned, "max-returned", "m", 100, "maximum number of return per resource")
	cmd.Flags().IntVarP(&page, "page", "p", 0, fmt.Sprintf("By default, its return %d at a time. You have to use page=2 to list through the results.", maxReturned))
	cmd.Flags().IntVarP(&loop, "loop", "L", 1, "loop count")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger Filename ")
	if database == "" {
		database = mongoDatabase
	}
	if database == "" {
		log.Fatal().Msg("mongodb database is missing")
	}
}

func retBlock(cmd *cobra.Command, args []string) {
	var blockf types.BlockF
	blockf.GetBlock(projectId, blockNumber)
	b, err := json.Marshal(&blockf.Block)
	if err == nil {
		utils.PrintJson(string(b))
	}
}

func retBlocks(cmd *cobra.Command, args []string) {

	var (
		blockf  types.BlockF
		options = blockfrost.APIClientOptions{
			ProjectID: projectId,
		}
		getOptions = types.GetBlockOptions{
			ApiClient: blockfrost.NewAPIClient(options),
			ApiQueryParams: blockfrost.APIQueryParams{
				Count: maxReturned,
				Page:  page,
			},
		}
	)
	n := 0
	for {
		blockf.GetBlocks(getOptions, blockNumber)
		if blockf.Blocks != nil {
			if printIt {
				for _, v := range blockf.Blocks {
					b, err := json.Marshal(&v)
					if err == nil {
						utils.PrintJson(string(b))
					}
				}
			} else {
				fmt.Println(len(blockf.Blocks), blockf.LastBlock)
			}
			n++
			if loop != 0 && n >= loop {
				return
			} else {
				blockNumber = int64(blockf.LastBlock)
			}
		} else {
			return
		}

	}
}

func saveBlocks(cmd *cobra.Command, args []string) {
	var (
		blockf  types.BlockF
		options = blockfrost.APIClientOptions{
			ProjectID: projectId,
		}
		getOptions = types.GetBlockOptions{
			ApiClient: blockfrost.NewAPIClient(options),
			ApiQueryParams: blockfrost.APIQueryParams{
				Count: maxReturned,
				Page:  page,
			},
		}
		req      *mongodb.MongoDB
		err      error
		ns       = []byte("blockf")
		key      = []byte("lastblock")
		badgerDB *db.BadgerDB
	)

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Error().Err(err).Msgf("logging to file %s", logger)
	}
	if badgerDB, err = OpenBdb(dataDir); err != nil {
		log.Error().Err(err).Msgf("Opening badger db %s", dataDir)
		return
	}
	v, err := badgerDB.Get(ns, key)
	if err == nil && blockNumber == 0 {
		//v1, _ := strconv.Atoi(string(v))
		blockNumber, _ = strconv.ParseInt(string(v), 10, 64)
	}
	flags := SetFlags()

	if req, _, err = lib.InitReq(flags, mongoUrl, filter); err == nil {
		req.Collection = "blockf"
	} else {
		log.Error().Msgf("Init connection to mongodb failed %v", err)
		return
	}
	l := 0 /* loop number */
	start := time.Now()
	for {
		start1 := time.Now()
		blockf.GetBlocks(getOptions, blockNumber)
		if blockf.Blocks != nil {
			if _, err := req.InsertMany(blockf.Blocks); err != nil {
				log.Error().Err(err).Msg("Insert many")
				return
			} else {
				lastBlock := blockf.LastBlock
				log.Info().Msgf("%d blocks were uploaded - Uploading next %d blocks from block # %d - Elapsed time %v/%v", len(blockf.Blocks), maxReturned, lastBlock, time.Since(start1), time.Since(start))
				//badgerDB.Set(ns, key, []byte(strconv.Itoa(lastBlock)))
				badgerDB.Set(ns, key, []byte(strconv.FormatInt(int64(lastBlock), 10)))
			}
			l++
			if loop != 0 && l >= loop {
				log.Info().Msgf("Last block number %d", blockf.LastBlock)
				return
			} else {
				blockNumber = int64(blockf.LastBlock)
			}
		} else {
			return
		}
	}
}
