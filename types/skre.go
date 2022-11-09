package types

/*
	Stake registration skre

    {"context":{"block_hash":"87c3ee56724e3a2d83ef2bbe7ff02e8ded6d9626114f6cebdd721e6e094e6421","block_number":7496918,"certificate_idx":0,"input_idx":null,"output_address":null,"output_idx":null,"slot":66225738,"timestamp":1657792029,"tx_hash":"17005f753eab3c9303be1e2c25ececaa01dcf102815c47c9c7fc328077a41fd2","tx_idx":9},
    "fingerprint":"66225738.skre.281964539736104184912796681177198854573",
    "stake_registration":{"credential":{"AddrKeyhash":"44d98060487b76e5f4a51e573d914dd5a2033b089e950c5181367aba"}}}


*/
type Skre struct {
	Context           Context           `json:"context" bson:"context"`
	Fingerprint       string            `json:"fingerprint" bson:"fingerprint"`
	StakeRegistration StakeRegistration `json:"stake_registration,omitempty" bson:"stake_registration,omitempty"`
}

type StakeRegistration struct {
	Credential Credential `json:"credential,omitempty" bson:"credential,omitempty"`
}
