package cmd

import (
	"github.com/paulmatencio/oura-go/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	listDBCmd = &cobra.Command{
		Use:   "listDB",
		Short: "list a given badgerDB namespace",
		Long:  ``,
		Run:   listDB,
	}
	delDBCmd = &cobra.Command{
		Use:   "delDB",
		Short: "delete a given badgerDB key",
		Long:  ``,
		Run:   delDB,
	}

	badgerDB *db.BadgerDB
	prefix   string
	err      error
	ns       string
)

func init() {
	rootCmd.AddCommand(listDBCmd)
	rootCmd.AddCommand(delDBCmd)
	initDBList(listDBCmd)
	initDBDelete(delDBCmd)
}

func initDBList(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "", "badger database - check the config file ")
	cmd.Flags().StringVarP(&ns, "name-space", "N", "", "badgerDB name space")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "badgerDB key prefix")
}

func initDBDelete(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "", "local database directory. If no a full path, it with be prefix")
	cmd.Flags().StringVarP(&ns, "name-space", "N", "", "badgerDB name space")
	cmd.Flags().StringVarP(&key, "key", "k", "", "badgerDB key")
}

func listDB(cmd *cobra.Command, args []string) {

	var badgerDir string
	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}
	if ns == "" {
		log.Error().Msg("the name space is missing")
		return
	}

	if dataDir == "" {
		badgerDir = GetBadgerPath("")
	} else {
		badgerDir = GetBadgerPath(dataDir)
	}
	if badgerDir == "" {
		log.Error().Msgf("Badger directory is missing")
		return
	} else {
		log.Info().Msgf("Badger directory %s", badgerDir)
		if badgerDB, err = OpenBdb(badgerDir); err != nil {
			log.Error().Msgf("Error opening badger db %v", err)
			return
		}
	}

	err := badgerDB.List([]byte(ns), []byte(prefix))
	if err != nil {
		log.Error().Msgf("Error listing badger db %v", err)
	}

}

func delDB(cmd *cobra.Command, args []string) {

	var badgerDir string
	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}

	if dataDir == "" {
		badgerDir = GetBadgerPath("")
	} else {
		badgerDir = GetBadgerPath(dataDir)
	}
	if badgerDir == "" {
		log.Error().Msgf("Badger directory is missing")
		return
	}
	if badgerDB, err = OpenBdb(badgerDir); err != nil {
		log.Error().Msgf("Error opening badger db %v", err)
		return
	}
	if ns == "" {
		log.Error().Msg("the name space is missing")
		return
	}

	if key == "" {
		log.Error().Msg("the key is missing")
		return
	}

	err := badgerDB.Delete([]byte(ns), []byte(key))
	if err != nil {
		log.Error().Msgf("Error %v deleting key  %s badger db", err, key)
	}

}
