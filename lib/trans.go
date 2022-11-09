package lib

import (
	"errors"
	"fmt"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func BuildTran(flags types.Options, req *mongodb.MongoDB, Filter interface{}) (Trans types.Trans, err error) {

	var (
		Tx   types.Tx
		req1 = mongodb.MongoDB{
			Option:     flags.MongoOptions,
			Database:   flags.MongoDatabase,
			Collection: flags.MongoCollection,
		}
		opt1 options.FindOneOptions
	)

	req1.Uri = req.Uri
	req1.Database = req.Database
	if req1.Client, err = req1.Connect(); err != nil {
		LogError(err, "connect req1")
		return
	}
	defer req1.DisConnect()
	req.Collection = "tx"
	if r, err1 := FindOne(&opt1, Tx, Filter, req); err1 == nil {
		filter1 := r.Context.TxHash
		filter2 := r.Context.BlockHash
		Filter2 := bson.M{"context.tx_hash": filter1, "context.block_hash": filter2}
		return BuildTran1(flags, r, req, &req1, Filter2)
	} else {
		err = err1
		log.Error().Err(err).Msgf("find one filer:%v", Filter)
	}

	return
}

func BuildTransSeq(flags types.Options, result []types.Tx, req *mongodb.MongoDB, req1 *mongodb.MongoDB) (total int, terror int) {
	/*
		flags is used by BuildTran1
	*/
	var (
		Trans []interface{}
		req2  = mongodb.MongoDB{
			Option:   req.Option,
			Uri:      req.Uri,
			Database: req.Database,
		}
		err error
	)

	if req2.Client, err = req2.Connect(); err != nil {
		LogError(err, "connect req2")
		return
	}
	defer req2.DisConnect()

	for _, r := range result {
		filter1 := r.Context.TxHash
		filter2 := r.Context.BlockHash
		log.Trace().Msgf("Sequential building transaction document %s %s\n", filter1, filter2)
		Filter := bson.M{"context.tx_hash": filter1, "context.block_hash": filter2}
		if trans, err := BuildTran1(flags, r, req, req1, Filter); err == nil {
			log.Trace().Msgf("Finger_print %s - Tx_hash %s\n", trans.Fingerprint, trans.Context.TxHash)
			Trans = append(Trans, trans)
			total += 1
		} else {
			LogError(err, "building transaction")
			terror += 1
		}
	}

	log.Info().Msgf("Sequential Bulk upload %d transaction documents", total)
	req2.Collection = "trans"
	if len(Trans) > 0 && flags.Upload {
		if _, err := req2.InsertMany(Trans); err != nil {
			LogError(err, "insert many")
		}
	}

	return
}

/*
	TODO -> Bug with append when appended elem is an address ?? change result []types.TxC  to result *[]types.TxB
    concurrent build upload
    for every tx in  result {
         trans = buildTrans (aggregate tx, utxo, stxi, mint, meta, etc ..)
         append trans to Trans
    }
    // bulk upload
    insertMany( trans)
*/

func BuildTransCon(flags types.Options, result []types.Tx, req *mongodb.MongoDB, req1 *mongodb.MongoDB) (total int, terror int) {
	/*
		flags is used by BuildTran1
	*/
	type Response struct {
		Trans types.Trans
		Err   error
	}
	var (
		ch    = make(chan *Response)
		Trans []interface{}

		// req2 used for write
		req2 = mongodb.MongoDB{
			Option:   req.Option,
			Uri:      req.Uri,
			Database: req.Database,
		}

		request, receive = 0, 0
		err              error
	)
	request = len(result)
	if req2.Client, err = req2.Connect(); err != nil {
		LogError(err, "connect req2")
		return
	}
	defer req2.DisConnect()
	for _, r := range result {
		go func(r types.Tx, req *mongodb.MongoDB, req1 *mongodb.MongoDB) {
			var resp Response
			filter1 := r.Context.TxHash
			filter2 := r.Context.BlockHash
			Filter := bson.M{"context.tx_hash": filter1, "context.block_hash": filter2}
			resp.Trans, resp.Err = BuildTran1(flags, r, req, req1, Filter)
			ch <- &resp
		}(r, req, req1)
	}

	receive = 0

	for {
		if len(result) == 0 {
			return total, terror
		}
		select {
		case rec := <-ch:
			receive++
			if rec.Err == nil {
				Trans = append(Trans, rec.Trans)
			} else {
				terror += 1
				log.Trace().Stack().Msgf("Error building transaction  %v %s", rec.Err, rec.Trans.Transaction.Hash)
			}
			if receive == request {
				total += len(Trans)
				log.Info().Msgf("Concurrent bulk uploading %d transaction documents", total)
				if flags.Upload {
					req2.Collection = "trans"
					if _, err := req2.InsertMany(Trans); err != nil {
						LogError(err, "insert many req2")
						terror += 1
					}
				}

				return total, terror
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Printf(".")
		}
	}
}

/*
called by BuildTran
*/

func BuildTran1(flags types.Options, tx types.Tx, req *mongodb.MongoDB, req1 *mongodb.MongoDB, Filter interface{}) (Trans types.Trans, err error) {

	var (
		assetCount  int64
		inputCount  = tx.Transaction.InputCount
		outputCount = tx.Transaction.OutputCount
		mintCount   = tx.Transaction.MintCount
		cip25       = false
	)
	Trans.Context = tx.Context
	fp := strings.Split(tx.Fingerprint, ".")
	Trans.Fingerprint = fp[0] + ".trans." + fp[2]
	Trans.Transaction = tx.Transaction

	/*
	   get Utxo outputs
	*/
	if outputCount > 0 {
		txOutputs, err := GetUtxoOutput(flags, outputCount, req, Filter)
		if err == nil {
			Trans.TxOutput = txOutputs
			Trans.TxMeta.OutputCount = len(txOutputs)
			assetCount += int64(len(txOutputs))
		} else {
			return Trans, err
		}
	} else {
		err := errors.New(fmt.Sprintf("missing utxo - Filter %v", Filter))
		return Trans, err
	}

	/*
		get stxi  inputs
	*/
	if inputCount > 0 {
		txInputs, err := GetStxiInput(flags, inputCount, req, req1, Filter)
		if err == nil {
			Trans.TxInput = txInputs
			Trans.TxMeta.InputCount = len(txInputs)
		} else {
			return Trans, err
		}
	} else {
		err := errors.New(fmt.Sprintf("missing stxi - Filter %v", Filter))
		return Trans, err
	}

	/*
		get Mint Asset
	*/

	if mintCount > 0 {
		minAssets, err := GetMint(req, mintCount, Filter)
		if err == nil {
			Trans.MintAsset = minAssets
			Trans.TxMeta.MintCount = int64(len(minAssets))
		}
	}

	/*
		get Meta for regular transaction
	*/

	req.Collection = "meta"
	var (
		m           types.Meta
		mm          []types.Meta
		findOptions options.FindOptions
	)

	mm, _ = FindAll(&findOptions, m, Filter, req)
	for _, v := range mm {
		meta := v.Metadata
		Trans.TxMeta.MetaCount = len(mm)
		if meta != nil {
			Trans.Metadata = append(Trans.Metadata, meta)
			metadata, err := PrimitiveDtoMap(v.Metadata)

			if err == nil {
				if _, ok := metadata["text_scalar"]; ok {
					var metaScalar types.MetaScalar
					data, err := bson.Marshal(v.Metadata)
					if err == nil {
						err = bson.Unmarshal(data, &metaScalar)
						if err == nil {
							Trans.MetaScalar = append(Trans.MetaScalar, metaScalar)
						} else {
							log.Error().Err(err)
							LogError(err, "bson unmarshal")
						}
					} else {
						LogError(err, "bson marshal")
					}
					Trans.TxMeta.MetaLabel = "text_scalar"
					continue
				}

				Trans.TxMeta.MetaLabel = fmt.Sprintf("%s", metadata["label"])
				if metadata["label"] == "721" {
					// Trans.TxMeta.MetaLabel = "721"
					cip25 = true
					continue
				}

				if metadata["label"] == "674" {
					var meta674 types.Meta674
					data, err := bson.Marshal(meta)
					if err == nil {
						err = bson.Unmarshal(data, &meta674)
						if err == nil {
							Trans.Meta674 = append(Trans.Meta674, meta674)
						} else {
							LogError(err, "bson unmarshal")
							var meta674a types.Meta674a
							err = bson.Unmarshal(data, &meta674a)
							if err == nil {
								meta674.Label = meta674a.Label
								meta674.MapJson.Msg = append(meta674.MapJson.Msg, meta674a.MapJson.Msg)
								Trans.Meta674 = append(Trans.Meta674, meta674)
							}
						}
					} else {
						LogError(err, "bson marshal")
					}
					continue
				}

				if metadata["label"] == "3322" {
					var meta3322 types.Meta3322
					data, err := bson.Marshal(meta)
					if err == nil {
						err = bson.Unmarshal(data, &meta3322)
						if err == nil {
							Trans.Meta3322 = append(Trans.Meta3322, meta3322)
						} else {
							LogError(err, "bson unmarshal")
						}
					} else {
						LogError(err, "bson marshal")
					}
					continue
				}
			} else {
				LogError(err, "Mongo primitive to map[string]string")
			}
		}
	}

	//  Get Collateral
	collateral, err := GetCollateral(req, req1, Filter)
	if err == nil {
		Trans.Collateral = collateral
		Trans.TxMeta.CollCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get collateral -  Filter %v", Filter))
			return Trans, err
		}
	}

	// datum

	datums, err := GetPlutusDatums(req, Filter)
	if err == nil {
		Trans.PlutusDatum = datums
		Trans.TxMeta.PlutusDatumCount = len(datums)
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get datums  - Filter %v", Filter))
			return Trans, err
		}
	}

	// Cip25 asset
	if cip25 {

		cip25pAssets, err := GetCip25p(req, mintCount, Filter)
		if err == nil {
			Trans.Cip25pAsset = cip25pAssets
			Trans.TxMeta.Cip25AssetCount = len(cip25pAssets)
		} else {
			if err != mongo.ErrNoDocuments {
				err = errors.New(fmt.Sprintf("get cip25p  - Filter %v", Filter))
				return Trans, err
			}
		}
	}

	/*
		plutus redeemer
	*/
	reDeemers, err := GetRedeemer(req, Filter)
	if err == nil {
		Trans.PlutusRedeemer = reDeemers
		Trans.TxMeta.PlutusRdmrCount = len(reDeemers)
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get plutus redeemers - Filter %v", Filter))
			return Trans, err
		}
	}

	//  Plutus native Witness
	//  witn

	nativeWitness, err := GetNativeWitness(req, Filter)
	if err == nil {
		Trans.NativeWitness = nativeWitness
		Trans.TxMeta.NativeWitnessesCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get native witness - Filter %v", Filter))
			return Trans, err
		}
	}

	//  Plutus Witness
	//  witp
	plutusWitness, err := GetPlutusWitness(req, Filter)
	if err == nil {
		Trans.PlutusWitness = plutusWitness
		Trans.TxMeta.PlutusWitnessesCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get plutus witness - Filter %v", Filter))
			return Trans, err
		}
	}

	// Pool registration

	poolRegistration, err := GetPoolRegistration(req, Filter)
	if err == nil {
		Trans.PoolRegistration = poolRegistration
		Trans.TxMeta.PoolRegistrationCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get pool registration - Filter %v", Filter))
			return Trans, err
		}
	}

	// Pool retirement
	poolRetirement, err := GetPoolRetirement(req, Filter)
	if err == nil {
		Trans.PoolRetirement = poolRetirement
		Trans.TxMeta.PoolRetirementCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get pool retirement - Filter %v", Filter))
			return Trans, err
		}
	}

	// Stake delegation
	stakeDelegation, err := GetStakeDelegation(req, Filter)
	if err == nil {
		Trans.StakeDelegation = stakeDelegation
		Trans.TxMeta.StakeDelegationCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get Stake delegation - Filter %v", Filter))
			return Trans, err
		}
	}

	// Stake registration
	stakeRegistration, err := GetStakeRegistration(req, Filter)
	if err == nil {
		Trans.StakeRegistration = stakeRegistration
		Trans.TxMeta.StakeDelegationCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get Stake registration - Filter %v", Filter))
			return Trans, err
		}
	}

	// Stake de-registration
	stakeDeregistration, err := GetStakeDeregistration(req, Filter)
	if err == nil {
		Trans.StakeDeregistration = stakeDeregistration
		Trans.TxMeta.StakeDeregistrationCount = 1
	} else {
		if err != mongo.ErrNoDocuments {
			err = errors.New(fmt.Sprintf("get Stake deregistration - Filter %v", Filter))
			return Trans, err
		}
	}
	return
}
