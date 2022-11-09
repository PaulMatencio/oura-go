package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

func GetMarketPlaces(mktPlace string, mPlaces []types.MktPlace) (mktPlaces []types.MktPlace) {
	/*
		accept generic marketplace
	*/
	for _, v := range mPlaces {
		if strings.Contains(v.Name, mktPlace) {
			mktPlaces = append(mktPlaces, v)
		}
	}
	return
}

/*
	Build a NFT transaction for a given filter ( tx_hash for instance)
    call buildNftTran1
    ( it is for testing purpose)
*/

func BuildNftTran(flags types.Options, req *mongodb.MongoDB,
	req1 *mongodb.MongoDB, Filter interface{}) {
	var (
		err  error
		req3 = mongodb.MongoDB{
			Uri:      req.Uri,
			Database: req.Database,
			Option:   req.Option,
		}
	)
	if req3.Client, err = req3.Connect(); err != nil {
		log.Error().Err(err).Msgf("mongodb connection")
		return
	}
	defer req3.DisConnect()
	if trans, err := BuildNftTrans1(flags, req, req1, &req3, Filter); err == nil {
		if b, err := bson.MarshalExtJSON(&trans, false, true); err == nil {
			if flags.PrintIt {
				utils.PrintJson(string(b))
			}
		} else {
			log.Error().Stack().Err(err).Msg("Parse transaction failed")
		}
	} else {
		if !strings.Contains(err.Error(), "skip") {
			log.Error().Err(err).Msgf("BuildNftTrans1")
		} else {
			log.Info().Err(err).Msgf("skip transaction")
		}
	}
	return
}

/*
			Build NFT transactions

	          find utxo events ( --limit)  where the utxo output.address == market place address
	          for every returned utxo event
	            call sequentially BuildNftTran1  to build a transaction
	          end for loop
*/
func BuildNftTransSeq(flags types.Options, mktpl types.MktPlace, objectId primitive.ObjectID,
	req *mongodb.MongoDB, req1 *mongodb.MongoDB) (total int, terror int,
	lastId primitive.ObjectID) {
	var (
		err         error
		r           types.Utxo
		rr          []types.Utxo
		findOptions options.FindOptions
		req3        = mongodb.MongoDB{
			Uri:      req.Uri,
			Database: req.Database,
			Option:   req.Option,
		}
	)
	if req3.Client, err = req3.Connect(); err != nil {
		log.Error().Err(err).Msgf("mongodb connection")
		return
	}
	defer req3.DisConnect()
	req.Collection = "utxo"
	fmt.Println("start objectId", objectId)
	Filter := bson.M{"tx_output.address": mktpl.Address, "_id": bson.M{"$gt": objectId}}
	findOptions.SetLimit(flags.Limit)
	if _, rr, err = Find(&findOptions, r, Filter, req); err == nil {
		for _, r := range rr {
			Filter1 := bson.M{"context.tx_hash": r.Context.TxHash, "context.block_hash": r.Context.BlockHash}
			if trans, err := BuildNftTrans1(flags, req, req1, &req3, Filter1); err == nil {
				if b, err := bson.MarshalExtJSON(&trans, false, true); err == nil {
					//	if b, err := json.Marshal(&trans); err == nil {
					if flags.PrintIt {
						utils.PrintJson(string(b))
					} else {
						// fmt.Println(trans.Fingerprint)
					}
				} else {
					terror++
					log.Error().Stack().Err(err).Msg("Parse transaction failed")
				}
			} else {
				// if err.Error() != "skip" {
				if !strings.Contains(err.Error(), "skip") {
					log.Error().Err(err).Msgf("building transaction %s", trans.Transaction.Hash)
				} else {
					log.Info().Err(err).Msgf("skip transaction %s", trans.Transaction.Hash)
				}
			}
		}
		total = len(rr)
		if len(rr) > 0 {
			lastId = rr[len(rr)-1].ID
			fmt.Println("lastId", lastId)
		}

	} else {
		if err != mongo.ErrNoDocuments {
			log.Error().Stack().Msgf("Filter %v - %s ", "Utxo", Filter, err)
		}
	}
	return
}

