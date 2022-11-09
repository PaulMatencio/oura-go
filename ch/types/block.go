package types

import (
	"context"
	"errors"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/rs/zerolog/log"
	"time"
	//"time"
)

type ChBlock struct {
	Table  string
	Batch  driver.Batch
	Blocks []Block
}

const CHBlock = `CREATE TABLE IF NOT EXISTS chblock (
        BlockNumber              UInt64,  
		Slot                     UInt64,
		BodySize                 UInt32,
		Epoch                    UInt32 ,
		EpochSlot                UInt32,
		Era                      String, 
		Hash                     String, 
		IssuerVkey               String,
		SlotLeader               String, 
		TxCount                  UInt32,
		Fees                     UInt64,        
		TotalOutput              UInt64,          
		InputCount               UInt32,          
		OutputCount              UInt32,           
		MintCount                UInt64,        
		MetaCount                UInt32,            
		NativeWitnessesCount     UInt32,           
		PlutusDatumCount         UInt32,           
		PlutusRdmrCount          UInt32,         
		PlutusWitnessesCount     UInt32,        
		Cip25AssetCount          UInt32,            
		Cip20Count               UInt32,              
		PoolRegistrationCount    UInt32,        
		PoolRetirementCount      UInt32,           
		StakeDelegationCount     UInt32,    
		StakeRegistrationCount   UInt32,       
		StakeDeregistrationCount UInt32,      
		Confirmations            UInt32,
		TimeStamp                DateTime                           
	) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(TimeStamp)
	ORDER BY (TimeStamp, BlockNumber);`

const CHBlockEpoch = `CREATE MATERIALIZED VIEW chblock_epoch
       ENGINE = AggregatingMergeTree()
       ORDER BY (SlotLeader,Epoch)
       AS SELECT
         SlotLeader,
         Epoch, 
         countState() block_number,
         sumState(Fees) AS sum_transaction_fees,
         sumState(TxCount) AS sum_transaction_count,
         sumState(TotalOutput) AS sum_total_output,
         sumState(BodySize) AS sum_body_size,
         sumState(PoolRegistrationCount) AS sum_pool_registration,
         sumState(PoolRetirementCount) AS sum_pool_retirement,
         sumState(StakeDelegationCount) AS sum_stake_delegation,  
         sumState(StakeRegistrationCount) AS sum_stake_registration,  
         sumState(StakeDeregistrationCount) AS sum_stake_deregistration  
       FROM  chblock
       GROUP BY Epoch, SlotLeader;`

const CHBlockDay = `CREATE MATERIALIZED VIEW cardano.chblock_day
       ENGINE = AggregatingMergeTree()
       ORDER BY (year,month,day)
       AS SELECT
         toStartOfYear(TimeStamp) year,
         toStartOfMonth(TimeStamp) month,
         toStartOfDay(TimeStamp) day,
         countState() block_number,
         sumState(Fees) AS sum_transaction_fees,
         sumState(TxCount) AS sum_transaction_count,
         sumState(TotalOutput) AS sum_total_output,
         sumState(BodySize) AS sum_body_size,
		 sumState(MintCount) AS sum_mint_count,
         sumState(MetaCount) AS sum_meta_count,
         sumState(Cip25AssetCount) AS sum_cip25_assets,
         sumState(Cip20Count) AS sum_cip20
       FROM cardano.chblock
       GROUP BY toStartOfYear(TimeStamp), toStartOfMonth(TimeStamp),toStartOfDay(TimeStamp);`

const CHBlockDayHour = `CREATE MATERIALIZED VIEW cardano.chblock_day_hour
       ENGINE = AggregatingMergeTree()
       ORDER BY (year,month,day,hour)
       AS SELECT
         toStartOfYear(TimeStamp) year,
         toStartOfMonth(TimeStamp) month,
         toStartOfDay(TimeStamp) day,
          toStartOfDay(TimeStamp) hour,
         countState() block_number,
         sumState(Fees) AS sum_transaction_fees,
         sumState(TxCount) AS sum_transaction_count,
         sumState(TotalOutput) AS sum_total_output,
         sumState(BodySize) AS sum_body_size,
		 sumState(MintCount) AS sum_mint_count,
         sumState(MetaCount) AS sum_meta_count,
         sumState(Cip25AssetCount) AS sum_cip25_assets,
         sumState(Cip20Count) AS sum_cip20
       FROM cardano.chblock
       GROUP BY toStartOfYear(TimeStamp), toStartOfMonth(TimeStamp),toStartOfDay(TimeStamp),toStartOfHour(TimeStamp);`

type Block struct {
	BlockNumber              uint64    `json:"block_number"`
	Slot                     uint64    `json:"slot"`
	BodySize                 uint32    `json:"body_size"`
	Epoch                    uint32    `json:"epoch"`
	EpochSlot                uint32    `json:"epoch_slot"`
	Era                      string    `json:"era"`
	Hash                     string    `json:"hash"`
	IssuerVkey               string    `json:"issuer_vkey"`
	SlotLeader               string    `json:"slot_leader"`
	TxCount                  uint32    `json:"tx_count"`
	Fees                     uint64    `json:"fees"`
	TotalOutput              uint64    `json:"total_output"`
	InputCount               uint32    `json:"input_count"`
	OutputCount              uint32    `json:"output_count"`
	MintCount                uint64    `json:"mint_count"`
	MetaCount                uint32    `json:"metadata_count"`
	NativeWitnessesCount     uint32    `json:"native_witnesses_count"`
	PlutusDatumCount         uint32    `json:"plutus_datum_count"`
	PlutusRdmrCount          uint32    `json:"plutus_redeemer_count"`
	PlutusWitnessesCount     uint32    `json:"plutus_witnesses_count"`
	Cip25AssetCount          uint32    `json:"cip25_asset_count"`
	Cip20Count               uint32    `json:"cip20_count"`
	PoolRegistrationCount    uint32    `json:"pool_registration_count"`
	PoolRetirementCount      uint32    `json:"pool_retirement_count"`
	StakeDelegationCount     uint32    `json:"stake_delegation_count"`
	StakeRegistrationCount   uint32    `json:"stake_registration_count"`
	StakeDeregistrationCount uint32    `json:"stake_deregistration_count"`
	Confirmations            uint32    `json:"confirmations"`
	TimeStamp                time.Time `json:"datetime,omitempty"`
}

