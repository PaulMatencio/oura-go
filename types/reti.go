package types

/*
	Reti
    Pool retirement
*/

type Reti struct {
	Context        Context        `json:"context" bson:"context"`
	Fingerprint    string         `json:"fingerprint" bson:"fingerprint" `
	PoolRetirement PoolRetirement `json:"pool_retirement,omitempty" bson:"pool_retirement,omitempty"`
}

type PoolRetirement struct {
	Epoch int    `json:"epoch,omitempty" bson:"epoch,omitempty"`
	Pool  string `json:"pool,omitempty" bson:"pool,omitempty"`
}

type PoolRetirementN struct {
	Context ContextN `json:"context" bson:"context"`
	Epoch   int      `json:"epoch,omitempty" bson:"epoch,omitempty"`
	Pool    string   `json:"pool,omitempty" bson:"pool,omitempty"`
	PoolId  string   `json:"pool_id,omitempty" bson:"pool_id,omitempty"`
}
