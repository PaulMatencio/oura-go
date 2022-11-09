package lib

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Total  int
	Terror int
}

func GetCip25a(req *mongodb.MongoDB, policy string, asset string) (cip25 types.Cip25p, err error) {
	var (
		Filter1     = bson.M{"cip25_asset.policy": policy, "cip25_asset.asset": asset}
		findOptions options.FindOptions
		r           types.Cip25p
		rr          []types.Cip25p
	)
	findOptions.SetLimit(1)
	req.Collection = "cip25p"
	if _, rr, err = Find(&findOptions, r, Filter1, req); err == nil {
		if len(rr) > 0 {
			return rr[0], nil
		}
	}
	return
}

func GetCip25c(req *mongodb.MongoDB, policy string, asset string) (cip25 types.Cip25p, err error) {
	var (
		opt1    options.FindOneOptions
		Filter1 = bson.M{"cip25_asset.policy": policy, "cip25_asset.asset": asset}
	)
	req.Collection = "cip25p"
	opt1.SetAllowPartialResults(false)
	cip25, err = FindOne(&opt1, cip25, Filter1, req)
	return
}

func GetCip25b(req *mongodb.MongoDB, asset types.Assets) (cip25 types.Cip25p, err error) {
	var (
		opt1    options.FindOneOptions
		Filter1 = bson.M{"cip25_asset.policy": asset.Policy, "cip25_asset.asset": asset.AssetASCII}
	)
	req.Collection = "cip25p"
	opt1.SetAllowPartialResults(false)
	cip25, err = FindOne(&opt1, cip25, Filter1, req)
	return
}

/*

	Build cip25p collection
    called by buildC25Start  ( cmd buildCip25)
    call  BuildC25

*/

func BuildC25Seq(flags types.Options, req *mongodb.MongoDB, req1 *mongodb.MongoDB, req2 *mongodb.MongoDB, Filter interface{}) (Tmetas int, Total int, Terror int, lastId primitive.ObjectID) {

	var (
		err         error
		m           types.Meta
		mm          []types.Meta
		findOptions options.FindOptions
	)
	log.Info().Msg("building cip25 sequentially ")
	lastId = primitive.NilObjectID
	req.Collection = "meta"
	findOptions.SetLimit(flags.Limit)
	if _, mm, err = Find(&findOptions, m, Filter, req); err == nil {

		for _, m1 := range mm {
			resp := BuildC25(flags, req1, req2, &m1)
			Total += resp.Total
			Terror += resp.Terror
		}
		if len(mm) > 0 {
			lastId = mm[len(mm)-1].ID
		} else {
			lastId = primitive.NilObjectID
		}
	}
	Tmetas = len(mm)
	return
}

/*
		Build cip25p collection in parallel
	    called by buildC25Start  ( cmd buildCip25)
	    call  BuildC25
*/
func BuildC25Con(flags types.Options, req *mongodb.MongoDB, req1 *mongodb.MongoDB, req2 *mongodb.MongoDB, Filter interface{}) (Tmetas int, Total int, Terror int, lastId primitive.ObjectID) {

	var (
		err         error
		m           types.Meta
		mm          []types.Meta
		ch          = make(chan *Response)
		findOptions options.FindOptions
	)
	log.Info().Msg("building cip25 in parallel")
	lastId = primitive.NilObjectID
	req.Collection = "meta"

	/*
			loop tru meta with label 721
		    for each meta document create as many as cip25 as there are assets
	*/
	findOptions.SetLimit(flags.Limit)
	if _, mm, err = Find(&findOptions, m, Filter, req); err == nil {

		request := len(mm)
		if request == 0 {
			return
		}
		for _, m1 := range mm {
			go func(m1 types.Meta, req *mongodb.MongoDB, req1 *mongodb.MongoDB, req2 *mongodb.MongoDB) {
				resp := BuildC25(flags, req1, req2, &m1)
				ch <- &resp
			}(m1, req, req1, req2)
		}
		receive := 0
		for {
			select {
			case rec := <-ch:
				receive++
				Total += rec.Total
				Terror += rec.Terror
				if receive == request {
					if len(mm) > 0 {
						lastId = mm[len(mm)-1].ID
					} else {
						lastId = primitive.NilObjectID
					}
					Tmetas = len(mm)
					return
				}
			case <-time.After(100 * time.Millisecond):
				fmt.Printf("n")
			}
		}
	}
	req1.DisConnect()
	req2.DisConnect()
	return
}

