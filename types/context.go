package types

import "github.com/jinzhu/copier"

type Context struct {
	BlockHash      string      `bson:"block_hash" json:"block_hash"`
	BlockNumber    int64       `bson:"block_number" json:"block_number"`
	CertificateIdx interface{} `bson:"certificate_idx" json:"certificate_idx"`
	InputIdx       interface{} `bson:"input_idx" json:"input_idx"`
	OutputAddress  interface{} `bson:"output_address" json:"output_address"`
	OutputIdx      interface{} `bson:"output_idx" json:"output_idx"`
	Slot           int64       `bson:"slot" json:"slot"`
	Timestamp      int64       `bson:"timestamp" json:"timestamp"`
	TxHash         string      `bson:"tx_hash" json:"tx_hash"`
	TxIdx          int64       `bson:"tx_idx" json:"tx_idx"`
}

type ContextN struct {
	BlockHash      string      `bson:"block_hash" json:"block_hash"`
	BlockNumber    int64       `bson:"block_number" json:"block_number"`
	CertificateIdx interface{} `bson:"certificate_idx" json:"certificate_idx"`
	Slot           int64       `bson:"slot" json:"slot"`
	Timestamp      int64       `bson:"timestamp" json:"timestamp"`
	TxHash         string      `bson:"tx_hash" json:"tx_hash"`
	TxIdx          int64       `bson:"tx_idx" json:"tx_idx"`
}

func (to *ContextN) CopyFrom(from *Context) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}
