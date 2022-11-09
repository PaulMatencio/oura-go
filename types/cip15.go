package types

type Cip15 struct {
	Cip15Asset  Cip15Asset `json:"cip15_asset" bson:"cip15_asset"`
	Context     Context    `json:"context" bson:"context"`
	Fingerprint string     `json:"fingerprint" bson:"fingerprint"`
}
type RawJSON struct {
	Num1 string `json:"1" bson:"1"`
	Num2 string `json:"2" bson:"2"`
	Num3 string `json:"3" bson:"3"`
	Num4 int    `json:"4" bson:"4"`
}
type Cip15Asset struct {
	Nonce         int     `json:"nonce" bson:"nonce"`
	RawJSON       RawJSON `json:"raw_json" bson:"raw_json"`
	RewardAddress string  `json:"reward_address" bson:"reward_address"`
	StakePub      string  `json:"stake_pub" bson:"stake_pub"`
	VotingKey     string  `json:"voting_key" bson:"voting_key"`
}
