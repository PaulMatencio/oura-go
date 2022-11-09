package lib

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

/*
	get stxi

*/

func GetStxiInput(flags types.Options, inputCount int, req *mongodb.MongoDB, req1 *mongodb.MongoDB, Filter interface{}) ([]types.TxInput, error) {

	var (
		err      error
		r        types.Stxi
		rr       []types.Stxi
		txInputs []types.TxInput
		errs     []error
		// cur         *mongo.Cursor
		findOptions options.FindOptions
		req2        = mongodb.MongoDB{
			Uri:      req.Uri,
			Database: req.Database,
			Option:   req.Option,
		}
	)
	//ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	req2.Client, _ = req2.Connect()
	defer req2.DisConnect()
	if inputCount > 0 {
		req1.Collection = "stxi"
		findOptions.SetLimit(int64(inputCount))
		for retry := 1; retry <= flags.MaxRetry; retry++ {
			if _, rr, err = Find(&findOptions, r, Filter, req1); err == nil {
				// defer cur.Close(ctx)

				req2.Collection = "utxo"
				if len(rr) > 0 {
					if !flags.Concurrent {
						txInputs, errs = GetStxiSeq(rr, &req2)
					} else {
						txInputs, errs = GetStxiSeq(rr, &req2)
					}
					if len(errs) > 0 {
						for _, err := range errs {
							log.Error().Err(err).Msg("Get Tx Input  address")
						}
					}
					break
				} else {
					if retry < flags.MaxRetry {
						log.Warn().Msgf("Filter %v, 0 of %d stxi is returned - retry %d ", Filter, *findOptions.Limit, retry)
					} else {
						log.Error().Msgf("Filter %v, 0 of %d stxi is returned - the number of retries are exceeded %d ", Filter, *findOptions.Limit, retry)
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
		}
	} else {
		err = errors.New("transaction input count is 0")
		LogError(err, fmt.Sprintf("collection %s - Filter %v", req.Collection, Filter))

	}
	return txInputs, err
}

func GetStxiSeq(rr []types.Stxi, req *mongodb.MongoDB) (txInputs []types.TxInput, errs []error) {

	/*
		req.collection is "utxo"
	*/

	var (
		m1   types.Utxo
		opt1 options.FindOneOptions
	)
	req.Collection = "utxo"
	for _, v := range rr {

		var (
			txID    = v.TxInput.TxID
			txInd   = v.TxInput.Index
			txInput types.TxInput
		)
		Filter1 := bson.M{
			"context.tx_hash":    txID,
			"context.output_idx": txInd,
		}
		// fmt.Println("filter", Filter1)
		if m1, err := FindOne(&opt1, m1, Filter1, req); err == nil {

			txInput.Address = m1.TxOutput.Address
			txInput.Amount = m1.TxOutput.Amount
			txInput.TxID = txID
			txInput.Index = txInd
			txInput.InputIdx = int64((v.Context.InputIdx).(float64))
			for _, v1 := range m1.TxOutput.Assets {
				v1.FingerPrint, _ = v1.SetFingerPrint()
				txInput.Assets = append(txInput.Assets, v1)
			}
			// b, _ := json.Marshal(&txInput)
			//fmt.Println("input Index", v.Context.InputIdx)
			// utils.PrintJson(string(b))
			txInputs = append(txInputs, txInput)
		} else {
			if err != mongo.ErrNoDocuments {
				err1 := errors.New(fmt.Sprintf("%s - Filter1 %v -  %v", req.Collection, Filter1, err))
				errs = append(errs, err1)
			}
		}
	}
	return

}

func GetUtxoOutput(flags types.Options, outputCount int, req *mongodb.MongoDB, Filter interface{}) (txOutputs []types.TxOutput, err error) {
	var (
		r           types.Utxo
		rr          []types.Utxo
		findOptions options.FindOptions
		//	cur         *mongo.Cursor
	)
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	req.Collection = "utxo"
	findOptions.SetLimit(int64(outputCount))
	if _, rr, err = Find(&findOptions, r, Filter, req); err == nil {
		// defer cur.Close(ctx)
		for _, v := range rr {
			for i, w := range v.TxOutput.Assets {
				if v.TxOutput.Assets[i].FingerPrint, err = v.SetFingerPrint(w); err != nil {
					log.Warn().Msgf("SetFingerPrint - Filter %v - Error %v", Filter, err)
				}
			}
			txOutputs = append(txOutputs, v.TxOutput)
		}
	}
	return
}

/*
Do not delete
*/

func GetCip25AssetsSeq(v types.Utxo, req *mongodb.MongoDB, Filter interface{}) (map[string]types.Cip25pAsset, []error) {

	var (
		cip25Asset = map[string]types.Cip25pAsset{}
		errs       []error
		err        error
		opt1       options.FindOneOptions
	)

	for i, w := range v.TxOutput.Assets {

		var fingerprint string
		if fingerprint, err = v.SetFingerPrint(w); err != nil {
			log.Warn().Msgf("SetFingerPrint - Filter %v - Error %v", Filter, err)
			errs = append(errs, err)
		} else {
			v.TxOutput.Assets[i].FingerPrint = fingerprint
			Filter1 := bson.M{"cip25_asset.policy": w.Policy, "cip25_asset.asset": w.AssetASCII}
			req.Collection = "cip25p"
			var cip25 types.Cip25p
			if cip25, err := FindOne(&opt1, cip25, Filter1, req); err == nil {
				cip25Asset[fingerprint] = cip25.Cip25Asset
			} else {
				errs = append(errs, err)
			}
		}

	}
	return cip25Asset, errs
}

func GetCip25AssetsCon(v types.Utxo, req *mongodb.MongoDB, Filter interface{}) (map[string]types.Cip25pAsset, []error) {

	type Response struct {
		Index       int
		FingerPrint string
		Cip25Asset  types.Cip25pAsset
		Err         error
	}
	var (
		err        error
		cip25Asset = map[string]types.Cip25pAsset{}
		ch2        = make(chan *Response)
		request    = 0
		receive    = 0
		errs       []error
		opt1       options.FindOneOptions
	)
	request = len(v.TxOutput.Assets)
	for i, w := range v.TxOutput.Assets {
		go func(i int, v types.Utxo, w types.Assets, req *mongodb.MongoDB) {
			var (
				resp        Response
				fingerprint string
			)
			if fingerprint, err = v.SetFingerPrint(w); err != nil {
				log.Warn().Msgf("SetFingerPrint - Filter %v - Error %v", Filter, err)
				resp.Err = errors.New(fmt.Sprintf("SetFingerPrint - Filter %v - Error %v", Filter, err))
			} else {
				resp.FingerPrint = fingerprint
				resp.Index = i
				Filter1 := bson.M{
					"cip25_asset.policy": w.Policy,
					"cip25_asset.asset":  w.AssetASCII,
				}
				req.Collection = "cip25p"
				var cip25 types.Cip25p
				if cip25, err := FindOne(&opt1, cip25, Filter1, req); err == nil {
					resp.Cip25Asset = cip25.Cip25Asset
				}
			}
			ch2 <- &resp
		}(i, v, w, req)
	}

	receive = 0
	if request == 0 {
		return cip25Asset, errs
	}
	for {
		select {
		case rec := <-ch2:
			receive++
			if rec.Err == nil {
				fingerprint := rec.FingerPrint
				v.TxOutput.Assets[rec.Index].FingerPrint = fingerprint
				cip25Asset[fingerprint] = rec.Cip25Asset

			} else {
				LogError(err, "find Cip25")
				errs = append(errs, err)
			}
			if receive == request {
				return cip25Asset, errs
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("a")
		}
	}
}

func GetStxiCon(rr []types.Stxi, req *mongodb.MongoDB) (txInputs []types.TxInput, errs []error) {

	type Response struct {
		TxInput types.TxInput
		Err     error
		K       int
	}
	var (
		m1               types.Utxo
		ch2              = make(chan *Response)
		request, receive = 0, 0
		opt1             options.FindOneOptions
	)
	request = len(rr)
	if request == 0 {
		return
	}
	for _, v := range rr {
		go func(v types.Stxi, req *mongodb.MongoDB, opt1 *options.FindOneOptions) {
			var (
				txID    = v.TxInput.TxID
				txInd   = v.TxInput.Index
				Filter1 = bson.M{"context.tx_hash": txID, "context.output_idx": txInd}
				resp    Response
			)

			if r1, err := FindOne(opt1, m1, Filter1, req); err == nil {
				resp.TxInput.Address = r1.TxOutput.Address
				resp.TxInput.Amount = r1.TxOutput.Amount
				resp.TxInput.TxID = txID
				resp.TxInput.Index = txInd
				for _, v1 := range r1.TxOutput.Assets {
					resp.TxInput.Assets = append(resp.TxInput.Assets, v1)
				}
			} else {
				if err != mongo.ErrNoDocuments {
					err = errors.New(fmt.Sprintf("%s - Filter1 %v -  %v", req.Collection, Filter1, err))
				}
			}

			ch2 <- &resp

		}(v, req, &opt1)
	}

	receive = 0
	for {
		select {
		case rec := <-ch2:
			receive++
			if rec.Err == nil {
				txInputs = append(txInputs, rec.TxInput)
			} else {
				errs = append(errs, rec.Err)
			}
			if receive == request {
				return txInputs, errs
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("s")
		}
	}
}

func GetMeta(req *mongodb.MongoDB, Filter interface{}) (metaDatas []interface{}) {

	req.Collection = "meta"
	var (
		m           types.Meta
		mm          []types.Meta
		findOptions options.FindOptions
	)
	mm, _ = FindAll(&findOptions, m, Filter, req)
	for _, v := range mm {
		meta := v.Metadata
		if meta != nil {
			metaDatas = append(metaDatas, meta)
		}
	}
	return
}

func GetMint(req *mongodb.MongoDB, mintCount int64, Filter interface{}) ([]types.Assets, error) {
	var (
		err         error
		r           types.Mint
		rr          []types.Mint
		mintAsset   []types.Assets
		findOptions options.FindOptions
	)
	findOptions.SetLimit(mintCount)
	if _, rr, err = Find(&findOptions, r, Filter, req); err == nil {
		for _, v := range rr {
			fingerprint, err := v.SetFingerPrint()
			if err == nil {
				v.Asset.FingerPrint = fingerprint
			}
			if v.Asset.AssetASCII == "" {
				assetAscii, err := hex.DecodeString(v.Asset.Asset)
				if err == nil {
					v.Asset.AssetASCII = string(assetAscii)
				}
			}
			mintAsset = append(mintAsset, v.Asset)
		}
		// Trans.TxMeta.MintCount = int64(len(rr))
	} else {
		log.Error().Err(err).Msgf("find collection %s", req.Collection)
		return mintAsset, err
	}
	return mintAsset, err
}

func GetCollateral(req *mongodb.MongoDB, req1 *mongodb.MongoDB, Filter interface{}) (types.Collateral, error) {

	var (
		err        error
		c          types.Coll
		collateral types.Collateral
		opt1       options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	req.Collection = "coll"
	if c, err := FindOne(&opt1, c, Filter, req); err == nil {
		if c.Collateral.TxID != "" {
			var (
				txID  = c.Collateral.TxID
				txInd = c.Collateral.Index
				m2    types.Utxo
			)
			Filter1 := bson.M{
				"context.tx_hash":    txID,
				"context.output_idx": txInd,
			}
			req1.Collection = "utxo"
			if r1, err := FindOne(&opt1, m2, Filter1, req1); err == nil {
				c.Collateral.Address = r1.TxOutput.Address
				c.Collateral.Amount = int(r1.TxOutput.Amount)
			} else {
				log.Warn().Msgf("Err %v - Could not find collateral utxo output  %v for tx_hash %s", err, Filter1, c.Context.TxHash)
			}
			collateral = c.Collateral
		}
	}
	return collateral, err
}

func GetPlutusDatums(req *mongodb.MongoDB, Filter interface{}) ([]types.PlutusDatum, error) {
	var (
		err         error
		t           types.Datum
		tt          []types.Datum
		pDatums     []types.PlutusDatum
		findOptions options.FindOptions
	)
	req.Collection = "dtum"
	findOptions.SetLimit(20)
	if _, tt, err = Find(&findOptions, t, Filter, req); err == nil {
		for _, v := range tt {
			if v.PlutusDatum.DatumHash != "" {
				pDatums = append(pDatums, v.PlutusDatum)
			}
		}
	}
	return pDatums, err
}

func GetPlutusDatas(req *mongodb.MongoDB, Filter interface{}) ([]types.PlutusData, error) {
	var (
		err         error
		t           types.Datum
		tt          []types.Datum
		pDatas      []types.PlutusData
		findOptions options.FindOptions
	)
	req.Collection = "dtum"
	findOptions.SetLimit(20)
	if _, tt, err = Find(&findOptions, t, Filter, req); err == nil {
		for _, v := range tt {
			if v.PlutusDatum.DatumHash != "" {
				var pData types.PlutusData
				if pData, err = utils.PrimitiveDtoStruct(pData, v.PlutusDatum.PlutusData); err == nil {
					pDatas = append(pDatas, pData)
				} else {
					log.Error().Err(err).Msgf("primitiveD to struct")
				}

			}
		}
	}
	return pDatas, err
}

func GetRedeemer(req *mongodb.MongoDB, Filter interface{}) ([]types.PlutusRedeemer, error) {

	var (
		err         error
		t           types.Redeemer
		tt          []types.Redeemer
		redeems     []types.PlutusRedeemer
		findOptions options.FindOptions
	)
	req.Collection = "rdmr"
	findOptions.SetLimit(20)

	if _, tt, err = Find(&findOptions, t, Filter, req); err == nil {
		for _, v := range tt {
			if v.PlutusRedeemer.ExUnitsMem > 0 {
				redeems = append(redeems, v.PlutusRedeemer)
			}
		}
	}
	return redeems, err

}

func GetPlutusWitness(req *mongodb.MongoDB, Filter interface{}) (types.PlutusWitness, error) {

	var (
		err  error
		t    types.Witp
		rpw  types.PlutusWitness
		opt1 options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	if r, err := FindOne(&opt1, t, Filter, req); err == nil {
		if r.PlutusWitness.ScriptHash != "" {
			rpw = r.PlutusWitness
		}
	}
	return rpw, err

}

func GetNativeWitness(req *mongodb.MongoDB, Filter interface{}) (types.NativeWitness, error) {

	var (
		err  error
		t    types.Witn
		wn   types.NativeWitness
		opt1 options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	if w, err := FindOne(&opt1, t, Filter, req); err == nil {
		if w.NativeWitness.PolicyID != "" {
			wn = w.NativeWitness
		}

	}
	return wn, err
}

func GetCip25p(req *mongodb.MongoDB, mintCount int64, Filter interface{}) ([]types.Cip25pAsset, error) {

	var (
		err         error
		c25         types.Cip25p
		cc25        []types.Cip25p
		cip25Assets []types.Cip25pAsset
		findOptions options.FindOptions
	)
	req.Collection = "cip25p"
	findOptions.SetLimit(mintCount)
	if _, cc25, err = Find(&findOptions, c25, Filter, req); err == nil {
		for _, v := range cc25 {
			v.Cip25Asset.FingerPrint, _ = v.SetFingerPrint(v.Cip25Asset)
			cip25Assets = append(cip25Assets, v.Cip25Asset)
		}
	}
	return cip25Assets, err
}

func GetPoolRegistration(req *mongodb.MongoDB, Filter interface{}) (types.PoolRegistration, error) {

	var (
		err  error
		p    types.Pool
		opt1 *options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	req.Collection = "pool"
	// var pn types.Pool
	if p, err = FindOne(opt1, p, Filter, req); err == nil {
		if p.PoolRegistration.PoolId != "" {
			return p.PoolRegistration, err
		}
	}
	return types.PoolRegistration{}, err
}

func GetPoolRetirement(req *mongodb.MongoDB, Filter interface{}) (types.PoolRetirement, error) {

	var (
		err  error
		p    types.Reti
		opt1 options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	req.Collection = "reti"
	// var pn types.Reti
	if p, err = FindOne(&opt1, p, Filter, req); err == nil {
		if p.PoolRetirement.Pool != "" {
			return p.PoolRetirement, err
		}
	}
	return types.PoolRetirement{}, err
}

func GetStakeDelegation(req *mongodb.MongoDB, Filter interface{}) (types.StakeDelegation, error) {
	var (
		err  error
		p    types.Dele
		opt1 options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	req.Collection = "dele"
	if p, err = FindOne(&opt1, p, Filter, req); err == nil {
		if p.StakeDelegation.PoolHash != "" {
			return p.StakeDelegation, err
		}
	}
	return types.StakeDelegation{}, err
}

func GetStakeRegistration(req *mongodb.MongoDB, Filter interface{}) (types.StakeRegistration, error) {
	var (
		err  error
		p    types.Skre
		opt1 options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	req.Collection = "skre"
	if p, err = FindOne(&opt1, p, Filter, req); err == nil {
		if p.StakeRegistration.Credential.AddrKeyhash != "" {
			return p.StakeRegistration, err
		}
	}
	return types.StakeRegistration{}, err
}

func GetStakeDeregistration(req *mongodb.MongoDB, Filter interface{}) (types.StakeDeregistration, error) {
	var (
		err  error
		p    types.Skde
		opt1 *options.FindOneOptions
	)
	opt1.SetAllowPartialResults(false)
	req.Collection = "skde"
	if p, err = FindOne(opt1, p, Filter, req); err == nil {
		if p.StakeDeregistration.Credential.AddrKeyhash != "" {
			return p.StakeDeregistration, err
		}
	}
	return types.StakeDeregistration{}, err
}
