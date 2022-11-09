package types

import (
	"eagain.net/go/bech32"
	"encoding/hex"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/blake2b"
	"hash"
	"log"
)

type IUtxo interface {
	Unmarshall(b []byte) error
	SetContext(utxo *Utxo) (err error)
}

type Utxo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `json:"context" bson:"context"`
	Fingerprint string             `json:"fingerprint" bson:"fingerprint"`
	TxOutput    TxOutput           `json:"tx_output" bson:"tx_output"`
}

type TxOutput struct {
	Address   string      `json:"address" bson:"address"`
	Amount    int64       `json:"amount" bson:"amount"`
	Assets    []Assets    `json:"assets" bson:"assets"`
	DatumHash interface{} `json:"datum_hash" bson:"datum_hash"`
	OutputIdx int64       `json:"output_idx,omitempty" bson:"output_idx,omitempty"`
}

func (utxo *Utxo) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &utxo)
}

func (utxo *Utxo) SetFingerPrint(asset Assets) (fingerPrint string, err error) {
	var assetId []byte
	assetId, err = hex.DecodeString(asset.Policy + asset.Asset)
	if err == nil {
		var hash hash.Hash
		hash, err = blake2b.New(20, nil)
		if err != nil {
			log.Println(err)
			return
		}
		hash.Write(assetId)
		fingerPrint, err = bech32.Encode("asset", hash.Sum(nil))
	}
	return
}
