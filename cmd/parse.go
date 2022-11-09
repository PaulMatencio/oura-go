package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	//"log"
	"os"
	//log "go.uber.org/zap"
	"github.com/rs/zerolog/log"
)

type pBlock struct {
	Number    int64     `json:"block"`
	Hash      string    `json:"hash"`
	Slot      int64     `json:"slot"`
	TxCount   int       `json:"tx_count"`
	BodySize  int       `json:"body_size"`
	Timestamp time.Time `json:"time"`
}

var (
	parseCmd = &cobra.Command{
		Use:   "parse",
		Short: "A brief description of your command",
		Long:  ``,
		Run:   parseFile,
	}
	parseBlockCmd = &cobra.Command{
		Use:   "parseBlock",
		Short: "A brief description of your command",
		Long:  ``,
		Run:   parseBlocks,
	}
	ptime time.Time
	check bool
)

func (pb *pBlock) Println(t time.Duration) (err error) {
	bl, err := json.Marshal(&pb)
	fmt.Printf("%s %v sec\n", string(bl), t.Seconds())
	return
}

func init() {
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(parseBlockCmd)
	initParse(parseCmd)
	initParseBlock(parseBlockCmd)
}

func initParse(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVarP(&num, "input-number", "n", 1, "input number")
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "database", "d", "cardano", "database name")
	cmd.Flags().BoolVarP(&concurrent, "concurrent", "C", false, "Concurrent upload")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger full file path")
	cmd.Flags().BoolVarP(&upload, "upload", "U", false, "upload to mongodb ")
	cmd.Flags().BoolVarP(&check, "check", "", false, "check if event are inserted")
}

func initParseBlock(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	cmd.Flags().IntVarP(&num, "input-number", "n", 1, "input number")
	cmd.Flags().StringVarP(&logger, "logger", "", "", "logger full file path")
}

func parseBlocks(cmd *cobra.Command, args []string) {

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf("logging to file %s - error %v ", logger, err)
	}
	log.Info().Msgf("Upload %v", upload)
	fmt.Printf("Upload %v", upload)

	if inputFile != "" {
		var fp, err = os.Open(inputFile)
		if err == nil {
			defer fp.Close()
			r := bufio.NewReader(fp)
			line, err := utils.Readline(r)
			for err == nil {
				if len(line) > 0 {
					parseBlock(line)
				}
				line, err = utils.Readline(r)
			}
		}
	} else {
		LogError(errors.New("missing input file"), "")
	}
}

func parseBlock(line string) {

	var (
		unknown types.Unknown
		blck    types.Blck
		typ     string
		b       = []byte(line)
		block   pBlock
	)

	if err = json.Unmarshal(b, &unknown); err == nil {
		if typ, err = utils.GetType(unknown.Fingerprint); err == nil {
			if typ == "blck" {
				if err = json.Unmarshal(b, &blck); err == nil {
					block.Number = blck.Context.BlockNumber
					block.Hash = blck.Context.BlockHash
					block.Slot = blck.Context.Slot
					block.TxCount = blck.Block.TxCount
					block.BodySize = blck.Block.BodySize
					block.Timestamp = time.Unix(blck.Context.Timestamp, 0)
					block.Println(block.Timestamp.Sub(ptime))
					ptime = block.Timestamp
				} else {
					LogError(err, fmt.Sprintf("unmarshal blck type %s", line))
				}
			}
		}
	} else {
		LogError(err, fmt.Sprintf("unmarshal unknown type %s", line))
	}
}

/*
Parse an input file
*/
func parseFile(cmd *cobra.Command, args []string) {
	var (
		client *mongo.Client
		err    error
		option = &clientOption
		req    = mongodb.MongoDB{
			Option:     option,
			Database:   database,
			Collection: collection,
		}
		start = time.Now()
	)
	if log.Logger, err = SetLogFile(logger); err != nil {
		LogError(err, fmt.Sprintf("set log to file %s", logger))
	}
	if concurrent {
		utils.SetCPU("50%")
	}
	flags := SetFlags()
	if inputFile != "" {
		log.Info().Msgf("Processing file:%v", inputFile)
		uri = "mongodb://" + mongoUrl + "/?ssl=" + mongoSSL
		req.Uri = uri
		if client, err = req.Connect(); err != nil {
			LogError(err, "Connect req")
			if upload {
				return
			}
		}
		defer req.DisConnect()
		req.Client = client
		if num <= 1 || check {
			err = ParseLines(flags, inputFile, &req)
		} else {
			if !concurrent {
				err = ParseLinesSeq(flags, inputFile, num, &req)
			} else {
				err = ParseLinesCon(flags, inputFile, num, &req)
			}
		}
		if err != nil {
			LogError(err, "Parse multiple lines")
		}
	} else {
		log.Warn().Msg("Missing input file")
	}

	fmt.Printf("Total elapsed time %v\n", time.Since(start))
}

