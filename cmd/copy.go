/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

// copyCmd represents the copy command
var (
	copyCmd = &cobra.Command{
		Use:   "copy",
		Short: "copy a document from one collection to another collection",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: copyDocs,
	}
	fromCol, fromDb, toCol, toDb, doc string
	docs                              []string
	Num                               bool
)

func init() {
	rootCmd.AddCommand(copyCmd)
	initCopy(copyCmd)
}

func initCopy(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&fromDb, "from-db", "", "cardano", "database name")
	cmd.Flags().StringVarP(&fromCol, "from-col", "", "", "from collection")
	cmd.Flags().StringVarP(&toDb, "to-db", "", "cardano", "database name")
	cmd.Flags().StringVarP(&toCol, "to-col", "", "", "to collection")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "field key to be used as filter")
	cmd.Flags().StringVarP(&doc, "documents", "", "", "list of documents separated by a comma")
	cmd.Flags().BoolVarP(&Num, "Num", "", false, "type of selection")
}

func copyDocs(cmd *cobra.Command, args []string) {
	if fromCol == "" {
		log.Error().Msg("Missing from collection")
		return
	}
	if toCol == "" {
		log.Error().Msg("Missing to collection")
		return
	}
	if doc == "" {
		log.Error().Msg("Missing to collection")
		return
	} else {
		docs = strings.Split(doc, ",")
	}
	if toDb == "" {
		log.Warn().Msg("Missing the target database. Source database will be used as target")
		toDb = fromDb
	}

	if req, err := InitCopy(mongoUrl); err == nil {
		req.Documents = docs
		req.Num = Num
		req.FromDatabase = fromDb
		req.FromCollection = fromCol
		req.ToDatabase = toDb
		req.ToCollection = toCol
		result := req.Copy(filter)
		for _, v := range result {
			fmt.Println(v.InsertedID)
		}
	}

}