func (cb *ChBlock) Drop(conn clickhouse.Conn) (err error) {

	if cb.Table == "" {
		err = errors.New("table is missing")
		return
	}
	var (
		query       = "DROP TABLE IF EXISTS " + cb.Table
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()
	err = conn.Exec(ctx, query)
	if err == nil {
		query = "DROP TABLE IF EXISTS " + cb.Table + "_epoch"
		err = conn.Exec(context.Background(), query)
		query = "DROP TABLE IF EXISTS " + cb.Table + "_hour"
		err = conn.Exec(context.Background(), query)
	}
	return
}

func (cb *ChBlock) CreateTable(conn clickhouse.Conn) (err error) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	err = conn.Exec(ctx, CHBlock)
	if err == nil {
		err = conn.Exec(context.Background(), CHBlockEpoch)
		err = conn.Exec(context.Background(), CHBlockDay)
		err = conn.Exec(context.Background(), CHBlockDayHour)
	}
	return
}

func (cb *ChBlock) PrepareBatch(conn clickhouse.Conn) (err error) {

	var ctx = context.Background()
	/*
		var ctx, _ = context.WithCancel(context.Background())
		defer cancel()
	*/
	var query = "INSERT INTO " + cb.Table
	cb.Batch, err = conn.PrepareBatch(ctx, query)
	return
}

func (cb *ChBlock) BulkInsert() error {
	batch := cb.Batch
	for _, v := range cb.Blocks {
		batch.Append(
			v.BlockNumber,
			v.Slot,
			v.BodySize,
			v.Epoch,
			v.EpochSlot,
			v.Era,
			v.Hash,
			v.IssuerVkey,
			v.SlotLeader,
			v.TxCount,
			v.Fees,
			v.TotalOutput,
			v.InputCount,
			v.OutputCount,
			v.MintCount,
			v.MetaCount,
			v.NativeWitnessesCount,
			v.PlutusDatumCount,
			v.PlutusRdmrCount,
			v.PlutusWitnessesCount,
			v.Cip25AssetCount,
			v.Cip20Count,
			v.PoolRegistrationCount,
			v.PoolRetirementCount,
			v.StakeDelegationCount,
			v.StakeRegistrationCount,
			v.StakeDeregistrationCount,
			v.Confirmations,
			v.TimeStamp,
		)
	}
	return batch.Send()
}

func (cb *ChBlock) BulkInsertStruct() error {
	batch := cb.Batch
	for _, v := range cb.Blocks {
		err := batch.AppendStruct(
			&v)
		if err != nil {
			log.Error().Msgf("Append block_number %s - slot_number %s  failed", v.BlockNumber, v.Slot)
		}
	}
	return batch.Send()
}

type BlockEpoch struct {
	Epoch                      uint32
	BlockNumber                uint64
	SumTransactionFees         uint32
	SumTransactionCount        uint32
	SumTotalOutput             uint32
	SumBodySize                uint32
	AverageTransactionPerBlock int
	AverageTransactionFees     int
	AverageTransactionCount    int
}

const QueryBlockEpoch = `
		SELECT  Epoch, 
		countMerge(block_number) AS count_block_number, 
		sumMerge(sum_transaction_fees) AS sum_transaction_fees,
		sumMerge(sum_transaction_count) AS sum_transaction_count, 
		sumMerge(sum_total_output) AS sum_total_output, 
		sumMerge(sum_body_size) AS sum_body_size,
		intDiv(	sum_transaction_count,count_block_number) AS avg_transaction_block,
		intDiv(	sum_body_size,count_block_number) AS avg_body_size,
		intDiv(sum_transaction_fees,sum_transaction_count) AS avg_transaction_fees,
		intDiv(sum_total_output,sum_transaction_count) AS avg_transaction_output
		FROM $1
 		GROUP BY $2 ORDER BY $3;
	`

const QueryBlockDay = `
      SELECT  day, 
		countMerge(block_number) AS count_block_number, 
		sumMerge(sum_transaction_fees) AS sum_transaction_fees,
		sumMerge(sum_transaction_count) AS sum_transaction_count, 
		sumMerge(sum_total_output) AS sum_total_output, 
		sumMerge(sum_body_size) AS sum_body_size,
		intDiv(	sum_transaction_count,count_block_number) AS avg_transaction_block,
		intDiv(	sum_body_size,count_block_number) AS avg_body_size,
		intDiv(sum_transaction_fees,sum_transaction_count) AS avg_transaction_fees,
		intDiv(sum_total_output,sum_transaction_count) AS avg_transaction_output
		FROM $1
 		GROUP BY $2 ORDER BY $3;
`

func (cb *BlockEpoch) Select(conn clickhouse.Conn, table string) (result []BlockEpoch, err error) {
	var (
		query       = QueryBlockEpoch
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()
	if err = conn.Select(ctx, &result, query, table, "Epoch", "Epoch"); err != nil {
		log.Error().Msgf("%v", err)
	}
	return
}