/*
	parse one at a time
    Just for testing ( it is very slow )
*/

func ParseLines(flags types.Options, filename string, req *mongodb.MongoDB) (err error) {

	var fp *os.File
	if fp, err = os.Open(filename); err == nil {

		defer fp.Close()
		r := bufio.NewReader(fp)
		text, err := utils.Readline(r)
		for err == nil {
			if len(text) > 0 {
				if check {
					err = CheckLine(text, req)
				} else {
					err = ParseLine(flags, text, req)
				}
			}
			text, err = utils.Readline(r)
		}
		counter.Print()
	}
	return
}

/*
	Parse multiple ( num) lines at a time
    Split the lines into group events
    For  every group of events
       if the  group ( type)  is selected  (check the yaml config file)
          call insertLines to insert the group of events
    End
*/

func ParseLinesSeq(flags types.Options, filename string, num int, req *mongodb.MongoDB) error {
	var (
		fp, err                 = os.Open(filename)
		set                     types.Events
		stop                    = false
		inserted, total, nerror int
	)
	if err == nil {
		set.Init()
		set.Keep(kEvent)
		// set.Selected(kEvent)
		fmt.Println(set.Kp)
		defer fp.Close()
		r := bufio.NewReader(fp)
		for !stop {
			linea, _ := SplitLns2(r, num)
			if len(linea) == 0 {
				stop = true
			} else {
				set.SplitTypes(linea)
				m := set.Ev
				for k, v := range m {
					if set.Has(k) {
						req.Collection = k
						nerror += lib.ParseLines(flags, k, v, req)
						inserted += len(v)
					}
					total += len(v)
					counter.Add(k, len(v))
				}
			}
			set.Init()
		}
	}
	counter.Print()
	fmt.Printf("Total number:%d - Total number inserted %d  -Total number of errors: %d\n", total, inserted, nerror)
	return err
}

/*
Parse multiple ( num) lines at a time
Split the lines into group events
For  every group of events

	if the  group ( type)  is selected  (check the yaml config file)
	call goroutine to insertLines to insert the group of events

End
*/

func ParseLinesCon(flags types.Options, filename string, num int, req *mongodb.MongoDB) error {

	var (
		fp, err                = os.Open(filename)
		set                    types.Events
		stop                   = false
		Inserted, Total, Error int
	)
	if err == nil {
		set.Init()
		set.Keep(kEvent)
		// set.Selected(kEvent)
		defer fp.Close()
		r := bufio.NewReader(fp)
		for !stop {
			linea, _ := SplitLns2(r, num)
			if len(linea) == 0 {
				stop = true
			} else {
				ntotal, ninserted, nerror := parseArrayCon(flags, req, linea, set)
				set.Init()
				Total += ntotal
				Inserted += ninserted
				Error += nerror
			}
			// set.Init()
		}
	}
	counter.Print()
	fmt.Printf("Total number:%d - Total number inserted %d  -Total number of errors: %d\n", Total, Inserted, Error)
	return err
}

