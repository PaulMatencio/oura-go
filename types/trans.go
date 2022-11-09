package types

import (
	"encoding/json"
	"github.com/paulmatencio/oura-go/utils"
)

type Trans struct {
	Context        Context          `json:"context" bson:"context"`
	Fingerprint    string           `json:"fingerprint" bson:"fingerprint"`
	TxMeta         TxMeta           `json:"tx_meta" bson:"tx_meta"`
	Transaction    Transaction      `json:"transaction" bson:"transaction"`
	TxInput        []TxInput        `json:"utxo_input" bson:"utxo_input"`
	TxOutput       []TxOutput       `json:"utxo_output" bson:"utxo_output"`
	MintAsset      []Assets         `json:"mint_asset,omitempty" bson:"mint_asset,omitempty"`
	Metadata       []interface{}    `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Meta674        []Meta674        `json:"meta_674,omitempty" bson:"meta_674,omitempty"`
	Meta3322       []Meta3322       `json:"meta_3322,omitempty" bson:"meta_3322,omitempty"`
	MetaScalar     []MetaScalar     `json:"meta_scalar,omitempty" bson:"meta_scalar,omitempty"`
	Collateral     Collateral       `json:"collateral,omitempty" bson:"collateral,omitempty"`
	PlutusWitness  PlutusWitness    `json:"plutus_witness,omitempty" bson:"plutus_witness,omitempty"`
	NativeWitness  NativeWitness    `json:"native_witness,omitempty" bson:"native_witness,omitempty"`
	PlutusRedeemer []PlutusRedeemer `json:"plutus_redeemer,omitempty" bson:"plutus_redeemer,omitempty"`
	PlutusDatum    []PlutusDatum    `json:"plutus_datum,omitempty" bson:"plutus_datum,omitempty"`
	Cip25pAsset    []Cip25pAsset    `json:"cip25_asset,omitempty" bson:"cip25_asset,omitempty"`
	// TransBase
	PoolRegistration    PoolRegistration    `json:"pool_registration,omitempty" bson:"pool_registration,omitempty"`
	PoolRetirement      PoolRetirement      `json:"pool_retirement,omitempty" bson:"pool_retirement,omitempty"`
	StakeRegistration   StakeRegistration   `json:"stake_registration,omitempty" bson:"stake_registration,omitempty"`
	StakeDeregistration StakeDeregistration `json:"stake_deregistration,omitempty" bson:"stake_deregistration,omitempty"`
	StakeDelegation     StakeDelegation     `json:"stake_delegation,omitempty" bson:"stake_delegation,omitempty"`
}

type TxMeta struct {
	InputCount               int    `bson:"input_count" json:"input_count"`
	OutputCount              int    `bson:"output_count" json:"output_count"`
	MintCount                int64  `bson:"mint_count" json:"mint_count"`
	CollCount                int    `bson:"coll_count" json:"coll_count"`
	MetaCount                int    `bson:"metadata_count" json:"metadata_count"`
	MetaLabel                string `bson:"meta_label" json:"meta_label"`
	NativeWitnessesCount     int    `bson:"native_witnesses_count" json:"native_witnesses_count"`
	PlutusDatumCount         int    `bson:"plutus_datum_count" json:"plutus_datum_count"`
	PlutusRdmrCount          int    `bson:"plutus_redeemer_count" json:"plutus_redeemer_count"`
	PlutusWitnessesCount     int    `bson:"plutus_witnesses_count" json:"plutus_witnesses_count"`
	Cip25AssetCount          int    `bson:"cip25_asset_count" json:"cip25_asset_count"`
	PoolRegistrationCount    int    `bson:"pool_registration_count" json:"pool_registration_count"`
	PoolRetirementCount      int    `bson:"pool_retirement_count" json:"pool_retirement_count"`
	StakeDelegationCount     int    `bson:"stake_delegation_count" json:"stake_delegation_count"`
	StakeRegistrationCount   int    `bson:"stake_registration_count" json:"stake_registration_count"`
	StakeDeregistrationCount int    `bson:"stake_deregistration_count" json:"stake_deregistration_count"`
}

func (tx *Trans) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &tx)
}

func (s *Trans) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := utils.SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
