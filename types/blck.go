package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blck struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
	Block       Block              `bson:"block" json:"block"`
}

type Block struct {
	BodySize     int         `bson:"body_size" json:"body_size"`
	CborHex      interface{} `bson:"cbor_hex" json:"cbor_hex"`
	Epoch        int         `bson:"epoch" json:"epoch"`
	EpochSlot    int         `bson:"epoch_slot" json:"epoch_slot"`
	Era          string      `bson:"era" json:"era"`
	Hash         string      `bson:"hash" json:"hash"`
	IssuerVkey   string      `bson:"issuer_vkey" json:"issuer_vkey"`
	Number       int         `bson:"number" json:"number"`
	PreviousHash string      `bson:"previous_hash" json:"previous_hash"`
	Slot         int         `bson:"slot" json:"slot"`
	Transactions interface{} `bson:"transactions" json:"transactions"`
	TxCount      int         `bson:"tx_count" json:"tx_count"`
}

func (b *Blck) GetEpoch() int {
	return b.Block.Epoch
}
func (b *Blck) GetEpochSlot() int {
	return b.Block.EpochSlot
}
func (b *Blck) GetSlot() int {
	return b.Block.Slot
}

/*

type BlckB struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     ContextB           `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
	Block       BlockB             `bson:"block" json:"block"`
}

type BlockB struct {
	BodySize     int         `bson:"body_size" json:"body_size"`
	CborHex      interface{} `bson:"cbor_hex" json:"cbor_hex"`
	Epoch        int         `bson:"epoch" json:"epoch"`
	EpochSlot    int         `bson:"epoch_slot" json:"epoch_slot"`
	Era          string      `bson:"era" json:"era"`
	Hash         string      `bson:"hash" json:"hash"`
	IssuerVkey   string      `bson:"issuer_vkey" json:"issuer_vkey"`
	Number       int         `bson:"number" json:"number"`
	PreviousHash string      `bson:"previous_hash" json:"previous_hash"`
	Slot         int         `bson:"slot" json:"slot"`
	Transactions interface{} `bson:"transactions" json:"transactions"`
	TxCount      int         `bson:"tx_count" json:"tx_count"`
}

func (to *BlckB) CopyFrom(from *Blck) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

*/

func (blck *Blck) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &blck)
}
