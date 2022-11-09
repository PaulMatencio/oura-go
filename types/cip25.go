package types

import (
	"eagain.net/go/bech32"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/blake2b"
	"hash"
)

type Cip25 struct {
	Cip25Asset  Cip25Asset `json:"cip25_asset" bson:"cip25_asset"`
	Context     Context    `json:"context" bson:"context"`
	Fingerprint string     `json:"fingerprint" bson:"fingerprint"`
}

type Cip25Asset struct {
	Asset       string      `json:"asset" bson:"asset"`
	Description interface{} `json:"description" bson:"description"`
	Image       interface{} `json:"image" bson:"image"`
	MediaType   string      `json:"media_type" bson:"media_type"`
	Name        string      `json:"name" bson:"name"`
	Policy      string      `json:"policy" bson:"policy"`
	FingerPrint string      `json:"fingerprint,omitempty" bson:"fingerprint,omitempty"`
	RawJSON     interface{} `json:"raw_json" bson:"raw_json"`
	Version     string      `json:"version" bson:"version"`
}

func (cip25 *Cip25) SetFingerPrint(asset Cip25Asset) (fingerPrint string, err error) {
	var assetId []byte
	assetId, err = hex.DecodeString(asset.Policy + hex.EncodeToString([]byte(asset.Asset)))
	if err == nil {
		var hash hash.Hash
		hash, err = blake2b.New(20, nil)
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}
		hash.Write(assetId)
		fingerPrint, err = bech32.Encode("asset", hash.Sum(nil))
	}
	return
}

func (asset *Cip25Asset) SetFingerPrint() (fingerPrint string, err error) {
	var assetId []byte
	assetId, err = hex.DecodeString(asset.Policy + hex.EncodeToString([]byte(asset.Asset)))
	if err == nil {
		var hash hash.Hash
		hash, err = blake2b.New(20, nil)
		if err != nil {
			log.Error().Msgf("%v", err)
			return
		}
		hash.Write(assetId)
		fingerPrint, err = bech32.Encode("asset", hash.Sum(nil))
	}
	return
}
