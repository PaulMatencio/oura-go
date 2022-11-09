package types

import (
	"context"
	"github.com/blockfrost/blockfrost-go"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

type BlockF struct {
	Block     interface{}
	Blocks    []interface{}
	LastBlock int
}

type BlockfB struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Time          int                `json:"time" bson:"time"`
	Height        int                `json:"height" bson:"height"`
	Hash          string             `json:"hash" bson:"hash"`
	Slot          int                `json:"slot" bson:"slot" `
	Epoch         int                `json:"epoch" bson:"epoch"`
	EpochSlot     int                `json:"epoch_slot" json:"epoch_slot"`
	SlotLeader    string             `json:"slot_leader" bson:"slot_leader"`
	Size          int                `json:"size" bson:"size"`
	TxCount       int                `json:"tx_count" bson:"tx_count"`
	Output        string             `json:"output" bson:"output"`
	Fees          string             `json:"fees" bson:"fees"`
	BlockVrf      string             `json:"block_vrf" bson:"block_vrf"`
	PreviousBlock string             `json:"previous_block"bson:"previous_block" `
	NextBlock     string             `json:"next_block" bson:"next_block"`
	Confirmations int                `json:"confirmations" bson:"confirmations"`
}

type GetBlockOptions struct {
	ApiClient      blockfrost.APIClient
	ApiQueryParams blockfrost.APIQueryParams
}

func (block *BlockF) GetBlock(projectId string, blockNumber int64) error {
	var (
		options = blockfrost.APIClientOptions{
			ProjectID: projectId,
		}
		blockf blockfrost.Block
		err    error
	)

	blockOptions := GetBlockOptions{
		ApiClient: blockfrost.NewAPIClient(options),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	blockf, err = blockOptions.ApiClient.Block(ctx, strconv.FormatInt(blockNumber, 10))
	if err == nil {
		/*
			var blockFB BlockfB
			copy(&blockf, &blockFB)
			block.Block = blockFB
		*/
		block.Block = blockf
		block.LastBlock = blockf.Height
	} else {
		log.Error().Msgf("Get block error: %v\n", err)
	}

	return err

}

func (block *BlockF) GetBlocks(getOptions GetBlockOptions, blockNumber int64) error {

	var (
		blockf []blockfrost.Block
		err    error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	block.Blocks = nil
	blockf, err = getOptions.ApiClient.BlocksNext(ctx, strconv.FormatInt(blockNumber, 10))
	if err == nil {
		for _, v := range blockf {
			var blocfB BlockfB
			copy(&v, &blocfB)
			block.Blocks = append(block.Blocks, blocfB)
		}
		if len(blockf) > 0 {
			lastBlock := blockf[len(blockf)-1 : len(blockf)][0]
			block.LastBlock = lastBlock.Height
		}

	} else {
		log.Error().Msgf("Get blocks error: %v\n", err)
	}

	return err

}

func copy(from *blockfrost.Block, to *BlockfB) {
	copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}
