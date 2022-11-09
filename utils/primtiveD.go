package utils

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func PrimitiveDtoStruct[T any](v T, in interface{}) (result T, err error) {
	var b []byte
	b, err = bson.Marshal(&in)
	if err == nil {
		err = bson.Unmarshal(b, &result)
	}
	return
}

func PrimitiveDtoMap(in interface{}) (Map map[string]interface{}, err error) {
	var b []byte
	b, err = bson.MarshalExtJSON(in, true, true)
	if err == nil {
		err = json.Unmarshal(b, &Map)
	}
	return
}