/*

  find utxo events ( --limit)  where the utxo output.address == market place address
  for every returned utxo event
     call in parallel  BuildNftTran1  to build a transaction
  end for loop

  wait until all transactions have been completed

*/

func BuildNftTransCon(flags types.Options, mktpl types.MktPlace, objectId primitive.ObjectID,
	req *mongodb.MongoDB, req1 *mongodb.MongoDB,
	req2 *mongodb.MongoDB) (total int, terror int, lastId primitive.ObjectID) {

	/*
				BuildNftTrans:
				Find  in collection UTXO  N  (  --limit) transaction where tx_output.address = market_place address
			    For each returned transaction (tx_hash) {
					 BuildNftTrans1
		        }
	*/
	type Response struct {
		Trans types.TransNft
		Err   error
	}
	var (
		err              error
		r                types.Utxo
		rr               []types.Utxo
		ch               = make(chan *Response)
		Trans            []interface{}
		request, receive = 0, 0
		findOptions      options.FindOptions
		req3             = mongodb.MongoDB{
			Uri:      req.Uri,
			Database: req.Database,
			Option:   req.Option,
		}
	)
	if req3.Client, err = req3.Connect(); err != nil {
		log.Error().Err(err).Msgf("mongodb connection")
		return
	}
	defer req3.DisConnect()
	req.Collection = "utxo"
	findOptions.SetLimit(flags.Limit)
	log.Info().Msgf("start from objectId %v", objectId)
	Filter := bson.M{"tx_output.address": mktpl.Address, "_id": bson.M{"$gt": objectId}}
	if _, rr, err = Find(&findOptions, r, Filter, req); err == nil {
		request = len(rr)
		for k, r := range rr {
			Filter1 := bson.M{"context.tx_hash": r.Context.TxHash, "context.block_hash": r.Context.BlockHash}
			go func(r types.Utxo, req *mongodb.MongoDB, req1 *mongodb.MongoDB, Filter1 bson.M, k int) {
				var resp Response
				resp.Trans, resp.Err = BuildNftTrans1(flags, req, req1, &req3, Filter1)
				ch <- &resp
			}(r, req, req1, Filter1, k)
		}
		receive = 0
		if len(rr) > 0 {
			r1 := rr[len(rr)-1]
			lastId = r1.ID
			log.Info().Msgf("next id %v\n", lastId)
		}
		for {
			if len(rr) == 0 {
				return 0, 0, lastId
			}

			select {
			case rec := <-ch:
				receive++
				if rec.Err == nil || rec.Err == mongo.ErrNoDocuments {
					Trans = append(Trans, rec.Trans)
				} else {
					if !strings.Contains(rec.Err.Error(), "skip") {
						terror += 1
						log.Error().Err(rec.Err).Msgf("building transaction %s", rec.Trans.Transaction.Hash)
					} else {
						log.Warn().Err(rec.Err).Msgf("skip trans %s", rec.Trans.Transaction.Hash)
					}
				}
				if receive == request {
					total += len(Trans)
					log.Info().Msgf("Concurrent bulk uploading %d transaction documents", total)
					/* upload  */
					req2.Collection = "trans_nft"
					if flags.Upload {
						if _, err := req2.InsertMany(Trans); err != nil {
							log.Error().Err(err).Msgf("many insert to %s", req2.Collection)
							terror += 1
						}
					}

					if flags.PrintIt {
						for _, v := range Trans {
							v1 := v.(types.TransNft)
							if b, err := json.Marshal(v1); err == nil {
								utils.PrintJson(string(b))
							} else {
								log.Error().Err(err).Msg("marshal trans")
							}
						}
					}
					//  save block number before returning
					return total, terror, lastId
				}
			case <-time.After(100 * time.Millisecond):
				fmt.Printf(".")
			}
		}
	} else {
		if err != mongo.ErrNoDocuments {
			log.Error().Err(err).Msgf("Find Filter %s", Filter)
		}
	}
	return total, terror, lastId
}

