package types

import (
	"eagain.net/go/bech32"
	"encoding/hex"
	"encoding/json"
	// "github.com/jinzhu/copier"
	"golang.org/x/crypto/blake2b"
	"hash"
	"log"
)

type Mint struct {
	// ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context `bson:"context" json:"context"`
	Fingerprint string  `bson:"fingerprint" json:"fingerprint"`
	Asset       Assets  `bson:"mint" json:"mint"`
}

/*
type MintB struct {
	// ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     ContextB `bson:"context" json:"context"`
	Fingerprint string   `bson:"fingerprint" json:"fingerprint"`
	Asset       AssetsB  `bson:"mint" json:"mint"`
}
*/

func (mint *Mint) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &mint)
}

/*
func (mint *MintB) SetFingerPrint() (err error) {
	var assetId []byte
	assetId, err = hex.DecodeString(mint.Asset.Policy + mint.Asset.Asset)
	if err == nil {
		var hash hash.Hash
		hash, err = blake2b.New(20, nil)
		if err != nil {
			log.Println(err)
			return
		}
		hash.Write(assetId)
		mint.Asset.FingerPrint, err = bech32.Encode("asset", hash.Sum(nil))

	}
	if mint.Asset.AssetASCII == "" {
		assetAscii, err := hex.DecodeString(mint.Asset.Asset)
		if err == nil {
			mint.Asset.AssetASCII = string(assetAscii)
		} else {
			log.Println(err)
		}
	}
	return
}

*/

func (mint *Mint) SetFingerPrint() (fingerPrint string, err error) {
	var assetId []byte
	assetId, err = hex.DecodeString(mint.Asset.Policy + mint.Asset.Asset)
	if err == nil {
		var hash hash.Hash
		hash, err = blake2b.New(20, nil)
		if err != nil {
			log.Println(err)
			return
		}
		hash.Write(assetId)
		// mint.Asset.FingerPrint, err = bech32.Encode("asset", hash.Sum(nil))
		return bech32.Encode("asset", hash.Sum(nil))
		/*     */
	}
	/*
		if mint.Asset.AssetASCII == "" {
			assetAscii, err := hex.DecodeString(mint.Asset.Asset)
			if err == nil {
				mint.Asset.AssetASCII = string(assetAscii)
			} else {
				log.Println(err)
			}
		}

	*/
	return

}

/*
func (to *MintB) CopyFrom(from *Mint) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

*/
