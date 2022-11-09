package lib

import "github.com/paulmatencio/oura-go/types"

func AddressesToMap(val []interface{}) {
	map1 := make(map[string]int64)
	for _, addr := range val {
		switch addr.(type) {
		case types.TxInput:
			addr1 := addr.(types.TxInput)
			map1[addr1.Address] = int64(addr1.Amount)
		}
	}
}