/*
		Build cip25p collection
	    is called by BuildC25seq and BuildC25Con
	    calling   buildCP25s
*/
func BuildC25(flags types.Options, req1 *mongodb.MongoDB, req2 *mongodb.MongoDB, m *types.Meta) (resp Response) {

	/*
				req1 is used fot insert
				call BuildCip25s ->  BuildCip25Tokens
			                ->  BuildCip25CH
			                ->  BuildCip25Std
		        req2 is used for accessing collection mint
	*/

	cip25s, _, e := BuildCip25s(flags, req2, m)
	if len(cip25s) == 0 {
		log.Warn().Msgf("tx_hash:%s - no CIP25 - Check the corresponding metadata or run with trace on", m.Context.TxHash)
		log.Trace().Msgf("tx_hash:%s - no CIP25 - Check the corresponding metadata: %v ", m.Context.TxHash, m.Metadata)
	} else {

		resp.Total = len(cip25s)
		resp.Terror = e
		if flags.PrintIt {
			for _, cip25 := range cip25s {
				b, err := json.Marshal(&cip25)
				if err == nil {
					utils.PrintJson(string(b))
				}
			}
		}

		if flags.CheckDup {
			for _, v := range cip25s {
				cip25 := v.(types.Cip25p)
				fmt.Println(cip25.Context.TxHash, cip25.Fingerprint)
			}
		}
		if flags.Upload && len(cip25s) > 0 {
			req1.Collection = "cip25p"
			// insertMany will retry  to insert one by one if duplicate key
			if _, err := req1.InsertMany(cip25s); err != nil {
				resp.Terror++
				log.Error().Err(err).Msg(MeInsertMany)
			}
		}
	}

	return
}

/*
     called by buildC25
	 call BuildCip25Tokens
          BuildCip25CH
          BuildCip25Std
*/

func BuildCip25s(flags types.Options, req *mongodb.MongoDB, m *types.Meta) (cip25s []interface{}, total int, terror int) {
	var (
		err      error
		t        map[string]interface{}
		metadata bson.Raw
	)
	metadata, err = bson.Marshal(m.Metadata)
	if err == nil {
		mapJson := metadata.Lookup("map_json")
		if mapJson.String() != "" {
			err = mapJson.Unmarshal(&t)
			if err == nil {
				/*
					tokens  metadata
				*/
				if _, ok := t["tokens"]; ok {
					cip25s = BuildCip25Tokens(m, t)
					return
				}
				/*
						Cardano Hub metadata
					    req1 will be used to get the correct poliicy Id
				*/
				if _, ok := t["chPolicyId"]; ok {
					cip25s = BuildCip25CH(flags, req, m, t)
					return
				}

				/*
						    cip25  metadata
							K  should be a  NFT policy
					        req1 will be used to access the mint collection in case
					        we need to get a valid asset name
				*/
				cip25s = BuildCip25Std(flags, req, m, t)
				return

			} else {
				log.Error().Err(err).Msgf("%s - tx_hash:%s  - mapJson: %v ", UeMapJson, m.Context.TxHash, t)
				terror++
			}
		} else {
			log.Warn().Msgf("tx_hash:%s - the metadata:%v does not contain a mapJson ", m.Context.TxHash, m.Metadata)
		}
	}
	return
}

/*
	called by buildCP25s
*/

