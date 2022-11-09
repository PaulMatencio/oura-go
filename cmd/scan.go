/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"bufio"
	"encoding/json"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

// scanCmd represents the scan command
var (
	scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "A brief description of your command",
		Long:  ``,
		Run:   scanFile,
	}
)

type ScanResult struct {
	Blck int `json:"number_blocks"`
	Tx   int `json:"number_transactions"`
	Asst int `json:"number_assets"`
	Oth  int
}

func init() {
	rootCmd.AddCommand(scanCmd)
	initScan(scanCmd)
}

func initScan(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVarP(&num, "input-number", "n", 1, "input number")
	cmd.Flags().BoolVarP(&printIt, "print", "p", false, "print the result")
}

/*
	Scan input lines
*/

func scanFile(cmd *cobra.Command, args []string) {

	var err error
	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf(" logging to file %s - error %v ", logger, err)
	}
	if inputFile != "" {
		if num <= 1 {
			if err = scanLines(inputFile); err != nil {
				LogError(err, "scan lines")
			}
		} else {
			if err = SplitLns(inputFile, num); err != nil {
				LogError(err, "Split multi lines")
			}
		}
	} else {
		log.Warn().Msgf("%v\n", "missing input file")
	}
}

/*
Scan line by line
*/
func scanLines(filename string) error {
	var (
		fp, err = os.Open(filename)
		// j       string
		unknown types.Unknown
		b       []byte
		Total   int
	)
	if err == nil {
		defer fp.Close()
		r := bufio.NewReader(fp)
		text, err := utils.Readline(r)
		for err == nil {
			if len(text) > 0 {
				Total++
				b = []byte(text)
				err = json.Unmarshal(b, &unknown)
				if err == nil {
					if unknown.Fingerprint != "" {
						typ, _ := utils.GetType(unknown.Fingerprint)
						counter.Increment(typ)
					}
				} else {
					log.Error().Msgf("%v", err)
				}

			}
			text, err = utils.Readline(r)
		}

		// print the counter
		log.Info().Msgf("Total number of events %v", Total)
		counter.Print()
	}
	return err
}

/*
    Scan n lines at a time  until end of file
	   - read n lines
	   - scan n lines
*/

func SplitLns(filename string, num int) error {
	var (
		fp, err = os.Open(filename)
		set     types.Events
		stop        = false
		Total   int = 0
		counter types.Counter
	)

	if err == nil {
		// set.Set1 = make(map[string][]string)
		set.Init()
		defer fp.Close()
		r := bufio.NewReader(fp)
		for !stop {
			linea, _ := SplitLns2(r, num)
			if len(linea) == 0 {
				stop = true
			} else {
				// SplitTypes(linea, set)
				set.SplitTypes(linea)
				// process  intermediate
				m := set.Ev
				total := 0
				for k, v := range m {
					// fmt.Printf("\t%s %d\n", k, len(v))
					total += len(v)
					counter.Add(k, len(v))
				}
				// log.Trace().Msgf("Total %d", total)
				Total += total
				// clear the set ( map)
				set.Init()

			}
		}
		counter.Print()
		// log.Info().Msgf("Total number of events %d - Total number of bocks/trans/assets %d/%d/%d", Total, numBlock, numTx, numAsst)
	}
	return err

}

/*
	called by SplitLns
*/

func SplitLns2(r *bufio.Reader, num int) (linea []string, err error) {
	var (
		k    int = 0
		text string
		stop bool = false
	)
	text, err = utils.Readline(r)
	for err == nil && !stop {
		if len(text) > 0 {
			linea = append(linea, text)
			k++
		}
		if k >= num {
			log.Info().Msgf("%d lines are being processed\n", len(linea))
			// fmt.Printf("%d lines are being processed\n", len(linea))
			stop = true
		} else {
			text, err = utils.Readline(r)
		}
	}
	return
}
