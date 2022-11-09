/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/paulmatencio/oura-go/db"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oura-go",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	counter                                                                   types.Counter
	myBdb                                                                     *db.BadgerDB
	badgerRoot, badgerDefault, badgerNft, badgerTran, badgerBlock, badgerPool string
	inputFile, config, mongoUrl, mongoSSL, uri                                string
	mongoDatabase, database, collection, key, value, configFileUsed           string
	projectId, server                                                         string
	timeOut                                                                   time.Duration = 10000
	loglevel, num                                                             int
	limit, fromBlock                                                          int64
	reg                                                                       *bsoncodec.Registry
	network                                                                   string
	chUser, chPassword, chAddrs, chDatabase, cDatabase                        string
	events                                                                    string
	kEvent                                                                    []string
	MktPlaces                                                                 []types.MktPlace
	loglevelHelp                                                                   = " -1 => Trace | 0 => Debug | 1 => Info | 2 => Warning | 3 => Error | 5 => Panic "
	retryRead, retryWrite                                                     bool = true, true

	clientOption = options.ClientOptions{
		// SocketTimeout: &timeOut
		RetryReads:  &retryRead,
		RetryWrites: &retryWrite,
		Registry:    reg,
	}
)

const (
	MaxRetry     int    = 5
	MeConnect    string = "MongoDB connection"
	MeDisconnect string = "MongoDB disconnection"
	MeInsertMany string = "MongoDB Insert many"
	MeInsertOne  string = "MongoDB Insert One"
	PeOperator   string = "Parse operator"
	MeInitReq    string = "Initialise MongoDB connection"
	UeMapJson    string = "Unmarshal map json"
	UeJsonString string = "UMarshal json string to rawJson"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	rootCmd.PersistentFlags().IntVarP(&loglevel, "loglevel", "l", 1, loglevelHelp)
	rootCmd.Flags().StringVarP(&config, "config", "c", "", "Full path of the config file; default $HOME/.oura-go/config.yaml")
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	var (
		configPath string
	)
	if len(config) > 0 {
		log.Printf("Setting Config file to %s", config)
		viper.SetConfigFile(config)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}
		configPath = filepath.Join(home, ".oura-go")
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config") // name of the config file without the extension
		viper.SetConfigType("yaml")
	}
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		configFileUsed = viper.ConfigFileUsed()
		log.Printf("Using config file: %s", configFileUsed)

	} else {
		log.Printf("Error %v  reading config file %s", err, viper.ConfigFileUsed())
		os.Exit(2)
	}
	network = viper.GetString("cardano.network")
	projectId = viper.GetString("blockfrost.projectId")
	server = viper.GetString("blockfrost.server")
	mongoUrl = viper.GetString("mongoDB.url")
	mongoSSL = viper.GetString("mongoDB.ssl")
	fromBlock = viper.GetInt64("mongoDB.fromBlock")
	mongoDatabase = viper.GetString("mongoDB.database")
	chUser = viper.GetString("ch.user")
	chPassword = viper.GetString("ch.password")
	chAddrs = viper.GetString("ch.urls")
	chDatabase = viper.GetString("ch.database")
	events = viper.GetString("event.keep")
	kEvent = keepEvent(events)
	setLoglevel(loglevel)
	SetLogShortFile()
	SetBadgerDatabases()
	// log.Println(chUser, chAddrs, chDatabase, loglevel)
	err = viper.UnmarshalKey("marketPlaces", &MktPlaces)
	reg = types.CardanoType()
	clientOption.Registry = reg
	clientOption.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

}

func setLoglevel(loglevel int) {
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339
	switch loglevel {
	case -1:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case 0:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 5:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}

func SetBadgerDatabases() {
	badgerRoot = viper.GetString("badger.root")
	badgerDefault = viper.GetString("badger.default")
	badgerNft = viper.GetString("badger.databases.nft")
	badgerTran = viper.GetString("badger.databases.trans")
	badgerBlock = viper.GetString("badger.databases.block")
	badgerPool = viper.GetString("badger.databases.pool")
}

func GetBadgerDatabases() {

}

func SetLogShortFile() {

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	return
}

func keepEvent(list string) (events []string) {
	return strings.Split(strings.TrimSpace(list), ",")
}