func BuildCip25Tokens(m *types.Meta, t map[string]interface{}) (cip25s []interface{}) {

	var (
		assets Assets
		asset  Asset
	)

	for k, itf := range t {
		v1 := reflect.ValueOf(itf)
		switch v1.Kind() {
		case reflect.Map:
			switch k {
			case "tokens":
				for _, key := range v1.MapKeys() {
					if key.Kind() == reflect.String {
						asset.Asset = key.Interface().(string)
					}
					v2 := reflect.ValueOf(v1.MapIndex(key))
					if v2.Kind() == reflect.Struct {
						asset.Name = Struct2String(v2)
					}

					assets.Asset = append(assets.Asset, asset)

				}
			case "asset":
				var ipfs, url string
				for _, key := range v1.MapKeys() {
					v2 := reflect.ValueOf(v1.MapIndex(key))
					if v2.Kind() == reflect.Struct {
						key1 := key.Interface().(string)
						switch key1 {
						case "assetDescription":
							assets.Description = v2.Interface()
						case "mimeType":
							assets.MediaType = Struct2String(v2)
						case "url":
							url = Struct2String(v2)
						case "ipfs":
							ipfs = Struct2String(v2)
						default:
							//fmt.Println(key.Kind(), v2.Kind())
						}
					}
				}
				assets.Image = url + ipfs
				assets.RawJson = itf
				// assets.Asset = append(assets.Asset, asset)

			case "policyScript":
			case "artistName":
			default:
			}
		case reflect.String:
			switch k {
			case "policyId":
				assets.Policy = v1.Interface().(string)
			case "cnftFormatVersion":
				assets.Version = v1.Interface().(string)
			default:
			}
		case reflect.Slice:
			switch k {
			case "publisher":
			default:
			}

		default:

		}
	}

	for _, asst := range assets.Asset {
		var (
			cip25      types.Cip25p
			cip25Asset types.Cip25pAsset
		)
		cip25Asset.Policy = assets.Policy
		cip25Asset.Asset = asst.Asset
		cip25Asset.Name = asst.Name
		cip25Asset.RawJSON = assets.RawJson
		cip25Asset.Image = assets.Image
		cip25Asset.MediaType = assets.MediaType
		cip25Asset.Description = assets.Description
		cip25Asset.Version = assets.Version
		cip25Asset.FingerPrint, _ = cip25Asset.SetFingerPrint()
		fingerPrint := Hash(cip25Asset.Policy + hex.EncodeToString([]byte(cip25Asset.Asset)) + m.Context.TxHash)
		cip25.Fingerprint = strconv.FormatInt(m.Context.Slot, 10) + ".cip25." + fingerPrint
		cip25.Context = m.Context
		cip25.Cip25Asset = cip25Asset
		cip25s = append(cip25s, cip25)
	}
	return
}

// Cardano hub.io
/*
	called by buildCP25s
*/

func BuildCip25CH(flags types.Options, req *mongodb.MongoDB, m *types.Meta, t map[string]interface{}) (cip25s []interface{}) {
	var (
		chStruct   types.ChMapJSON
		rawJson    map[string]interface{}
		cip25      types.Cip25p
		cip25Asset types.Cip25pAsset
	)

	cip25.Context = m.Context
	if jsonString, err := json.Marshal(t); err == nil {
		json.Unmarshal(jsonString, &chStruct)
		cip25Asset.Asset = chStruct.Name
		//  get the PolicyId from the corresponding  mint  of the transaction
		policyId, err := GetPolicyID(req, flags, m)
		if err == nil {
			if policyId != "" {
				cip25Asset.Policy = policyId
			} else {
				log.Error().Err(errors.New("policyID is null")).Msgf("getPolicyId - tx-hash %s", m.Context.TxHash)
				return
			}
		} else {
			log.Error().Err(err).Msgf("getPolicyId - tx-hash %s", m.Context.TxHash)
			return
		}
		cip25Asset.Name = chStruct.Name
		cip25Asset.Website = chStruct.Website
		cip25Asset.FingerPrint, _ = cip25Asset.SetFingerPrint()
		cip25Asset.Description = chStruct.Description
		cip25Asset.Image = chStruct.Image
		cip25Asset.MediaType = chStruct.MediaType
		err = json.Unmarshal([]byte(jsonString), &rawJson)
		if err == nil {
			cip25Asset.RawJSON = rawJson
		} else {
			log.Error().Err(err).Msg(UeJsonString)
		}
	}
	cip25.Cip25Asset = cip25Asset
	// fingerPrint := Hash(cip25Asset.Policy + cip25Asset.Asset + m.Context.TxHash)
	fingerPrint := Hash(cip25Asset.Policy + hex.EncodeToString([]byte(cip25Asset.Asset)) + m.Context.TxHash)
	cip25.Fingerprint = strconv.FormatInt(m.Context.Slot, 10) + ".cip25." + fingerPrint
	cip25s = append(cip25s, cip25)

	return
}

