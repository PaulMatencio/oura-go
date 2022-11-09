package types

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlckN struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `json:"context" bson:"context"`
	Fingerprint string             `json:"fingerprint" bson:"fingerprint"`
	Block       BlockN             `json:"block" bson:"block"`
}

type BlockN struct {
	BodySize                 int           `bson:"body_size" json:"body_size"`
	CborHex                  interface{}   `bson:"cbor_hex" json:"cbor_hex"`
	Epoch                    int           `bson:"epoch" json:"epoch"`
	EpochSlot                int           `bson:"epoch_slot" json:"epoch_slot"`
	Era                      string        `bson:"era" json:"era"`
	Hash                     string        `bson:"hash" json:"hash"`
	NextHash                 string        `json:"next_hash" bson:"next_hash"`
	IssuerVkey               string        `bson:"issuer_vkey" json:"issuer_vkey"`
	Number                   int           `bson:"number" json:"number"`
	PreviousHash             string        `bson:"previous_hash" json:"previous_hash"`
	Slot                     int           `bson:"slot" json:"slot"`
	SlotLeader               string        `json:"slot_leader" bson:"slot_leader"`
	Transactions             []Transaction `bson:"transactions" json:"transactions"`
	TxMeta                   []TxMeta      `bson:"tx_meta" json:"tx_meta"`
	TxCount                  int           `bson:"tx_count" json:"tx_count"`
	Fees                     int64         `bson:"fees" json:"fees"`
	TotalOutput              int64         `bson:"total_output" json:"total_output"`
	InputCount               int           `bson:"input_count" json:"input_count"`
	OutputCount              int           `bson:"output_count" json:"output_count"`
	MintCount                int64         `bson:"mint_count" json:"mint_count"`
	MetaCount                int           `bson:"metadata_count" json:"metadata_count"`
	NativeWitnessesCount     int           `bson:"native_witnesses_count" json:"native_witnesses_count"`
	PlutusDatumCount         int           `bson:"plutus_datum_count" json:"plutus_datum_count"`
	PlutusRdmrCount          int           `bson:"plutus_redeemer_count" json:"plutus_redeemer_count"`
	PlutusWitnessesCount     int           `bson:"plutus_witnesses_count" json:"plutus_witnesses_count"`
	Cip25AssetCount          int           `bson:"cip25_asset_count" json:"cip25_asset_count"`
	Cip20Count               int           `bson:"cip20_count" json:"cip20_count"`
	PoolRegistrationCount    int           `bson:"pool_registration_count" json:"pool_registration_count"`
	PoolRetirementCount      int           `bson:"pool_retirement_count" json:"pool_retirement_count"`
	StakeDelegationCount     int           `bson:"stake_delegation_count" json:"stake_delegation_count"`
	StakeRegistrationCount   int           `bson:"stake_registration_count" json:"stake_registration_count"`
	StakeDeregistrationCount int           `bson:"stake_deregistration_count" json:"stake_deregistration_count"`
	Confirmations            int           `json:"confirmations" bson:"confirmations"`
}

func (blckn *BlckN) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &blckn)
}

/*
func (to *BlckN) CopyFrom(from *BlckB) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

*/

func (blckn *BlckN) CopyFrom(from *Blck) {
	copier.CopyWithOption(&blckn, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

/*  */
