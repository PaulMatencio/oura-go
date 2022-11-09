package utils

import (
	"fmt"
	"reflect"
)

func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {

	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // Non-structural return error
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// Traversing structure fields
	// Specify the tagName value as the key in the map; the field value as the value in the map
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}

	return out, nil
}

func GetFingerPrint(doc interface{}) (fingerprint string) {
	v := reflect.ValueOf(doc)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			if key.Kind() == reflect.String {
				if key.Interface() == "fingerprint" {
					fingerprint = fmt.Sprintf("%s", reflect.ValueOf(v.MapIndex(key)).Interface())
					break
				}
			}
		}
	}
	return
}
