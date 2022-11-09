package types

import (
	"encoding/json"
)

type IStxi interface {
	Unmarshall(b []byte) error
	CopyFrom(from *TxInput)
}

type Stxi struct {
	Context     Context `bson:"context" json:"context"`
	Fingerprint string  `bson:"fingerprint" json:"fingerprint"`
	TxInput     TxInput `bson:"tx_input" json:"tx_input"`
}

type TxInput struct {
	Index    int64    `bson:"index" json:"index"` /* output index */
	TxID     string   `bson:"tx_id" json:"tx_id"`
	InputIdx int64    `bson:"input_idx,omitempty" json:"input_idx,omitempty"`
	Address  string   `bson:"address" json:"address"`
	Amount   int64    `bson:"amount" json:"amount"`
	Assets   []Assets `bson:"assets" json:"assets"`
}

/*
type StxiB struct {
	Context     ContextB `bson:"context" json:"context"`
	Fingerprint string   `bson:"fingerprint" json:"fingerprint"`
	TxInput     TxInputB `bson:"tx_input" json:"tx_input"`
}

type TxInputB struct {
	Index   int64     `bson:"index" json:"index"`
	TxID    string    `bson:"tx_id" json:"tx_id"`
	Address string    `bson:"address" json:"address"`
	Amount  int       `bson:"amount" json:"amount"`
	Assets  []AssetsB `bson:"assets" json:"assets"`
}
*/

func (stxi *Stxi) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &stxi)
}

/*
func (to *StxiB) CopyFrom(from *Stxi) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

*/
