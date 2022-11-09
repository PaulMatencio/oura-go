package cmd

import (
	//"encoding/json"
	"fmt"
	"github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/types"
	// options2 "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/paulmatencio/oura-go/utils"
	"github.com/spf13/cobra"
	// "log"
	"github.com/rs/zerolog/log"
)

// scanCmd represents the scan command
var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list limited number of documents of a given collection",
		Long:  ``,
		Run:   list,
	}
	listAllCmd = &cobra.Command{
		Use:   "listAll",
		Short: "list all documents of a given collection",
		Long:  ``,
		Run:   listAll,
	}
	opVal map[string]string
)

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(listAllCmd)
	initList(listCmd)
	initList(listAllCmd)
}

func initList(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "database", "d", "cardano", "database name")
	cmd.Flags().StringVarP(&collection, "collection", "c", "", "collection name")
	cmd.Flags().StringVarP(&filter, "filter-value", "f", "", "filters ")
	cmd.Flags().Int64VarP(&limit, "limit", "", 10, "max returned document")
}

func list(cmd *cobra.Command, args []string) {

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf(" logging to file %s - error %v ", logger, err)
	}
	if collection == "" {
		fmt.Println("collection is missing ")
		return
	}
	flags := SetFlags()
	var findOptions options.FindOptions
	findOptions.SetLimit(flags.Limit)
	if req, Filter, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
		defer req.DisConnect()
		req.Collection = collection
		switch collection {
		case "tx":
			var T []types.Tx
			// List(T, Filter, req)
			lib.List(&findOptions, T, Filter, req)
		case "utxo":
			var T []types.Utxo
			// List(T, Filter, req)
			lib.List(&findOptions, T, Filter, req)
		case "stxi":
			var T []types.Stxi
			// List(T, Filter, req)
			lib.List(&findOptions, T, Filter, req)
		case "asst":
			var T []types.OutputAsset
			lib.List(&findOptions, T, Filter, req)
			//List(T, Filter, req)
		case "meta":
			var T []types.Meta
			lib.List(&findOptions, T, Filter, req)
			//List(T, Filter, req)
		case "mint":
			var T []types.Mint
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "coll":
			var T []types.Coll
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "dtum":
			var T []types.Datum
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "witp":
			var T []types.Witp
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "blck":
			var T []types.Blck
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "witn":
			var T []types.Witn
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "rdmr":
			var T []types.Redeemer
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "trans":
			var T []types.Trans
			lib.List(&findOptions, T, Filter, req)
		case "trans_nft":
			var T []types.TransNft
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "cip25":
			var T []types.Cip25
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "cip25p":
			var T []types.Cip25p
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "pool":
			var T []types.Pool
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "dele":
			var T []types.Dele
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "reti":
			var T []types.Reti
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "skre":
			var T []types.Skre
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "skde":
			var T []types.Skde
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "block":
			var T []types.BlckN
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		case "blockf":
			var T []types.BlockfB
			lib.List(&findOptions, T, Filter, req)
			// List(T, Filter, req)
		default:
			log.Warn().Msgf("%v not on the list", collection)
		}

	} else {
		log.Error().Msgf("InitReq error : %v", err)
	}

}

func listAll(cmd *cobra.Command, args []string) {

	if collection == "" {
		fmt.Println("collection is missing ")
		return
	}
	flags := SetFlags()
	if req, _, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
		switch collection {
		case "tx":
			var T []types.Tx
			ListAll(T, filter, req)
		case "utxo":
			var T []types.Utxo
			ListAll(T, filter, req)
		case "stxi":
			var T []types.Stxi
			ListAll(T, filter, req)
		case "asst":
			var T []types.OutputAsset
			ListAll(T, filter, req)
		case "meta":
			var T []types.Meta
			ListAll(T, filter, req)
		case "coll":
			var T []types.Coll
			ListAll(T, filter, req)
		case "mint":
			var T []types.Mint
			ListAll(T, filter, req)
		case "dtum":
			var T []types.Datum
			ListAll(T, filter, req)
		case "blck":
			var T []types.Blck
			ListAll(T, filter, req)
		case "witn":
			var T []types.Witn
			ListAll(T, filter, req)
		case "witp":
			var T []types.Witp
			ListAll(T, filter, req)
		case "rdmr":
			var T []types.Redeemer
			ListAll(T, filter, req)
		case "trans":
			var T []types.Trans
			ListAll(T, filter, req)
		case "cip25":
			var T []types.Cip25
			ListAll(T, filter, req)
		case "cip25p":
			var T []types.Cip25p
			ListAll(T, filter, req)
		case "pool":
			var T []types.Pool
			ListAll(T, filter, req)
		case "dele":
			var T []types.Dele
			ListAll(T, filter, req)
		case "reti":
			var T []types.Reti
			ListAll(T, filter, req)
		case "skre":
			var T []types.Skre
			ListAll(T, filter, req)

		case "skde":
			var T []types.Skde
			ListAll(T, filter, req)
		case "block":
			var T []types.BlckN
			ListAll(T, filter, req)
		case "blockf":
			var T []types.BlockfB
			ListAll(T, filter, req)
		default:
			log.Warn().Stack().Msgf("%v not on the list", collection)
		}

		defer req.DisConnect()
	}
}
