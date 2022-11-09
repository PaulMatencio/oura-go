package types

/*
	Stake delegation

*/

type Dele struct {
	Context         Context         `json:"context" bson:"context" `
	Fingerprint     string          `json:"fingerprint" bson:"fingerprint"`
	StakeDelegation StakeDelegation `json:"stake_delegation,omitempty" bson:"stake_delegation,omitempty"`
}

type Credential struct {
	AddrKeyhash string `json:"AddrKeyhash,omitempty" bson:"AddrKeyhash,omitempty"` // stake_key

}
type CredentialN struct {
	AddrKeyhash string `json:"AddrKeyhash,omitempty" bson:"AddrKeyhash,omitempty"` // stake_key
	StakeKey    string `json:"stake_key,omitempty" bson:"stake_key,omitempty"`
}

type StakeDelegation struct {
	Credential Credential `json:"credential,omitempty" bson:"credential,omitempty"` // stake key
	PoolHash   string     `json:"pool_hash,omitempty" bson:"pool_hash,omitempty"`   // Pool id
}

type StakeDelegationN struct {
	Context    ContextN    `json:"context" bson:"context" `
	Credential CredentialN `json:"credential,omitempty" bson:"credential,omitempty"` // stake key
	PoolHash   string      `json:"pool_hash,omitempty" bson:"pool_hash,omitempty"`   // Pool id
	PoolId     string      `json:"pool_id,omitempty" bson:"pool_id,omitempty"`
}
