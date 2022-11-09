package lib

import (
	"encoding/json"
	"fmt"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
)

/*
	generic
*/

func InsertOne[TB any](flags types.Options, req *mongodb.MongoDB, b []byte, out *TB) {
	if err := json.Unmarshal(b, out); err == nil {
		if flags.Upload {
			if _, err := req.InsertOne(out); err != nil {
				LogError(err, "insert one")
			}
		}
	} else {
		log.Error().Err(err).Msgf("unmarshall %s", string(b))
	}
}

func ParseLines(flags types.Options, typ string, lines []string, req *mongodb.MongoDB) (nerr int) {
	var (
		collection = typ
	)
	switch collection {
	case "asst":
		var tb types.OutputAsset
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "utxo":
		var tb types.Utxo
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "stxi":
		var tb types.Stxi
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "meta":
		var tb types.Meta
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "tx":
		var tb types.Tx
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "mint":
		var tb types.Mint
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "coll":
		var tb types.Coll
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "dtum":
		var tb types.Datum
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "blck":
		var tb types.Blck
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "rdmr":
		var tb types.Redeemer
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "witn":
		var tb types.Witn
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "witp":
		var tb types.Witp
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "pool":
		var tb types.Pool
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "reti":
		var tb types.Reti
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "skde":
		var tb types.Skde
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "skre":
		var tb types.Skre
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "dele":
		var tb types.Dele
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "cip25":
		var tb types.Cip25
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "cip15":
		var tb types.Cip15
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "scpt":
		var tb types.Scpt
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	case "trans":
		var tb types.Trans
		nerr = InsertManyInterface(flags, collection, lines, req, tb)
	default:
	}
	return
}

/*
	Using generic
*/

func InsertManyGeneric[tb any](flags types.Options, collection string, lines []string, req *mongodb.MongoDB, in tb) int {

	var (
		nerror = 0
		ins    []interface{}
	)
	/*
		prepare the array of document
	*/
	for _, l := range lines {
		if err := json.Unmarshal([]byte(l), &in); err == nil {
			ins = append(ins, in)
		} else {
			LogError(err, fmt.Sprintf("unmarshal %s", l))
			nerror++
		}
	}
	/*
		inserting  the array of documents
	*/
	req.Collection = collection
	if len(ins) > 0 {
		if flags.Upload {
			if result, err := req.InsertMany(ins); err != nil {
				if result != nil {
					log.Error().Err(err).Msgf("insert many to collection:%s - number of docs:%d - number of insert ids:%d", collection, len(ins), len(result.InsertedIDs))
				} else {
					log.Error().Err(err).Msgf("insert many to collection:%s - number of docs:%d - number of insert ids:%d", collection, len(ins), 0)
				}
				nerror++
			}
		} else {
			ListFingerPrint(req.Collection, ins)
		}
	}
	return nerror
}

/*
	using Interface
    It is useful since the generic version has some issues with the array of documents
    ( golang bug ?)
*/

func InsertManyInterface(flags types.Options, collection string, lines []string, req *mongodb.MongoDB, in interface{}) int {

	var (
		nerror = 0
		ins    []interface{}
	)

	/*
	 prepare the array of documents
	*/
	for _, l := range lines {
		if err := json.Unmarshal([]byte(l), &in); err == nil {
			ins = append(ins, in)
		} else {
			log.Error().Err(err).Msgf("unmarshal %s", l)
			nerror++
		}
	}

	/*
		insert the  array of documents
	*/

	req.Collection = collection
	if len(ins) > 0 {
		if flags.Upload {
			if result, err := req.InsertMany(ins); err != nil {
				if result != nil {
					log.Error().Err(err).Msgf("insert many to collection:%s - number of docs:%d - number of insert ids:%d", collection, len(ins), len(result.InsertedIDs))
				} else {
					log.Error().Err(err).Msgf("insert many to collection:%s - number of docs:%d - number of insert ids:%d", collection, len(ins), 0)
				}
				nerror++
			}
		} else {
			ListFingerPrint(collection, ins)

		}
	}
	return nerror
}

func InsertMany(flags types.Options, docs []interface{}, req *mongodb.MongoDB) (nerr int) {

	// req.Collection = collection
	if len(docs) > 0 {
		if flags.Upload {
			if result, err := req.InsertMany(docs); err != nil {
				if result != nil {
					log.Error().Err(err).Msgf("insert many to collection:%s - number of docs:%d - number of insert ids:%d", req.Collection, len(docs), len(result.InsertedIDs))
				} else {
					log.Error().Err(err).Msgf("insert many to collection:%s - number of docs:%d - number of insert ids:%d", req.Collection, len(docs), 0)
				}
				nerr++
			}
		} else {
			ListFingerPrint(req.Collection, docs)
		}
	}
	return

}

/*
used only for testing
*/

func ListFingerPrint(collection string, docs []interface{}) {
	for _, doc := range docs {
		fmt.Println(collection, utils.GetFingerPrint(doc))
	}
}
