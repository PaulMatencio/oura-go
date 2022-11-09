package lib

import (
	"github.com/paulmatencio/oura-go/types"
)

func AssetsToMap(assets []types.Assets) map[string]types.Assets {
	map1 := make(map[string]types.Assets)
	for _, asset := range assets {
		key1 := asset.AssetASCII
		map1[key1] = asset
	}
	return map1
}
