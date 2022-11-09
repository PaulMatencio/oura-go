package utils

import (
	"errors"
	"fmt"
	"strings"
)

func GetType(fingerPrint string) (typ string, err error) {
	t := strings.Split(fingerPrint, ".")
	if len(t) >= 1 {
		typ = t[1]
	} else {
		err = errors.New(fmt.Sprintf("invalid finger print %s", fingerPrint))
	}
	return
}
