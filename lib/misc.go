package lib

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func LogError(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func Struct2String(v reflect.Value) string {
	var n interface{}
	mapstructure.Decode(v.Interface(), &n)
	return fmt.Sprintf("%s", n)
}

func Hash(h string) (fingerprint string) {
	fp := sha1.Sum([]byte(h))
	return fmt.Sprintf("%x", fp)
}

func GetPolicyID(req *mongodb.MongoDB, flags types.Options, m *types.Meta) (policyId string, err error) {
	var (
		mint types.Mint
		opt1 options.FindOneOptions
	)
	req.Collection = "mint"
	Filter := bson.M{"context.tx_hash": m.Context.TxHash}
	if mint, err = FindOne(&opt1, mint, Filter, req); err == nil {
		policyId = mint.Asset.Policy
		return
	} else {
		//log.Error().Err(err).Msg("Find one")
		log.Error().Err(err).Msgf("Find one collection: %s - filter %v ", req.Collection, Filter)

	}
	return

}

func GetAssetName(req *mongodb.MongoDB, flags types.Options, m *types.Meta) (assetName string, err error) {
	var (
		mint types.Mint
		opt1 options.FindOneOptions
	)
	req.Collection = "mint"
	Filter := bson.M{"context.tx_hash": m.Context.TxHash}
	if mint, err = FindOne(&opt1, mint, Filter, req); err == nil {
		if asst, err1 := hex.DecodeString(mint.Asset.Asset); err == nil {
			assetName = string(asst)
			err = err1
		}
		return
	} else {
		log.Error().Err(err).Msgf("Find one collection: %s - filter %v ", req.Collection, Filter)
	}
	return

}

func PrimitiveDtoMap(in interface{}) (Map map[string]interface{}, err error) {
	var b []byte
	b, err = bson.MarshalExtJSON(in, true, true)
	if err == nil {
		err = json.Unmarshal(b, &Map)
	} else {
		log.Error().Err(err).Msg("unmarshal mongodb primitive to map")
	}
	return
}

/*
	Check if a given event was inserted
*/
func CheckOne(unknown types.Unknown, in []byte, req *mongodb.MongoDB) (fingerprint string, err error) {
	var (
		filter interface{}
		opt1   options.FindOneOptions
	)
	if err = json.Unmarshal(in, &unknown); err == nil {
		fingerprint = unknown.Fingerprint
		filter = bson.M{"fingerprint": fingerprint}
		unknown, err = FindOne(&opt1, unknown, filter, req)
	}
	return
}

func ParseFloat(str string) (float64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return val, nil
	}
	str = strings.Replace(str, ",", "", -1)
	pos := strings.IndexAny(str, "eE")
	if pos < 0 {
		return strconv.ParseFloat(str, 64)
	}

	var baseVal float64
	var expVal int64

	baseStr := str[0:pos]
	baseVal, err = strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return 0, err
	}

	expStr := str[(pos + 1):]
	expVal, err = strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return baseVal * math.Pow10(int(expVal)), nil
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("the provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}
