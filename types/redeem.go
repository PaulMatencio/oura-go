package types

type Redeemer struct {
	Context        Context        `json:"context" bson:"context"`
	Fingerprint    string         `json:"fingerprint" bson:"fingerprint"`
	PlutusRedeemer PlutusRedeemer `json:"plutus_redeemer" bson:"plutus_redeemer"`
}
type PlutusData struct {
	Constructor int           `json:"constructor" bson:"constructor"`
	Fields      []interface{} `json:"fields" bson:"fields"`
}

type PlutusRedeemer struct {
	ExUnitsMem   int        `json:"ex_units_mem" bson:"ex_units_mem"`
	ExUnitsSteps int        `json:"ex_units_steps" bson:"ex_units_steps"`
	InputIdx     int        `json:"input_idx" bson:"input_idx"`
	PlutusData   PlutusData `json:"plutus_data" bson:"plutus_data"`
	Purpose      string     `json:"purpose" bson:"purpose"`
}

/*
type RedeemerB struct {
	Context        ContextB        `json:"context" bson:"context"`
	Fingerprint    string          `json:"fingerprint" bson:"fingerprint"`
	PlutusRedeemer PlutusRedeemerB `json:"plutus_redeemer" bson:"plutus_redeemer"`
}

type PlutusDataB struct {
	Constructor int           `json:"constructor" bson:"constructor"`
	Fields      []interface{} `json:"fields" bson:"fields"`
}

type PlutusRedeemerB struct {
	ExUnitsMem   int         `json:"ex_units_mem" bson:"ex_units_mem"`
	ExUnitsSteps int         `json:"ex_units_steps" bson:"ex_units_steps"`
	InputIdx     int         `json:"input_idx" bson:"input_idx"`
	PlutusData   PlutusDataB `json:"plutus_data" bson:"plutus_data"`
	Purpose      string      `json:"purpose" bson:"purpose"`
}

*/
