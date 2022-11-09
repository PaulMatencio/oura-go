package cmd

import (
	"fmt"
	lib "github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	buildNftTransCmd = &cobra.Command{
		Use:   "buildNftTrans",
		Short: "Build all transactions for a given utxo address",
		Long:  ``,
		Run:   BuildNftTrans,
	}
	buildNftTranCmd = &cobra.Command{
		Use:   "buildNftTran",
		Short: "Build a Nft transaction for a given tx_hash",
		Long:  ``,
		Run:   BuildNftTran,
	}
)

var (
	filter, marketPlace, dataDir, logger string
	printIt, listPlaces                  bool
	concurrent, upload                   bool
	bulk                                 int
	fromId                               string
	objectId                             primitive.ObjectID
)

func init() {

	rootCmd.AddCommand(buildNftTransCmd)
	initBuildNftTrans(buildNftTransCmd)
	rootCmd.AddCommand(buildNftTranCmd)
	initBuildNftTran(buildNftTranCmd)
	if database == "" {
		database = mongoDatabase
	}
	if database == "" {
		log.Fatal().Msg("mongodb database is missing")
	}

}

func initBuildNftTran(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&marketPlace, "market-place", "", "", "market place name ")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().BoolVarP(&listPlaces, "list-market-places", "L", false, "List NFT marketplaces ")
	cmd.Flags().BoolVarP(&printIt, "print", "p", true, "print the result")
}

func initBuildNftTrans(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "badger", "local database directory. If no a full path, it with be prefix by home directory")
	cmd.Flags().StringVarP(&marketPlace, "market-place", "", "", "market place name ")
	cmd.Flags().StringVarP(&fromId, "from-id", "", "", "from tx object id")
	cmd.Flags().Int64VarP(&limit, "limit", "", 1000, "max returned")
	cmd.Flags().IntVarP(&loop, "loop", "", 1, "max number of loop")
	cmd.Flags().BoolVarP(&listPlaces, "list-market-places", "L", false, "List NFT marketplaces ")
	cmd.Flags().BoolVarP(&concurrent, "concurrent", "C", false, "concurrent upload")
	cmd.Flags().BoolVarP(&upload, "upload", "U", false, "upload to mongodb ")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger Filename ")
	cmd.Flags().BoolVarP(&printIt, "print", "p", true, "print the result")
}

/*
	Build an NFT transaction
*/
func BuildNftTran(cmd *cobra.Command, args []string) {

	flags := SetFlags()
	mktPlaces := lib.GetMarketPlaces(marketPlace, MktPlaces)
	if len(mktPlaces) > 0 {
		flags.SetMarketPlace(mktPlaces[0])
		if req, Filter, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
			var (
				req1 mongodb.MongoDB
				req2 mongodb.MongoDB
			)
			req1.Option, req2.Option = req.Option, req.Option
			req1.Uri, req2.Uri = req.Uri, req.Uri
			req1.Database, req2.Database = req.Database, req.Database
			if req1.Client, err = req1.Connect(); err != nil {
				log.Error().Err(err).Msg("req1 connection")
				return
			}
			defer req1.DisConnect()
			if req2.Client, err = req2.Connect(); err != nil {
				log.Error().Err(err).Msg("req2 connection")
				return
			}
			defer req2.DisConnect()
			log.Info().Msgf("Building nft transaction %v  for market-place %v ", Filter, mktPlaces)
			lib.BuildNftTran(flags, req, &req1, Filter)
		}
	}
}

/*
	Build multiple NFT transactions  for a specific marketplace
       get the smart contract address of  the given marketplace
       resume  the execution context ( get next mongoDB object ID) .
           open the nft badger db database  ( "nft" )
           get the last object id for this marketplace  ( badger database)
       if concurrent -> BuildNftTransCon
       else -> BuildNftTransSeq
*/
func BuildNftTrans(cmd *cobra.Command, args []string) {

	var (
		badger = "nft"               /* check config.yaml */
		ns     = []byte("trans_nft") /*name space */
		key    []byte
	)
	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}
	if listPlaces {
		for _, v := range MktPlaces {
			fmt.Printf("Name: %s Address: %s\n", v.Name, v.Address)
		}
		return
	}
	badgerDir := GetBadgerPath(badger)
	if badgerDir == "" {
		log.Error().Msgf("configure the badger %s", badger)
		return
	} else {
		log.Info().Msgf("Badger directory %s - name space %s", badgerDir, string(ns))
		if badgerDB, err = OpenBdb(badgerDir); err != nil {
			log.Error().Err(err).Msgf("Opening %s failed", badgerDir)
			return
		}
	}
	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}

	log.Info().Msgf("Upload %v", upload)
	flags := SetFlags()

	if marketPlace != "" {
		/*
			get the  marketplace attributes   for  the given marketplace
			The given marketplace can be generic  since a marketplace  may have multiple addresses
		*/
		mktPlaces := lib.GetMarketPlaces(marketPlace, MktPlaces)
		if len(mktPlaces) > 0 {
			for _, mktpl := range mktPlaces {
				flags.SetMarketPlace(mktpl)
				key = []byte(mktpl.Name)
				if fromId != "" {
					objectId, err = primitive.ObjectIDFromHex(fromId)
				} else {
					value, err := badgerDB.Get(ns, key)
					if err == nil {
						objectId, err = primitive.ObjectIDFromHex(string(value))
					}
				}

				log.Trace().Msgf("Reading from objectId  %v\n", objectId)

				// mongodb connection
				if req, _, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
					var (
						req1 mongodb.MongoDB
						req2 mongodb.MongoDB
					)

					req1.Option, req2.Option = req.Option, req.Option
					req1.Uri, req2.Uri = req.Uri, req.Uri
					req1.Database, req2.Database = req.Database, req.Database

					if req1.Client, err = req1.Connect(); err != nil {
						log.Error().Err(err).Msg("req1 connection")
						return
					}
					// defer	req1.DisConnect()
					if req2.Client, err = req2.Connect(); err != nil {
						log.Error().Err(err).Msg("req2 connection")
						return
					}

					// defer req1.DisConnect()
					Total, Terror, nLoop := 0, 0, 0
					start := time.Now()
					for {
						var (
							start1        = time.Now()
							total, terror int
							// objectId      primitive.ObjectID
						)
						fmt.Println("Object Id", objectId)
						if concurrent {
							total, terror, objectId = lib.BuildNftTransCon(flags, mktpl, objectId, req, &req1, &req2)
						} else {
							total, terror, objectId = lib.BuildNftTransSeq(flags, mktpl, objectId, req, &req1)
						}
						nLoop++
						if upload && objectId != primitive.NilObjectID {
							badgerDB.Set(ns, key, []byte((objectId).Hex()))
						}
						log.Info().Msgf("total documents uploaded:%d,total errors:%d,last objectId:%v,elapsed time: %v ms", total, terror, objectId, time.Since(start1).Milliseconds())
						if (nLoop <= loop || loop == 0) && total > 0 {
							Total += total
							Terror += terror
						} else {
							log.Info().Msgf("Total documents uploaded:%d,Total errors:%d,last objecID: %v,Elapsed time: %v ms", Total, Terror, objectId, time.Since(start).Milliseconds())
							break
						}
					}
					req.DisConnect()
					req1.DisConnect()
					req2.DisConnect()

				} else {
					log.Error().Err(err).Msg("InitReq error")
				}
			}
		} else {
			log.Warn().Msgf("Add <%s> market place to the config file %s", marketPlace, configFileUsed)
		}
	}
}