/*
	called by buildCP25s
*/

func BuildCip25Std(flags types.Options, req *mongodb.MongoDB, m *types.Meta, t map[string]interface{}) (cip25s []interface{}) {

	var (
		assets         Assets
		mAsset, mImage string
		RawJson        map[string]interface{}
	)

	for k, v := range t {
		var (
			asset Asset
			// rawJson interface{}
		)
		asset.Policy = k /*    */
		if len(asset.Policy) > 56 {
			asset.Policy = asset.Policy[:56]
		}
		v1 := reflect.ValueOf(v)
		switch v1.Kind() {
		case reflect.Map:
			for _, k2 := range v1.MapKeys() {
				v2 := v1.MapIndex(k2)
				val2 := reflect.ValueOf(v2.Interface())
				// fmt.Println(val2.Kind(), k2, val2)
				switch val2.Kind() {
				case reflect.Map:
					asset.RawJSON = val2.Interface()
					asset.Asset = (k2.Interface()).(string)
					asset.Source = "map"
					fingerPrint := Hash(k + hex.EncodeToString([]byte(asset.Asset)) + m.Context.TxHash)
					asset.Fingerprint = strconv.FormatInt(m.Context.Slot, 10) + ".cip25." + fingerPrint
					for _, k3 := range val2.MapKeys() {
						v3 := val2.MapIndex(k3)
						val3 := reflect.ValueOf(v3.Interface())
						//	tv3 := reflect.TypeOf(v3.Interface()).String()
						key3 := k3.Interface().(string)
						switch strings.ToLower(key3) {
						case "image":
							if val3.Kind() == reflect.String {
								asset.Image = v3.Interface().(string)
							} else {
								asset.Image = v3.Interface()
							}
						case "name":
							/*if tv3 == "string" { */
							if val3.Kind() == reflect.String {
								asset.Name = v3.Interface().(string)
							} else {
								asset.Name = v3.Interface()
							}
						case "description":
							var description []string
							if val3.Kind() == reflect.String {
								description = append(description, v3.Interface().(string))
							} else {
								asset.Name = v3.Interface()
							}
							asset.Description = description
						case "copyright", "copyrights":
							if val3.Kind() == reflect.String {
								asset.CopyRight = v3.Interface().(string)
							} else {
								asset.CopyRight = v3.Interface()
							}
						case "publisher":
							{
								if val3.Kind() == reflect.String {
									asset.Publisher = v3.Interface().(string)
								} else {
									asset.Publisher = v3.Interface()
								}
							}
						case "mediatype", "mime":
							if val3.Kind() == reflect.String {
								asset.MediaType = v3.Interface().(string)
							} else {
								if val3.Kind() == reflect.Slice {
									var n []string
									mapstructure.Decode(v3.Interface(), &n)
									if len(n) > 0 {
										asset.MediaType = n[0]
									}
								} else {
									log.Fatal().Msgf("tx_hash:%s - val3.kind(): %s - k2:%s - v3:%v ", m.Context.TxHash, val3.Kind(), k2, v3.Interface())
								}
							}

						case "website", "url":
							if val3.Kind() == reflect.String {
								asset.Website = v3.Interface().(string)
							} else {
								asset.Website = v3.Interface()
							}
						case "project":
							if val3.Kind() == reflect.String {
								asset.Project = v3.Interface().(string)
							} else {
								asset.Project = v3.Interface()
							}

						case "version":
							if val3.Kind() == reflect.Float64 {
								strconv.Itoa(int(v3.Interface().(float64)))
								asset.Version = strconv.Itoa(int(v3.Interface().(float64)))
							}
							if v3.Kind() == reflect.String {
								asset.Version = v3.Interface().(string)
							}
						default:
							log.Trace().Msgf("tx_hash:%s - val2.Kind():%v - val3.kind():%v - key:%s - v2:%v ", m.Context.TxHash, val2.Kind(), val3.Kind(), k2, v3.Interface())
						}
					}
					//  check if policy is present and valid before appending the asset
					if len(asset.Policy) == 56 {
						assets.Asset = append(assets.Asset, asset)
					}

				case reflect.Slice:
					// v2 is a Slice
					if k2.Interface().(string) == "attributes" {
						break
					}
					// val2 := reflect.ValueOf(v2.Interface())
					for i := 0; i < val2.Len(); i++ {
						v3 := val2.Index(i)
						val3 := reflect.ValueOf(v3.Interface())
						asset.RawJSON = val3.Interface()
						switch val3.Kind() {
						case reflect.Map:
							for _, k3 := range val3.MapKeys() {
								v4 := val3.MapIndex(k3)
								reflect.ValueOf(v4.Interface())
								// asset.Asset = k2.Interface().(string)
								//fingerPrint := Hash(k + hex.EncodeToString([]byte(asset.Asset)) + m.Context.TxHash)
								//asset.Fingerprint = strconv.FormatInt(m.Context.Slot, 10) + ".cip25." + fingerPrint
								switch strings.ToLower(k3.Interface().(string)) {
								case "image":
									asset.Image = v4.Interface().(string)
								case "name":
									if reflect.ValueOf(v4.Interface()).Kind() == reflect.String {
										asset.Name = v4.Interface().(string)
									} else {
										asset.Name = v4.Interface()
										// log.Fatal().Msgf("tx_hash %s", m.Context.TxHash)
									}
									asset.Asset = k2.Interface().(string)
									fingerPrint := Hash(k + hex.EncodeToString([]byte(asset.Asset)) + m.Context.TxHash)
									asset.Fingerprint = strconv.FormatInt(m.Context.Slot, 10) + ".cip25." + fingerPrint
								case "publisher":
									asset.Publisher = v4.Interface()
								}
							}
							if len(asset.Policy) == 56 {
								assets.Asset = append(assets.Asset, asset)
							}
						// case reflect.String:

						default:
							log.Warn().Msgf("tx_hash:%s - K2: %s - val2 Kind:%v - val3 Kind:%v ", m.Context.TxHash, k2.Interface().(string), reflect.TypeOf(v2.Interface()).Kind(), reflect.TypeOf(v3.Interface()).Kind())
						}
						//fmt.Println(v3.Kind(), reflect.TypeOf(v3.Interface()).Kind())
					}
				case reflect.String:
					k2s := k2.Interface().(string)
					switch strings.ToLower(k2s) {
					case "description":
						var description []string
						if v2.Kind() == reflect.String {
							description = append(description, v2.Interface().(string))
						}
						if v2.Kind() == reflect.Slice {
							for i := 0; i < v2.Len(); i++ {
								description = append(description, v2.Index(i).Interface().(string))
							}
						}
						asset.Description = description
					case "name", "asset":
						mAsset = v2.Interface().(string)
						RawJson = t
					case "image":
						mImage = v2.Interface().(string)
					default:
						log.Trace().Msgf("tx_hash:%s - v2.Kind():%v - v2:%v", m.Context.TxHash, val2.Kind(), val2)
					}
				default:
					log.Trace().Msgf("tx_hash:%s - v2.Kind():%v - v2:%v", m.Context.TxHash, val2.Kind(), val2)
				}
			}
			// fmt.Println(mAsset, mImage, len(asset.Asset), asset.Asset)
			if mAsset != "" && len(asset.Asset) == 0 {
				asset.Asset = mAsset
				asset.Name = mAsset
				asset.Policy = k
				fingerPrint := Hash(k + hex.EncodeToString([]byte(asset.Asset)) + m.Context.TxHash)
				asset.Fingerprint = strconv.FormatInt(m.Context.Slot, 10) + ".cip25." + fingerPrint
				asset.RawJSON = RawJson
				if mImage != "" {
					asset.Image = mImage
				}
				if len(asset.Policy) == 56 {
					assets.Asset = append(assets.Asset, asset)
				}
			}
		case reflect.Slice:

			for i := 0; i < v1.Len(); i++ {
				v2 := v1.Index(i)
				val2 := reflect.ValueOf(v2.Interface())
				// fmt.Println("kx:", kx, m.Context.TxHash)
				// switch vx.Kind().String() {
				switch val2.Kind() {

				case reflect.Map:
					// asset.RawJSON = val2.Interface()
					// asset.Asset = k2
					//fmt.Println("....  isSlice of Map", m.Context.TxHash)
					fingerPrint := Hash(k + hex.EncodeToString([]byte(asset.Asset)) + m.Context.TxHash)
					asset.Fingerprint = strconv.FormatInt(m.Context.Slot, 10) + ".cip25." + fingerPrint
					for _, k2 := range val2.MapKeys() {
						asset.Asset = (k2.Interface()).(string)
						asset.Source = "map"
						v3 := val2.MapIndex(k2)
						val3 := reflect.ValueOf(v3.Interface())
						asset.RawJSON = v3.Interface()
						switch val3.Kind() {
						case reflect.Map:
							for _, k3 := range val3.MapKeys() {
								// fmt.Println(k3, val3.MapIndex(k3))
								val4 := val3.MapIndex(k3)
								switch strings.ToLower(k3.Interface().(string)) {
								case "image":
									if reflect.ValueOf(val4.Interface()).Kind() == reflect.String {
										asset.Image = val4.Interface().(string)
									} else {
										asset.Image = val4.Interface()
									}
								case "name":
									if reflect.ValueOf(val4.Interface()).Kind() == reflect.String {
										asset.Name = val4.Interface().(string)
									} else {
										asset.Name = val4.Interface()
									}
								case "mediatype":
									if reflect.ValueOf(val4.Interface()).Kind() == reflect.String {
										asset.MediaType = val4.Interface().(string)
									}
								case "publisher":
									if reflect.ValueOf(val4.Interface()).Kind() == reflect.String {
										asset.Publisher = val4.Interface().(string)
									} else {
										asset.Publisher = val4.Interface()
									}
								case "version":
									if val4.Kind() == reflect.Float64 {
										strconv.Itoa(int(val4.Interface().(float64)))
										asset.Version = strconv.Itoa(int(val4.Interface().(float64)))
									}
									if val4.Kind() == reflect.String {
										asset.Version = val4.Interface().(string)
									}
								case "description":
									var description []string
									if val4.Kind() == reflect.String {
										description = append(description, val4.Interface().(string))
									}
									if val4.Kind() == reflect.Slice {
										for i := 0; i < val4.Len(); i++ {
											description = append(description, val4.Index(i).Interface().(string))
										}
									}
									asset.Description = description
								}
							}
						default:
							log.Warn().Msgf("tx_hash:%s - k2:%s - val3.kind():%v ", m.Context.TxHash, k2, val3.Kind())
						}

					}
					// there are some bug in the metadata .

					if len(asset.Policy) == 56 {
						assets.Asset = append(assets.Asset, asset)
					}

				//case value is a string:
				case reflect.String:
					val2 := v2.Interface().(string)
					switch strings.ToLower(k) {
					case "publisher":
						assets.Publisher = val2
					case "description":
						assets.Description = val2
					default:
						log.Trace().Msgf("tx_hash: %s - v2.Kind(): %v - Key: %s  ", m.Context.TxHash, v2.Kind(), k)
					}
				default:
					log.Trace().Msgf("tx_hash:%s - v2.Kind():%v - Key: %s", m.Context.TxHash, v2.Kind())
				}
			}

		case reflect.String:
			val1 := v1.Interface().(string)
			switch strings.ToLower(k) {
			case "description":
				assets.Description = val1
			case "publisher":
				assets.Publisher = val1
			case "version":
				assets.Version = val1
			case "policyLink":
				assets.PolicyLink = val1
			case "copyright":
				assets.CopyRight = val1
			case "artist":
				assets.Artist = val1
			default:
				log.Trace().Msgf("tx_hash:%s - v1.Kind(): %v - Bytes %v", m.Context.TxHash, v1.Kind(), v1.Interface())
			}
		case reflect.Float64:
			switch strings.ToLower(k) {
			case "version":
				assets.Version = strconv.Itoa(int(v1.Interface().(float64)))
			default:
				log.Info().Msgf("tx_hash:%s - v1.Kind(): %v - Bytes %v", m.Context.TxHash, v1.Kind(), v1.Interface())
			}
		default:
			log.Trace().Msgf("tx_hash:%s - v1.Kind():%v - Bytes:%v ", m.Context.TxHash, v1.Kind(), v1)
		}
	}

	var (
		cip25      types.Cip25p
		cip25Asset types.Cip25pAsset
	)

	cip25.Context = m.Context
	for _, asst := range assets.Asset {

		if asst.Source == "map" {
			cip25Asset.Asset = asst.Asset
			if len(asst.Asset) == 0 {
				log.Warn().Msgf("Check asset name:%s - tx_hash:%s", asst.Asset, m.Context.TxHash)
				if asst1, err := GetAssetName(req, flags, m); err == nil {
					cip25Asset.Asset = asst1
				} else {
					log.Error().Err(err).Msg("GetAssetName")
				}
			}
		} else {
			if len(asst.Asset) < 3 && len(asst.Fingerprint) > 0 {
				log.Warn().Msgf("Invalid asset name ? in metadata: %s - The asset name will be taken from its mint - tx_hash %s", asst.Asset, m.Context.TxHash)
				if asst1, err := GetAssetName(req, flags, m); err == nil {
					cip25Asset.Asset = asst1
				} else {
					log.Error().Err(err).Msg("GetAssetName")
				}
			}
		}

		// cip25Asset.Asset = asst.Asset
		cip25Asset.Name = asst.Name
		cip25Asset.Policy = asst.Policy
		cip25Asset.FingerPrint, _ = cip25Asset.SetFingerPrint()
		cip25Asset.Image = asst.Image
		cip25Asset.MediaType = asst.MediaType
		cip25Asset.Publisher = asst.Publisher

		if asst.PolicyScript != nil {
			cip25Asset.PolicyScript = asst.PolicyScript
		} else {
			cip25Asset.PolicyScript = assets.PolicyScript
		}
		if asst.Description != nil {
			cip25Asset.Description = asst.Description
		} else {
			cip25Asset.Description = assets.Description
		}
		if asst.Artist != "" {
			cip25Asset.Artist = asst.Artist
		} else if assets.Artist != "" {
			cip25Asset.Artist = assets.Artist
		}
		if asst.Publisher == nil {
			cip25Asset.Publisher = assets.Publisher
		}

		if asst.Website != nil {
			cip25Asset.Website = asst.Website
		}

		if asst.Version != "" {
			cip25Asset.Version = asst.Version
		} else if assets.Version != "" {
			cip25Asset.Version = assets.Version
		}

		if asst.PolicyLink != "" {
			cip25Asset.PolicyLink = asst.PolicyLink
		} else if assets.PolicyLink != "" {
			cip25Asset.PolicyLink = assets.PolicyLink
		}

		if asst.CopyRight != nil {
			cip25Asset.Copyright = asst.CopyRight
		} else if assets.CopyRight != "" {
			cip25Asset.Copyright = assets.CopyRight
		}
		cip25Asset.RawJSON = asst.RawJSON
		cip25.Cip25Asset = cip25Asset
		cip25.Fingerprint = asst.Fingerprint
		if cip25.Fingerprint != "" {
			cip25s = append(cip25s, cip25)
		} else {
			log.Warn().Msgf("tx_hash %s -  missing CIP25 fingerprint", m.Context.TxHash)
		}
	}

	return
}
