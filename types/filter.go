package types

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
)

type Filter struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	Operator string      `json:"operator"`
}

type Filters struct {
	FilterArray []Filter
}

/*
	bson.D{{"context.tx_hash", "fb89828ee1f40cf0c163e0e3d236c067ece2e17ce93d0cf955efbd4c0b45fac9"}}
	$and: [{$or: [{"context.tx_hash": "fb89828ee1f40cf0c163e0e3d236c067ece2e17ce93d0cf955efbd4c0b45fac9"},
	{"context.tx_hash": "6522a4bdd577788c49b3791e929bc02d1d45e2010b2a1c181aa3cd892ff07f4c"}]},{VIP: true}]

    bson.D{{"fingerprint", bson.D{{"$gte", "55813901.tx.274720608624015412428366917628365783646"}}}}
*/

func (f *Filters) ParseOp(input string, ops map[string]string) (filters Filters, err error) {
	var (
		filter Filter
		valid  = "eq,lte,gte,lt,gt,nin,in,"
	)
	fils := strings.Split(input, ":")

	for _, v := range fils {
		fil := strings.Split(v, ",")
		if len(fil) >= 3 {
			if op, ok := ops[fil[1]]; ok {
				filter.Key = fil[0]
				filter.Operator = op
				filter.Value = fil[2]
				if len(fil) > 3 && strings.ToUpper(fil[3]) == "N" {
					if n, err := strconv.Atoi(fil[2]); err == nil {
						filter.Value = n
					}
				}
			} else {
				err = errors.New(fmt.Sprintf("operator %s is not valid. Valid are %s )", fil[1], valid))
			}
		} else {
			err = errors.New(fmt.Sprintf("operator is missing. Valid operator are %s", valid))
		}
		filters.FilterArray = append(filters.FilterArray, filter)

	}
	return
}

func (f *Filters) ValidOp() map[string]string {
	opValid := make(map[string]string)
	opValid["ep"] = "$eq"
	opValid["lt"] = "$lt"
	opValid["lte"] = "$lte"
	opValid["gt"] = "$gt"
	opValid["gte"] = "$gte"
	opValid["ne"] = "$ne"
	opValid["in"] = "$in"
	opValid["nin"] = "$nin"
	opValid[""] = ""
	return opValid
}

func (f *Filters) ValidLogicalOp() map[string]string {
	opValid := make(map[string]string)
	opValid["and"] = "$and"
	opValid["or"] = "$or"
	opValid[""] = ""
	return opValid
}

func (f *Filters) BuildFilter() bson.D {
	/*
	   bson.D{{"title", "The Room"}}).Decode(&result)
	*/
	var fil bson.D
	for _, v := range f.FilterArray {
		if v.Operator == "" {
			fmt.Println(v.Key, v.Value)
			fil = bson.D{{v.Key, v.Value}}
		} else {
			fil = bson.D{{v.Key, bson.D{{v.Operator, v.Value}}}}
		}
	}
	return fil
}
