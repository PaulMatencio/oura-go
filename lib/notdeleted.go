package lib

/*

	 DO NOT DELETE

	 Build the Plutus Witness segment


	req.Collection = "witp"
	wn, err := GetPlutusWitness(req, Filter)
	if err == nil {
		Trans.SetPlutusWitness(wn)
		Trans.SetPlutusWitnessCount(1)
	} else {
		if err != mongo.ErrNoDocuments {
			log.Error().Err(err).Msgf("collection %s - Filter %v", req.Collection, Filter)
			return Trans, err
		}
	}

*/

/*
		Build the Native Witness segment
	    DO NOT DELETE

	req.Collection = "witn"
	nw, err := GetNativeWitness(req, Filter)
	if err == nil {
		Trans.SetNativeWitness(nw)
		Trans.SetNativeWitnessCount(1)
	} else {
		if err != mongo.ErrNoDocuments {
			log.Error().Err(err).Msgf("collection %s - Filter %v", req.Collection, Filter)
			return Trans, err
		}
	}
*/

/*      DO NOT DELETE

/*
func getPriceJpstore(plutusDatas []types.PlutusData) (price int64, err error) {

	var flt float64
	for _, data := range plutusDatas {
		var dt types.PlutusJpstore
		b, _ := bson.MarshalExtJSON(data, true, true)
		if err = json.Unmarshal(b, &dt); err == nil {
			if len(dt.Fields) > 1 {
				v := dt.Fields[1].Int.NumberDouble
				flt, err = ParseFloat(v)
				if err == nil {
					price = int64(flt)

				}
			} else {
				err = errors.New("could not get the token price")
				log.Warn().Err(err).Msgf("Fields %v", dt.Fields)
			}
		}
	}
	return
}
func getPriceSpacebudz(plutusDatas []types.PlutusData) (price int64, err error) {

	for _, data := range plutusDatas {
		var dt types.PlutusSpacebudz
		b, _ := bson.MarshalExtJSON(data, true, true)
		if err = json.Unmarshal(b, &dt); err == nil {
			if len(dt.Fields) > 0 {
				price := dt.Fields[0].Fields[2].Int
				fmt.Println(price)
			}
		}
	}
	return
}

func getPriceAdapix(plutusDatas []types.PlutusData) (price int64, err error) {
	var pr float64
	for _, data := range plutusDatas {
		b, _ := bson.MarshalExtJSON(data, true, true)
		var dt types.Adapix
		err = json.Unmarshal(b, &dt)
		if err == nil {
			pr, err = ParseFloat(dt.Int.NumberDouble)
			price = int64(pr)
		}
	}
	return
}

func getPriceCnftio(txInput []types.TxInput, fromAddr string) (price int64, err error) {

	for _, v := range txInput {
		if v.Address == fromAddr && len(v.Assets) == 0 {
			price = int64(v.Amount)
			return
		}
	}
	return
}
*/
