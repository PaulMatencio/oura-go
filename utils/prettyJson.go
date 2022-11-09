package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyJson(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func PrintJson(str string) (err error) {
	var j string
	if j, err = PrettyJson(str); err == nil {
		fmt.Printf("%s", j)
	}
	fmt.Printf("\n")
	return
}