func parseArrayCon(flags types.Options, req *mongodb.MongoDB, linea []string, set types.Events) (total int, inserted int, nerror int) {

	type Response struct {
		Collection string
		Total      int
		Error      int
		Inserted   int
	}
	var (
		ch = make(chan *Response)
	)

	set.SplitTypes(linea)

	m := set.Ev
	request := 0
	for k, v := range m {
		if set.Has(k) {
			request += 1
			go func(k string, v []string, req *mongodb.MongoDB) {
				var resp Response
				req.Collection = k
				resp.Total = len(v)
				resp.Error = lib.ParseLines(flags, k, v, req)
				resp.Inserted = len(v)
				resp.Collection = k
				ch <- &resp
			}(k, v, req)
		} else {
			counter.Add(k, len(v))
			total += len(v)
		}
	}
	receive := 0

	for {
		select {
		case rec := <-ch:
			receive++
			inserted += rec.Inserted
			nerror += rec.Error
			total += rec.Total
			counter.Add(rec.Collection, rec.Total)
			if receive == request {
				// fmt.Printf("Total number:%d - Total number inserted %d  -Total number of errors: %d\n", total, inserted, nerror)
				return
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Printf(".")

		}
	}
}

/*
	called by parseLines
    parse  the given  input line
*/

func ParseLine(flags types.Options, line string, req *mongodb.MongoDB) (err error) {
	var (
		unknown types.Unknown
		typ     string
		set     types.Events
	)
	set.Keep(kEvent)
	b := []byte(line)
	err = json.Unmarshal(b, &unknown)
	if err == nil {
		if typ, err = insertLine(flags, unknown.Fingerprint, b, req, set); err != nil {
			LogError(err, fmt.Sprintf("Parse type: %s", typ))
		} else {
			counter.Increment(typ)
		}
	} else {
		log.Warn().Msgf("%s %v", line, err)
	}
	return
}

/*
Check if the events in a given  input file are all inserted
*/

func CheckLine(line string, req *mongodb.MongoDB) (err error) {
	var (
		unknown     types.Unknown
		fingerprint string
		collection  string
		set         types.Events
	)
	set.Keep(kEvent)
	b := []byte(line)
	if err = json.Unmarshal(b, &unknown); err == nil {
		if collection, err = utils.GetType(unknown.Fingerprint); err == nil {
			kp := set.Kp
			if _, ok := kp[collection]; ok {
				req.Collection = collection
				log.Trace().Msgf("find fingerprint %s in collection %s", unknown.Fingerprint, collection)
				if fingerprint, err = checkOne(b, req, unknown); err == nil {
					counter.Increment(collection)
					log.Trace().Msgf("found fingerprint %s in collection %s", fingerprint, collection)
				} else {
					LogError(err, fmt.Sprintf("find fingerprint %s", fingerprint))
				}
			}
		}
	}
	return
}

/*
	called by parseLine
    insert  events  to its corresponding collection  ( event)
*/

func insertLine(flags types.Options, fingerPrint string, b []byte, req *mongodb.MongoDB, events types.Events) (typ string, err error) {

	if typ, err = utils.GetType(fingerPrint); err != nil {
		return
	}
	collection := typ
	kp := events.Kp
	if _, ok := kp[collection]; !ok {
		return
	}
	req.Collection = typ
	switch collection {
	case "asst":
		var tb types.OutputAsset
		lib.InsertOne(flags, req, b, &tb)
	case "stxi":
		var tb types.Stxi
		lib.InsertOne(flags, req, b, &tb)
	case "utxo":
		var tb types.Utxo
		lib.InsertOne(flags, req, b, &tb)
	case "meta":
		var tb types.Meta
		lib.InsertOne(flags, req, b, &tb)
	case "tx":
		var tb types.Tx
		lib.InsertOne(flags, req, b, &tb)
	case "dtum":
		var tb types.Datum
		lib.InsertOne(flags, req, b, &tb)
	case "mint":
		var tb types.Mint
		lib.InsertOne(flags, req, b, &tb)
	case "coll":
		var tb types.Coll
		lib.InsertOne(flags, req, b, &tb)
	case "blck":
		var tb types.Blck
		lib.InsertOne(flags, req, b, &tb)
	case "rdmr":
		var tb types.Redeemer
		lib.InsertOne(flags, req, b, &tb)
	case "witp":
		var tb types.Witp
		lib.InsertOne(flags, req, b, &tb)
	case "witn":
		var tb types.Witn
		lib.InsertOne(flags, req, b, &tb)
	case "pool":
		var tb types.Pool
		lib.InsertOne(flags, req, b, &tb)
	case "reti":
		var tb types.Reti
		lib.InsertOne(flags, req, b, &tb)
	case "skde":
		var tb types.Skde
		lib.InsertOne(flags, req, b, &tb)
	case "skre":
		var tb types.Skre
		lib.InsertOne(flags, req, b, &tb)
	case "dele":
		var tb types.Dele
		lib.InsertOne(flags, req, b, &tb)
	case "cip25":
		var tb types.Cip25
		lib.InsertOne(flags, req, b, &tb)
	case "cip15":
		var tb types.Cip15
		lib.InsertOne(flags, req, b, &tb)
	case "scpt":
		var tb types.Scpt
		lib.InsertOne(flags, req, b, &tb)
	case "trans":
		var tb types.Trans
		lib.InsertOne(flags, req, b, &tb)
	default:
	}
	return
}

/*
	Insert one doc at a time
*/

func insertOne[TB any](req *mongodb.MongoDB, b []byte, out *TB) {
	if err := json.Unmarshal(b, out); err == nil {
		if upload {
			if _, err := req.InsertOne(out); err != nil {
				LogError(err, "insert one")
			}
		}
	} else {
		log.Error().Err(err).Msgf("unmarshall %s", string(b))
	}
}

/*
Check if a given event was inserted
*/

func checkOne(in []byte, req *mongodb.MongoDB, unknown types.Unknown) (fingerprint string, err error) {
	var (
		filter interface{}
		opt1   options.FindOneOptions
	)
	if err = json.Unmarshal(in, &unknown); err == nil {
		fingerprint = unknown.Fingerprint
		filter = bson.M{"fingerprint": fingerprint}
		unknown, err = lib.FindOne(&opt1, unknown, filter, req)
	}
	return
}
