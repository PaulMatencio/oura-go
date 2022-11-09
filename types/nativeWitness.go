package types

// badger db

type Witn struct {
	Context       Context       `bson:"context" json:"context"`
	Fingerprint   string        `bson:"fingerprint" json:"fingerprint"`
	NativeWitness NativeWitness `bson:"native_witness,omitempty" json:"native_witness,omitempty"`
}

type NativeWitness struct {
	PolicyID   string     `bson:"policy_id,omitempty" json:"policy_id,omitempty"`
	ScriptJSON ScriptJSON `bson:"script_json,omitempty" json:"script_json,omitempty"`
}

type ScriptJSON struct {
	Scripts []Scripts `bson:"scripts,omitempty" json:"scripts,omitempty"`
	Type    string    `bson:"type,omitempty" json:"type,omitempty"`
}

type Scripts struct {
	KeyHash string `bson:"keyHash,omitempty" json:"keyHash,omitempty"`
	Type    string `bson:"type,omitempty" json:"type,omitempty"`
}
