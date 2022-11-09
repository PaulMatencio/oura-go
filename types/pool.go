package types

import (
	"eagain.net/go/bech32"
	"encoding/hex"
	"github.com/jinzhu/copier"
)

/*
	Pool registration


 Pool Meta
    {"context":{"block_hash":"c338693b27586cf4d3bdb5408dc153968e646fcb3faae6bbf183e5a6df64ac45","block_number":7496919,"certificate_idx":null,"input_idx":null,"output_address":null,"output_idx":null,"slot":66225759,"timestamp":1657792050,"tx_hash":"2ee78b3f0fe96dd85a9736036a5ede8c14caf93e7871fb0d1f17b1626ec1e4ee","tx_idx":9},
"fingerprint":"66225759.meta.217324739661216049004081452023260812839",
"metadata":{"label":"3322","map_json":{"customData":"","lockAmount":"505000000","lockEpochNo":350,"lockNumEpochs":36,"lockTimestamp":1657791929980,"poolId":"pool16u5trhs7ysdh2e8aadk0p8xp52sczl8s29ay6ythv3sxujmdfeu","pubKeyHash":"94e56e5218a56f640cd166832869e00c8260039a1d747678578c534f","rewards":{"pctMoreRewards":"18.36","refPoolROS":"3.6285","refPoolRewards":"9115499","svPoolROS":"4.2921","svPoolRewards":"10799798"},"unlockEpochNo":386,"unlockSlotNo":81388801,"unlockTimestamp":1672955092000}}}
  Pool  registration

{"context":{"block_hash":"66702a9ad9970ecb863de14ae07232d5a03026894ef81c58f78cd7a99cc84ddf","block_number":7499268,"certificate_idx":0,"input_idx":null,"output_address":null,"output_idx":null,"slot":66275448,"timestamp":1657841739,"tx_hash":"b121de5191c92a6af0de64e4403748fb7a5a6344c2c3242344f0ea930538170e","tx_idx":0},
"fingerprint":"66275448.pool.162326797182025811629202575620615872923",
"pool_registration":{"cost":340000000,"margin":0.0,"operator":"9075ebe17e99f4b6d87e01869bc40f37cbe88ffa29f5deb8d61805a7","pledge":250000000000,"pool_metadata":"https://www.doorkstaking.com/doorkMetaV1.json","pool_metadata_hash":"8ff2738ff2207e135021d0e07e2a032f823cb2ab929ba54d730069ad4ab10f8c","pool_owners":["2936d659c601491dd8d6347af0d111af0e5daa7a7f01afd86f3a5c58","3ee92ade2bab1bfaaa9758b61c3f9d3ac23a9602cb44bfc65c7fe81e"],"relays":["node1.acmestaking.com:55444","node2.acmestaking.com:55444"],"reward_account":"e12936d659c601491dd8d6347af0d111af0e5daa7a7f01afd86f3a5c58","vrf_keyhash":"3aec7fae374d05548eeb0fb12682bebe5cb86dcfaed5c68f8f4369f00bfc0433"}}

Find all delegation for a given pool id per epoch
	Find block per epoch
	get timestamp begin_epoch - end_epoch
  	for each pool ( with pool id)
   	  look for delegation where pool.operator && context.epoch == dele.stake_delegation.pool_hash
   	  add
*/

type Pool struct {
	Context          Context          `json:"context" bson:"context"`
	Fingerprint      string           `json:"fingerprint" bson:"fingerprint"`
	PoolRegistration PoolRegistration `json:"pool_registration,omitempty" bson:"pool_registration,omitempty"`
}

type PoolRegistration struct {
	Cost             Cardano  `json:"cost,omitempty" bson:"cost,omitempty"`
	Margin           float64  `json:"margin,omitempty" bson:"margin,omitempty"`
	Operator         string   `json:"operator,omitempty" bson:"operator,omitempty"`
	PoolId           string   `json:"pool_id,omitempty" bson:"pool_id,omitempty"`
	Pledge           Cardano  `json:"pledge,omitempty" bson:"pledge,omitempty"`
	PoolMetadata     string   `json:"pool_metadata,omitempty" bson:"pool_metadata,omitempty"`
	PoolMetadataHash string   `json:"pool_metadata_hash,omitempty" bson:"pool_metadata_hash,omitempty"`
	PoolOwners       []string `json:"pool_owners,omitempty" bson:"pool_owners,omitempty"`
	Relays           []string `json:"relays,omitempty" bson:"relays,omitempty"`
	RewardAccount    string   `json:"reward_account,omitempty" bson:"reward_account,omitempty"`
	VrfKeyhash       string   `json:"vrf_keyhash,omitempty" bson:"vrf_keyhash,omitempty"`
}

