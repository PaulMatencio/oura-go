package types

type Witp struct {
	Context       Context       `bson:"context" json:"context"`
	Fingerprint   string        `bson:"fingerprint" json:"fingerprint"`
	PlutusWitness PlutusWitness `bson:"plutus_witness,omitempty" json:"plutus_witness,omitempty"`
}

type PlutusWitness struct {
	ScriptHash string `bson:"script_hash,omitempty" json:"script_hash,omitempty"`
	ScriptHex  string `bson:"script_hex,omitempty" json:"script_hex,omitempty"`
}
