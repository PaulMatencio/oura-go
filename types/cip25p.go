package types

import (
	"eagain.net/go/bech32"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/blake2b"
	"hash"
)

type Cip25p struct {
	Cip25Asset  Cip25pAsset `json:"cip25_asset" bson:"cip25_asset"`
	Context     Context     `json:"context" bson:"context"`
	Fingerprint string      `json:"fingerprint" bson:"fingerprint"`
}

type Cip25pAsset struct {
	Asset        string      `json:"asset" bson:"asset"`
	Description  interface{} `json:"description" bson:"description"`
	Image        interface{} `json:"image" bson:"image"`
	MediaType    string      `json:"media_type" bson:"media_type"`
	Name         interface{} `json:"name" bson:"name"`
	Policy       string      `json:"policy" bson:"policy"`
	FingerPrint  string      `json:"fingerprint,omitempty" bson:"fingerprint,omitempty"`
	Twitter      string      `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Website      interface{} `json:"website,omitempty" bson:"website,omitempty"`
	Artist       string      `json:"artist,omitempty" bson:"artist,omitempty"`
	Publisher    interface{} `json:"publisher,omitempty" bson:"publisher,omitempty"`
	Project      interface{} `json:"project,omitempty" bson:"project,omitempty"`
	PolicyScript interface{} `json:"poly_script,omitempty" bson:"poly_script,omitempty"`
	PolicyLink   string      `json:"policy_link,omitempty" bson:"policy_link,omitempty"`
	Copyright    interface{} `json:"copyright,omitempty" bson:"copyright,omitempty"`
	RawJSON      interface{} `json:"raw_json" bson:"raw_json"`
	Version      string      `json:"version" bson:"version"`
}

func (cip25 *Cip25p) SetFingerPrint(asset Cip25pAsset) (fingerPrint string, err error) {
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

func (asset *Cip25pAsset) SetFingerPrint() (fingerPrint string, err error) {
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