type PoolN struct {
	Context             Context               `json:"context" bson:"context"`
	Fingerprint         string                `json:"fingerprint" bson:"fingerprint"`
	PoolRegistration    PoolRegistrationN     `json:"pool_registration,omitempty" bson:"pool_registration,omitempty"`
	PoolRetirement      PoolRetirementN       `json:"pool_retirement,omitempty" bson:"pool_retirement,omitempty"`
	StakeDelegation     []StakeDelegationN    `json:"stake_delegation,omitempty" bson:"stake_delegation,omitempty"`
	StakeRegistration   []StakeRegistration   `json:"stake_registration,omitempty" bson:"stake_registration,omitempty"`
	StakeDeregistration []StakeDeregistration `json:"stake_deregistration,omitempty" bson:"stake_deregistration,omitempty"`
}

type PoolRegistrationN struct {
	Cost             Cardano  `json:"cost,omitempty" bson:"cost,omitempty"`
	Margin           float64  `json:"margin,omitempty" bson:"margin,omitempty"`
	Operator         string   `json:"operator,omitempty" bson:"operator,omitempty"`
	PoolId           string   `json:"pool_id,omitempty" bson:"pool_id,omitempty"`
	Pledge           Cardano  `json:"pledge,omitempty" bson:"pledge,omitempty"`
	PoolMetadata     string   `json:"pool_metadata,omitempty" bson:"pool_metadata,omitempty"`
	PoolMeta         PoolMeta `json:"pool_meta,omitempty" bson:"pool_meta,omitempty"`
	PoolMetadataHash string   `json:"pool_metadata_hash,omitempty" bson:"pool_metadata_hash,omitempty"`
	PoolOwners       []string `json:"pool_owners,omitempty" bson:"pool_owners,omitempty"`
	Relays           []string `json:"relays,omitempty" bson:"relays,omitempty"`
	RewardAccount    string   `json:"reward_account,omitempty" bson:"reward_account,omitempty"`
	VrfKeyhash       string   `json:"vrf_keyhash,omitempty" bson:"vrf_keyhash,omitempty"`
}

type PoolMeta struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Ticker      string `json:"ticker" bson:"ticker"`
	Homepage    string `json:"homepage,omitempty" bson:"homepage,omitempty"`
	Extended    string `json:"extended,omitempty" bson:"extended,omitempty"`
}

func (to *PoolN) CopyFrom(from *Pool) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

func (p *Pool) GetPoolId() (poolId string, err error) {
	pid, err := hex.DecodeString(p.PoolRegistration.Operator)
	if err == nil {
		return bech32.Encode("pool", pid)
	} else {
		return "", err
	}
}

func (p *Pool) SetPoolId() (err error) {
	v, err := hex.DecodeString(p.PoolRegistration.Operator)
	if err == nil {
		poolId, err := bech32.Encode("pool", v)
		if err == nil {
			p.PoolRegistration.PoolId = poolId
		}
	}
	return
}

func (p *Pool) SetPoolOwner() (err error) {
	for k, v := range p.PoolRegistration.PoolOwners {
		v1, err := hex.DecodeString(v)
		if err == nil {
			poolOwner, err := bech32.Encode("stake", v1)
			if err == nil {
				p.PoolRegistration.PoolOwners[k] = poolOwner
			}
		}
	}
	return
}

func (p *PoolN) SetPoolOwner(network string) (err error) {
	var ex = "e1"
	if network != "mainnet" {
		ex = "e0"
	}
	for k, v := range p.PoolRegistration.PoolOwners {
		v1, err := hex.DecodeString(ex + v)
		if err == nil {
			poolOwner, err := bech32.Encode("stake", v1)
			if err == nil {
				p.PoolRegistration.PoolOwners[k] = poolOwner
			}
		}
	}
	return
}

func (p *PoolN) SetRewardAccount() (err error) {
	v, err := hex.DecodeString(p.PoolRegistration.RewardAccount)
	if err == nil {
		rewardAccount, err := bech32.Encode("stake", v)
		if err == nil {
			p.PoolRegistration.RewardAccount = rewardAccount
		}
	}
	return
}

func (p *Pool) SetRewardAccount() (err error) {
	v, err := hex.DecodeString(p.PoolRegistration.RewardAccount)
	if err == nil {
		rewardAccount, err := bech32.Encode("stake", v)
		if err == nil {
			p.PoolRegistration.RewardAccount = rewardAccount
		}
	}
	return
}
