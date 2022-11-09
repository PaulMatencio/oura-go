package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coll struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Collateral  Collateral         `bson:"collateral" json:"collateral"`
	Context     Context            `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
}

type Collateral struct {
	Index   int    `bson:"index" json:"index"`
	TxID    string `bson:"tx_id" json:"tx_id"`
	Address string `bson:"address,omitempty" json:"address,omitempty"`
	Amount  int    `bson:"amount,omitempty" json:"amount,omitempty"`
}

/*
type CollB struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Collateral  CollateralB        `bson:"collateral" json:"collateral"`
	Context     ContextB           `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
}

type CollateralB struct {
	Index   int    `bson:"index" json:"index"`
	TxID    string `bson:"tx_id" json:"tx_id"`
	Address string `bson:"address,omitempty" json:"address,omitempty"`
	Amount  int    `bson:"amount,omitempty" json:"amount,omitempty"`
}
func (to *CollB) CopyFrom(from *Coll) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

*/

func (call *Coll) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &call)
}
