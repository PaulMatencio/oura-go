package types

import (
	"eagain.net/go/bech32"
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
	"hash"
	"log"

	// "github.com/jinzhu/copier"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOutputAsset interface {
	Unmarshall(b []byte) error
	SetFingerPrint() (err error)
}

type OutputAsset struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
	Asset       Assets             `bson:"output_asset" json:"output_asset"`
}

type Assets struct {
	Amount      int    `bson:"amount" json:"amount"`
	Asset       string `bson:"asset" json:"asset"`
	AssetASCII  string `bson:"asset_ascii" json:"asset_ascii"`
	Policy      string `bson:"policy" json:"policy"`
	FingerPrint string `bson:"finger_print,omitempty" json:"finger_print,omitempty"`
}

func (asst *OutputAsset) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &asst)
}

func (asst *OutputAsset) SetFingerPrint() (err error) {
	var assetId []byte
	assetId, err = hex.DecodeString(asst.Asset.Policy + asst.Asset.Asset)
	if err == nil {
		var hash hash.Hash
		hash, err = blake2b.New(20, nil)
		if err != nil {
			log.Println(err)
			return
		}
		hash.Write(assetId)
		asst.Asset.FingerPrint, err = bech32.Encode("asset", hash.Sum(nil))
	}
	if asst.Asset.AssetASCII == "" {
		assetAscii, err := hex.DecodeString(asst.Asset.Asset)
		if err == nil {
			asst.Asset.AssetASCII = string(assetAscii)
		} else {
			log.Println(err)
		}
	}

	return
}

func (asst *Assets) SetFingerPrint() (fingerprint string, err error) {
	var assetId []byte
	assetId, err = hex.DecodeString(asst.Policy + asst.Asset)
	if err == nil {
		var hash hash.Hash
		hash, err = blake2b.New(20, nil)
		if err != nil {
			log.Println(err)
			return
		}
		hash.Write(assetId)
		fingerprint, err = bech32.Encode("asset", hash.Sum(nil))
	}
	return
}
