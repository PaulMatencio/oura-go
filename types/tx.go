package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITx interface {
	Unmarshall(b []byte) error
	SetContext(tx *Tx) (err error)
}

type Tx struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
	Transaction Transaction        `bson:"transaction" json:"transaction"`
}

type Transaction struct {
	Fee             int64       `bson:"fee" json:"fee"`
	Hash            string      `bson:"hash" json:"hash""`
	InputCount      int         `bson:"input_count" json:"input_count"`
	Inputs          interface{} `bson:"inputs" json:"inputs"`
	Metadata        interface{} `bson:"metadata" json:"metadata"`
	Mint            interface{} `bson:"mint" json:"mint"`
	MintCount       int64       `bson:"mint_count" json:"mint_count"`
	NativeWitnesses interface{} `bson:"native_witnesses" json:"native_witnesses"`
	NetworkID       interface{} `bson:"network_id" json:"network_id"`
	OutputCount     int         `bson:"output_count" json:"output_count"`
	Outputs         interface{} `bson:"outputs" json:"outputs"`
	PlutusData      interface{} `bson:"plutus_data,omitempty" json:"plutus_data,omitempty" `
	PlutusRedeemers interface{} `bson:"plutus_redeemers,omitempty" json:"plutus_redeemers,omitempty"`
	PlutusWitnesses interface{} `bson:"plutus_witnesses,omitempty" json:"plutus_witnesses,omitempty"`
	TotalOutput     int64       `bson:"total_output" json:"total_output"`
	//	TTL                   Cardano     `bson:"ttl" json:"ttl"`
	TTL                   interface{} `bson:"ttl" json:"ttl"`
	ValidityIntervalStart int64       `bson:"validity_interval_start" json:"validity_interval_start"`
	VkeyWitnesses         interface{} `bson:"vkey_witnesses" json:"vkey_witnesses"`
	Withdrawals           interface{} `bson:"withdrawals,omitempty" json:"withdrawals,omitempty"`
}

/*

type TxB struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     ContextB           `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
	Transaction TransactionB       `bson:"transaction" json:"transaction"`
}

type TransactionB struct {
	Fee                   int64       `bson:"fee" json:"fee"`
	Hash                  string      `bson:"hash" json:"hash""`
	InputCount            int         `bson:"input_count" json:"input_count"`
	Inputs                interface{} `bson:"inputs" json:"inputs"`
	Metadata              interface{} `bson:"metadata" json:"metadata"`
	Mint                  interface{} `bson:"mint" json:"mint"`
	MintCount             int64       `bson:"mint_count" json:"mint_count"`
	NativeWitnesses       interface{} `bson:"native_witnesses" json:"native_witnesses"`
	NetworkID             interface{} `bson:"network_id" json:"network_id"`
	OutputCount           int         `bson:"output_count" json:"output_count"`
	Outputs               interface{} `bson:"outputs" json:"outputs"`
	PlutusData            interface{} `bson:"plutus_data" json:"plutus_data" `
	PlutusRedeemers       interface{} `bson:"plutus_redeemers" json:"plutus_redeemers"`
	PlutusWitnesses       interface{} `bson:"plutus_witnesses" json:"plutus_witnesses"`
	TotalOutput           int64       `bson:"total_output" json:"total_output"`
	TTL                   int64       `bson:"ttl" json:"ttl"`
	ValidityIntervalStart int64       `bson:"validity_interval_start" json:"validity_interval_start"`
	VkeyWitnesses         interface{} `bson:"vkey_witnesses" json:"vkey_witnesses"`
}

*/

func (tx *Tx) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &tx)
}

/*
func (tx *TxB) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &tx)
}

func (to *TxB) CopyFrom(from *Tx) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

*/