/*
	BuildNftTrans1:
				Build NFT transaction ( tx_hash )  -> nft_trans {
						 add  tx event
                         add  tx_outputs event
                         add  tx_inputs event
                         add  tx_collateral event
                         add  tx_plutusdata event
                         add  tx_redeemer event
                         build and add  the summary
                         add  cip25p ( NFT metadata)
			    }

*/

func BuildNftTrans1(flags types.Options, req *mongodb.MongoDB, req1 *mongodb.MongoDB, req2 *mongodb.MongoDB, Filter interface{}) (Trans types.TransNft, err error) {

	var (
		inputCount, outputCount int
		mintCount               int64
		tx                      types.Tx
		opt1                    options.FindOneOptions
		summary                 types.TransNftSummary
	)

	/*
		add  the corresponding transaction market_place
	*/
	Trans.SetMarketPlace(flags.MarketPlace)
	for retry := 1; retry <= flags.MaxRetry; retry++ {
		req.Collection = "tx"
		if tx, err = FindOne(&opt1, tx, Filter, req); err == nil {
			inputCount = tx.Transaction.InputCount
			outputCount = tx.Transaction.OutputCount
			mintCount = tx.Transaction.MintCount
		} else {
			log.Error().Err(err).Msgf("collection %s - Filter %v", req.Collection, Filter)
			return Trans, err
		}
		if inputCount > 0 {
			break
		} else {
			err = errors.New("transaction input count is 0")
			if retry < flags.MaxRetry {
				log.Warn().Msgf("collection %s - Filter %v - retries %d", req.Collection, Filter, retry)
			} else {
				log.Error().Err(err).Msgf("collection %s - Filter %v - number of retries are exceeded %d", req.Collection, Filter, retry)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	/*
		add the corresponding transaction event
	*/
	Trans.SetTransaction(tx.Transaction)
	if mintCount > 0 {
		err := errors.New("skip mint transaction")
		return Trans, err
	}
	/*
		add  a new  transaction fingerprint ( unique key of the collection)
	*/
	fp := strings.Split(tx.Fingerprint, ".")
	if len(fp) != 3 {
		log.Error().Msgf("tx_hash %s - Wrong transaction finger print %s", tx.Context.TxHash, tx.Fingerprint)
		return Trans, err
	}
	Trans.SetFingerPrint(fp[0] + ".trans." + fp[2])
	/*
	   add transaction transaction
	*/
	Trans.SetContext(tx.Context)
	/*
	   add the corresponding utxo output event
	*/
	if outputCount > 0 {
		txOutputs, err := GetUtxoOutput(flags, outputCount, req1, Filter)
		if err == nil {
			Trans.SetTxOutput(txOutputs)
			Trans.SetOutputCount(len(txOutputs))

		} else {
			return Trans, err
		}
	}

	/*
	   add  the corresponding utxo input event
	*/
	if inputCount > 0 {
		txInputs, err := GetStxiInput(flags, inputCount, req, req1, Filter)
		if err == nil {
			Trans.SetTxInput(txInputs)
			Trans.SetInputCount(len(txInputs))
		} else {
			return Trans, err
		}
	}
	/*
	 add  the corresponding Collateral event
	*/
	req.Collection = "coll"
	collateral, err := GetCollateral(req, req1, Filter)
	if err == nil {
		Trans.SetCollateral(collateral)
		Trans.SetCollateralCount(1)
	} else {
		if err != mongo.ErrNoDocuments {
			log.Error().Err(err).Msgf("collection %s - Filter %v", req.Collection, Filter)
			return Trans, err
		}
	}
	/*
			add corresponding plutus data

		req.Collection = "dtum"
		dt, err := GetPlutusDatas(req, Filter)
		if err == nil {
			Trans.SetPlutusData(dt)
			Trans.SetDataCount(len(dt))
		} else {
			if err != mongo.ErrNoDocuments {
				log.Error().Err(err).Msgf("collection %s - Filter %v", req.Collection, Filter)
				return Trans, err
			}
		}

	*/

	/*
		add the corresponding redeemer event
	*/
	req.Collection = "rdmr"
	rd, err := GetRedeemer(req, Filter)
	if err == nil {
		Trans.SetRedeemer(rd)
		Trans.SetRedeemerCount(len(rd))
	} else {
		if err != mongo.ErrNoDocuments {
			log.Error().Err(err).Msgf("collection %s - Filter %v", req.Collection, Filter)
			return Trans, err
		}
	}

	/*
		Build Summary
	*/
	switch flags.MarketPlace.Name {
	case "jpgstore",
		"adapix",
		"cnftio",
		"spacebudz":
		if summary, err = BuildNftSummary(Trans); err == nil {
			if cip25pAsset, assetMintedBy, err1 := getCip25p(req2, Filter, summary); err1 == nil {
				Trans.SetCip25Asset(cip25pAsset)
				summary.SetAssetMintedBy(assetMintedBy)
			} else {
				err = err1
			}
			Trans.SetNftSummary(summary)
		}
	case "cryptodino":
		if summary, cip25pAsset, err1 := cryptoDino(req2, Filter, Trans); err1 == nil {
			Trans.SetNftSummary(summary)
			Trans.SetCip25Asset(cip25pAsset)
		} else {
			err = err1
		}
	default:
	}

	return Trans, err
}

func getCip25p(req *mongodb.MongoDB, filter interface{}, summary types.TransNftSummary) (cip25Asset types.Cip25pAsset,
	mintedBy types.AssetMintedBy, err error) {
	var cip25 types.Cip25p
	if cip25, err = GetCip25a(req, summary.AssetPolicy, summary.AssetAscii); err == nil {
		cip25Asset = cip25.Cip25Asset
		mintedBy.TxHash = cip25.Context.TxHash
		mintedBy.Timestamp = cip25.Context.Timestamp
		return
	} else {
		switch err {
		case mongo.ErrNoDocuments:
			log.Warn().Err(err).Msgf("collection %s - Filter %v", req.Collection, filter)
			err = nil
		default:
			err = errors.New(fmt.Sprintf("collection %s - Filter %v", req.Collection, filter))
			return
		}
	}
	return
}

func cryptoDino(req *mongodb.MongoDB, filter interface{}, Trans types.TransNft) (summary types.TransNftSummary, cip25Asset types.Cip25pAsset, err error) {

	var (
		filter2   = bson.M{"context.tx_hash": Trans.Context.TxHash}
		metas     []types.Meta405
		operation string
	)
	for _, v := range GetMeta(req, filter2) {
		var meta types.Meta405
		meta, err = utils.PrimitiveDtoStruct(meta, v)
		operation = meta.MapJson.Op
		if operation == "buy" {
			metas = append(metas, meta)
		}
	}
	if len(metas) > 0 {
		if summary, err = BuildNftSummary1(Trans, metas); err == nil {

			if cip25Asset, assetMintedBy, err := getCip25p(req, filter, summary); err == nil {
				Trans.SetCip25Asset(cip25Asset)
				summary.SetAssetMintedBy(assetMintedBy)
			}
			Trans.SetNftSummary(summary)
		}
	} else {
		err = errors.New(fmt.Sprintf("skip %s transaction", operation))
		return
	}
	return
}

func BuildNftSummary(trans types.TransNft) (summary types.TransNftSummary, err error) {

	type Asst struct {
		address       string
		assetAscii    string
		AssetQuantity int
		AssetPolicy   string
		FingerPrint   string
		asset         types.Assets
	}

	var (
		inMap, outMap                                                 = make(map[string]Asst), make(map[string]Asst)
		tInput, tOutput, marketFee, otherOutput, otherSpent, returned int64
	)
	summary.SetMarketPlace(trans.GetMarketPlace())
	summary.SetTimestamp(trans.Context.Timestamp)
	summary.SetTxHash(trans.Context.TxHash)

	if len(trans.PlutusRedeemer) > 0 {
		summary.SetPurpose(trans.PlutusRedeemer[0].Purpose)
	}

	fee := trans.Transaction.Fee
	summary.SetFee(fee)
	/*
		txOutput.[]Assets struct -> map[string]Asst
	*/
	mapout := make(map[string]int64)
	for _, txOutput := range trans.TxOutput {
		key := txOutput.Address
		tOutput += txOutput.Amount
		if len(txOutput.Assets) > 0 {
			assetMap := AssetsToMap(txOutput.Assets)
			for k, v := range assetMap {
				var a Asst
				a.address = key
				a.AssetQuantity = v.Amount
				a.AssetPolicy = v.Policy
				a.FingerPrint = v.FingerPrint
				a.asset = v
				outMap[k] = a
			}
		}
		mapout[key] = txOutput.Amount
	}
	summary.SetTotalOutput(tOutput)

	/*
		txInput.[]Assets struct -> map[string]Asst
	*/
	for _, txInput := range trans.TxInput {
		key := txInput.Address
		tInput += txInput.Amount
		if len(txInput.Assets) > 0 {
			assetMap := AssetsToMap(txInput.Assets)
			for k, v := range assetMap {
				var a Asst
				a.address = key
				a.AssetQuantity = v.Amount
				a.AssetPolicy = v.Policy
				a.FingerPrint = v.FingerPrint
				a.asset = v
				inMap[k] = a
			}
			/*
					amount of every input address which is not present
				    in  output
			*/
			if _, ok := mapout[txInput.Address]; !ok {
				otherSpent += txInput.Amount
			}
		}
	}

	summary.SetTotalInput(tInput)
	summary.SetOtherSpent(otherSpent)
	/*
		Look for the selling asset and the buyer address
	*/
	for k, v := range inMap {
		if val, ok := outMap[k]; ok {
			if v.address != val.address {
				summary.SetAssetName(k)
				/* take the amount from the utxo output */
				summary.SetAssetQuantity(int64(val.AssetQuantity))
				summary.SetAssetPolicy(v.AssetPolicy)
				summary.SetAssetFingerPrint(v.FingerPrint)
				summary.SetFromAddress(v.address)
				summary.SetToAddress(val.address)
				break
			}
		}
	}

	// total amount that were not going to neither the buyer nor to the marketplace
	// this should be  the amount received buy the seller
	var (
		maxOtherOutput int64
		sellerAddress  string
	)
	for _, txOutput := range trans.TxOutput {
		if txOutput.Address == summary.ToAddress {
			returned += txOutput.Amount
		}
		if txOutput.Address == trans.GetMarketPlace().Address {
			marketFee += txOutput.Amount
		}
		if txOutput.Address != summary.ToAddress && txOutput.Address != summary.MarketPlace.Address {
			if txOutput.Amount >= maxOtherOutput {
				maxOtherOutput = txOutput.Amount
				sellerAddress = txOutput.Address
			}
			otherOutput += txOutput.Amount
		}
	}
	summary.SetAdaBack(returned)
	summary.SetAdaSpent(tInput - returned - otherSpent)
	summary.SetMarketReceived(marketFee)
	summary.SetSellerReceived(maxOtherOutput)
	summary.SetSellerAddress(sellerAddress)
	summary.SetOtherReceived(otherOutput - maxOtherOutput)
	switch summary.GetMarketPlaceName() {
	case "cnftio":
		summary.SetAdaSpent(tInput - returned)
		summary.SetAssetPrice(tInput - otherSpent)
	default:
		summary.SetAssetPrice(otherOutput + marketFee)
	}
	return
}

/*

 */

func BuildNftSummary1(trans types.TransNft, metas []types.Meta405) (summary types.TransNftSummary, err error) {

	type Asst struct {
		address     string
		assetAscii  string
		AssetPolicy string
		FingerPrint string
		asset       types.Assets
	}

	var (
		inMap, outMap                                                             = make(map[string]Asst), make(map[string]Asst)
		tInput, tOutput, marketFee, otherOutput, otherSpent, assetPrice, returned int64
	)
	summary.SetMarketPlace(trans.GetMarketPlace())
	summary.SetTimestamp(trans.Context.Timestamp)
	summary.SetTxHash(trans.Context.TxHash)
	if len(metas) > 0 {
		if assetPrice, err = strconv.ParseInt(metas[0].MapJson.Price, 10, 64); err == nil {
			summary.SetAssetPrice(assetPrice)
		}

	}

	if len(trans.PlutusRedeemer) > 0 {
		summary.SetPurpose(trans.PlutusRedeemer[0].Purpose)
	}

	fee := trans.Transaction.Fee
	summary.SetFee(fee)

	mapout := make(map[string]int64)
	for _, txOutput := range trans.TxOutput {
		key := txOutput.Address
		tOutput += txOutput.Amount
		if len(txOutput.Assets) > 0 {
			assetMap := AssetsToMap(txOutput.Assets)
			for k, v := range assetMap {
				var a Asst
				a.address = key
				a.AssetPolicy = v.Policy
				a.FingerPrint = v.FingerPrint
				a.asset = v
				outMap[k] = a
			}
		}
		mapout[key] = txOutput.Amount
	}

	summary.SetTotalOutput(tOutput)

	for _, txInput := range trans.TxInput {
		key := txInput.Address
		tInput += txInput.Amount
		if len(txInput.Assets) > 0 {
			assetMap := AssetsToMap(txInput.Assets)
			for k, v := range assetMap {
				var a Asst
				a.address = key
				a.AssetPolicy = v.Policy
				a.FingerPrint = v.FingerPrint
				a.asset = v
				inMap[k] = a
			}
			/*
					amount of every input address which is not present
				    in  output
			*/
			if _, ok := mapout[txInput.Address]; !ok {
				otherSpent += txInput.Amount
			}
		}
	}

	summary.SetTotalInput(tInput)
	summary.SetOtherSpent(otherSpent)
	/*
		Look for the selling asset and the buyer address
	*/
	for k, v := range inMap {
		if val, ok := outMap[k]; ok {
			if v.address != val.address {
				summary.SetAssetName(k)
				summary.SetAssetPolicy(v.AssetPolicy)
				summary.SetAssetFingerPrint(v.FingerPrint)
				summary.SetFromAddress(v.address)
				summary.SetToAddress(val.address)
				break
			}
		}
	}

	// total amount that were not going to neither the buyer nor to the marketplace
	// this should be  the amount received buy the seller
	var (
		maxOtherOutput int64
		// sellerAddress  string
	)
	for _, txOutput := range trans.TxOutput {
		if txOutput.Address == summary.ToAddress {
			returned += txOutput.Amount
		}
		if txOutput.Address == trans.GetMarketPlace().Address {
			marketFee += txOutput.Amount
		}
		if txOutput.Address != summary.ToAddress && txOutput.Address != summary.MarketPlace.Address {
			if txOutput.Amount >= maxOtherOutput {
				maxOtherOutput = txOutput.Amount
				// sellerAddress = txOutput.Address
			}
			otherOutput += txOutput.Amount
		}
	}
	summary.SetAdaBack(returned)
	// summary.SetAdaSpent(tInput - back - otherSpent)
	summary.SetMarketReceived(marketFee)
	summary.SetSellerReceived(maxOtherOutput)
	// summary.SetSellerAddress(sellerAddress)
	// summary.SetOtherReceived(otherOutput - maxOtherOutput)
	return

}
