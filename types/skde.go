package types

/*
	Stake deregistration

     {"context":{"block_hash":"c9499fb9e9c8b690ca3ea52c07a1a5d5d1602eec61bd6be95e91f111edab90cd",
     "block_number":7496854,"certificate_idx":0,"input_idx":null,"output_address":null,"output_idx":null,"slot":66224503,"timestamp":1657790794,"tx_hash":"fdb373128edd4ad7924b1edba3394f9fad5b0afe680d4a7db69645ea436e2891","tx_idx":17},
     "fingerprint":"66224503.skde.291828108761531886258053653292847118406",
     "stake_deregistration":{"credential":{"AddrKeyhash":"f9d5a371f44bee8c65be315c6f80f2cfc4dbb48745137cd464a2403e"}}}

*/

type Skde struct {
	Context             Context             `json:"context" bson:"context"`
	Fingerprint         string              `json:"fingerprint" bson:"fingerprint"`
	StakeDeregistration StakeDeregistration `json:"stake_deregistration,omitempty" bson:"stake_deregistration,omitempty"`
}

type StakeDeregistration struct {
	Credential Credential `json:"credential,omitempty" bson:"credential,omitempty"`
}
