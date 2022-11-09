package types

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
	//"time"
)

type ChTrans struct {
	Table string
	Batch driver.Batch
	Trans []Trans
}

type Trans struct {
	TxHash      uint64    `json:"tx_hash"`
	BlockNumber uint64    `json:"block_number"`
	Fees        uint64    `json:"fees"`
	InputCount  uint32    `json:"input_count"`
	OutputCount uint32    `json:"output_count"`
	MintCount   uint64    `json:"mint_count"`
	MetaCount   uint32    `json:"metadata_count"`
	MetaLabel   string    `json:"metadata_label"`
	MetaMsg     string    `json:"metadata_msg"`
	Stxis       []Utx     `json:"tx_Input"`
	Utxos       []Utx     `json:"tx_output"`
	TimeStamp   time.Time `json:"datetime,omitempty"`
}

type Utx struct {
	Address string `bson:"address" json:"address"`
	Amount  int    `bson:"amount" json:"amount"`
	// Assets  []Assets `bson:"assets" json:"assets"`
}
