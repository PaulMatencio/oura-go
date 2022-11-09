package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/rs/zerolog/log"
)

type BlockAggregating struct {
	Epoch                      uint32
	BlockNumber                uint64
	SumTransactionFees         uint32
	SumTransactionCount        uint32
	SumTotalOutput             uint32
	AverageTransactionPerBlock int
	AverageTransactionFees     int
	AverageTransactionCount    int
}

const QueryBlockAggregating = `
		SELECT  Epoch, 
		countMerge(block_number) AS block_number, 
		sumMerge(sum_transaction_fees) AS sum_transaction_fees,
		sumMerge(sum_transaction_count) AS sum_transaction_count, 
		sumMerge(sum_total_output) AS sum_total_output, 
		intDiv(	sum_transaction_count,block_number) AS avg_transaction_block,
		intDiv(sum_transaction_fees,sum_transaction_count) AS avg_transaction_fees,
		intDiv(sum_total_output,sum_transaction_count) AS avg_transaction_output
		FROM $1
 		GROUP BY $2 ORDER BY $3;
	`

func (cb *BlockAggregating) Select(conn clickhouse.Conn, table string) (result []BlockAggregating, err error) {
	query := QueryBlockAggregating
	if err = conn.Select(context.Background(), &result, query, table, "Epoch", "Epoch"); err != nil {
		log.Error().Msgf("%v", err)
	}
	return
}
