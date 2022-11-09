package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Datum struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
	PlutusDatum PlutusDatum        `bson:"plutus_datum" json:"plutus_datum"`
}

type PlutusDatum struct {
	DatumHash  string      `bson:"datum_hash" json:"datum_hash"`
	PlutusData interface{} `bson:"plutus_data" json:"plutus_data"`
}

func (dtum *Datum) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &dtum)
}

/*
	Datum for JpStore
*/

type PlutusJpstore struct {
	Constructor Constructor `json:"constructor" bson:"constructor"`
	Fields      []Fields    `json:"fields" bson:"fields"`
}

type Int struct {
	NumberDouble string `json:"$numberDouble" bson:"$numberDouble"`
}

type Fields struct {
	Bytes string `json:"bytes,omitempty" bson:"bytes,omitempty"`
	Int   Int    `json:"int,omitempty" bson:"int,omitempty"`
}
type Constructor struct {
	NumberDouble string `json:"$numberDouble" bson:"$numberDouble"`
}

type Adapix struct {
	Int Int `json:"int,omitempty" bson:"int,omitempty"`
}

type PlutusSpacebudz struct {
	Constructor Constructor `json:"constructor"`
	Fields      []struct {
		Constructor Constructor `json:"constructor"`
		Fields      []struct {
			Bytes string `json:"bytes,omitempty"`
			Int   Int    `json:"int,omitempty"`
		} `json:"fields"`
	} `json:"fields"`
}
