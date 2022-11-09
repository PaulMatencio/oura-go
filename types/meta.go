package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMeta interface {
	Unmarshall(b []byte) error
}

type Meta struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `bson:"context" json:"context"`
	Fingerprint string             `bson:"fingerprint" json:"fingerprint"`
	Metadata    interface{}        `bson:"metadata" json:"metadata"`

	// Metadata Metadata `bson:"metadata" json:"metadata"`
}

type Metadata struct {
	Label      string      `bson:"label" json:"label"`
	TextScalar string      `bson:"text_scalar,omitempty" json:"text_scalar,omitempty"`
	MapJson    interface{} `bson:"map_json,omitempty" json:"map_json,omitempty"`
}

type MetaScalar struct {
	Label      string `bson:"label" json:"label"`
	TextScalar string `bson:"text_scalar,omitempty" json:"text_scalar,omitempty"`
}

type MapJSON struct {
	Label     string   `bson:"label,omitempty" json:"label,omitempty"`
	ExtraData []string `bson:"extraData,omitempty" json:"extraData,omitempty"`
	Msg       []string `bson:"msg,omitempty" json:"msg,omitempty"`
}

type Meta721 struct {
	Label string `bson:"label" json:"label"`
	// MapJson map[string]map[string]interface{}
	MapJson map[string]interface{}
}

//  crypto dino

type Meta405 struct {
	Label   string     `bson:"label" json:"label"`
	MapJson MapJSON405 `bson:"map_json,omitempty" json:"map_json,omitempty"`
}
type MapJSON405 struct {
	Price  string   `json:"price" bson:"price"`
	Addr   []string `json:"addr" bson:"addr"`
	Asset  string   `json:"asset" bson:"asset"`
	Op     string   `json:"op" bson:"op" `
	Policy string   `json:"policy" bson:"policy"`
}

type NFTAssets struct {
	Asset       string      `bson:"asset" json:"asset"`
	AssetASCII  string      `bson:"asset_ascii" json:"asset_ascii"`
	Policy      string      `bson:"policy" json:"policy"`
	FingerPrint string      `bson:"finger_print,omitempty" json:"finger_print,omitempty"`
	RawJson     interface{} `bson:"raw_json,omitempty" json:"raw_json,omitempty"`
}

type Meta674 struct {
	Label   string `bson:"label" json:"label"`
	MapJson struct {
		ExtraData []string `bson:"extraData,omitempty" json:"extraData,omitempty"`
		Msg       []string `bson:"msg,omitempty" json:"msg,omitempty"`
	} `json:"map_json" bson:"map_json"`
}

type Meta674a struct {
	Label   string `bson:"label" json:"label"`
	MapJson struct {
		Msg string `bson:"msg,omitempty" json:"msg,omitempty"`
	} `json:"map_json" bson:"map_json"`
}

type Meta3322 struct {
	Label   string `json:"label" bson:"label"`
	MapJson struct {
		PubKeyHash string `json:"pubKeyHash" bson:"pubKeyHash" `
		Rewards    struct {
			RefPoolROS     string `json:"refPoolROS" bson:"refPoolROS"`
			RefPoolRewards string `json:"refPoolRewards" bson:"refPoolRewards"`
			SvPoolROS      string `json:"svPoolROS" bson:"svPoolROS"`
			SvPoolRewards  string `json:"svPoolRewards" bson:"svPoolRewards" `
			PctMoreRewards string `json:"pctMoreRewards" bson:"pctMoreRewards"`
		} `json:"rewards" bson:"rewards"`
		UnlockSlotNo    int    `json:"unlockSlotNo" bson:"unlockSlotNo"`
		LockAmount      string `json:"lockAmount" bson:"lockAmount" `
		LockNumEpochs   int    `json:"lockNumEpochs" bson:"lockNumEpochs"`
		LockTimestamp   int64  `json:"lockTimestamp" bson:"lockTimestamp"`
		UnlockEpochNo   int    `json:"unlockEpochNo" bson:"unlockEpochNo"`
		UnlockTimestamp int64  `json:"unlockTimestamp" bson:"unlockTimestamp"`
		CustomData      string `json:"customData" bson:"customData"`
		LockEpochNo     int    `json:"lockEpochNo" bson:"lockEpochNo"`
		PoolID          string `json:"poolId" bson:"poolId"`
	} `json:"map_json" bson:"map_json"`
}

func (tx *Meta) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &tx)
}
