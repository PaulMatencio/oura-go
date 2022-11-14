package types

import (
	"encoding/json"
)

type TransNft struct {
	Context        Context          `json:"context" bson:"context"`
	Fingerprint    string           `json:"fingerprint" bson:"fingerprint"`
	TxNftMeta      TxNftMeta        `json:"tx_meta" bson:"tx_meta"`
	Transaction    Transaction      `json:"transaction" bson:"transaction"`
	TxInput        []TxInput        `json:"utxo_input" bson:"utxo_input"`
	TxOutput       []TxOutput       `json:"utxo_output" bson:"utxo_output"`
	Collateral     Collateral       `json:"collateral,omitempty" bson:"collateral,omitempty"`
	PlutusRedeemer []PlutusRedeemer `json:"plutus_redeemer,omitempty" bson:"plutus_redeemer,omitempty"`
	PlutusData     []PlutusData     `json:"plutus_data,omitempty" bson:"plutus_data,omitempty"`
	Cip25pAsset    Cip25pAsset      `json:"cip25_asset,omitempty" bson:"cip25_asset,omitempty"`
	TransSummary   TransNftSummary  `json:"summary" bson:"summary,omitempty"`
}

type TxNftMeta struct {
	MarketPlace     MktPlace `bson:"market_place" json:"market_place"`
	InputCount      int      `bson:"input_count" json:"input_count"`
	OutputCount     int      `bson:"output_count" json:"output_count"`
	CollCount       int      `bson:"coll_count" json:"coll_count"`
	PlutusDataCount int      `bson:"plutus_data_count" json:"plutus_data_count"`
	PlutusRdmrCount int      `bson:"plutus_redeemer_count" json:"plutus_redeemer_count"`
	Cip25AssetCount int      `bson:"cip25_asset_count" json:"cip25_asset_count"`
}

func (trans *TransNft) Unmarshall(b []byte) error {
	return json.Unmarshal(b, &trans)
}

func (trans *TransNft) SetMarketPlace(mkpl MktPlace) {
	trans.TxNftMeta.MarketPlace = mkpl
}

func (trans *TransNft) GetMarketPlace() MktPlace {
	return trans.TxNftMeta.MarketPlace
}

func (trans *TransNft) SetContext(context Context) {
	trans.Context = context
}

func (trans *TransNft) SetTransaction(tx Transaction) {
	trans.Transaction = tx
}

func (trans *TransNft) SetFingerPrint(fp string) {
	trans.Fingerprint = fp
}

func (trans *TransNft) SetOutputCount(count int) {
	trans.TxNftMeta.OutputCount = count
}

func (trans *TransNft) SetTxOutput(txOutput []TxOutput) {
	trans.TxOutput = txOutput
}

func (trans *TransNft) GetTxOutput() []TxOutput {
	return trans.TxOutput
}

func (trans *TransNft) SetInputCount(count int) {
	trans.TxNftMeta.InputCount = count
}

func (trans *TransNft) SetTxInput(txInput []TxInput) {
	trans.TxInput = txInput
}

func (trans *TransNft) GetTxInput() []TxInput {
	return trans.TxInput
}

func (trans *TransNft) SetCollateralCount(count int) {
	trans.TxNftMeta.CollCount = count
}

func (trans *TransNft) SetCollateral(collateral Collateral) {
	trans.Collateral = collateral
}

/*
func (trans *TransNft) SetDatumCount(count int) {
	trans.TxNftMeta.PlutusDatumCount = count
}
*/

func (trans *TransNft) SetDataCount(count int) {
	trans.TxNftMeta.PlutusDataCount = count
}

/*
func (trans *TransNft) SetPlutusDatum(dt []PlutusDatum) {
	trans.PlutusDatum = dt
}

*/

func (trans *TransNft) SetPlutusData(dt []PlutusData) {
	trans.PlutusData = dt
}

func (trans *TransNft) SetRedeemerCount(count int) {
	trans.TxNftMeta.PlutusRdmrCount = count
}

func (trans *TransNft) SetRedeemer(rd []PlutusRedeemer) {
	trans.PlutusRedeemer = rd
}

/*
func (trans *TransNft) SetPlutusWitnessCount(count int) {
	trans.TxNftMeta.PlutusWitnessesCount = count
}


func (trans *TransNft) SetPlutusWitness(wn PlutusWitness) {
	trans.PlutusWitness = wn
}

/*
func (trans *TransNft) SetNativeWitnessCount(count int) {
	trans.TxNftMeta.NativeWitnessesCount = count
}



func (trans *TransNft) SetNativeWitness(wn NativeWitness) {
	trans.NativeWitness = wn
}

*/

func (trans *TransNft) SetNftSummary(s TransNftSummary) {
	trans.TransSummary = s
}

func (trans *TransNft) SetCip25Asset(s Cip25pAsset) {
	trans.Cip25pAsset = s
}

/*
func (trans *TransNft) SetCip25s(cip25s []Cip25p) {
	for _, v := range cip25s {
		fp := v.Cip25Asset.FingerPrint
		trans.Cip25pAsset[fp] = v.Cip25Asset
	}
}
*/

type TransNftSummary struct {
	MarketPlace      MktPlace      `bson:"market_place" json:"market_place"`
	Timestamp        int64         `bson:"timestamp" json:"timestamp"`
	TxHash           string        `json:"tx_hash" bson:"tx_hash"`
	Purpose          string        `json:"purpose" bson:"purpose"`
	TotalInput       int64         `json:"total_input" bson:"total_input"`
	TotalOutput      int64         `json:"total_output" bson:"total_output"`
	AssetAscii       string        `json:"asset_name" bson:"asset_name"`
	AssetQuantity    int64         `json:"asset_quantity" bson:"asset_quantity"`
	AssetPolicy      string        `json:"asset_policy" bson:"asset_policy"`
	AssetFingerPrint string        `json:"fingerprint" bson:"fingerprint"`
	AssetMintedBy    AssetMintedBy `json:"asset_minted_by" bson:"asset_minted_by"`
	FromAddress      string        `json:"source_address" bson:"source_address"`
	ToAddress        string        `json:"buyer_address" bson:"buyer_address"`
	SellerAddress    string        `json:"seller_address" bson:"seller_address"`
	AssetPrice       int64         `json:"asset_price" bson:"asset_price"`
	AdaBack          int64         `json:"amount_returned" bson:"amount_returned"`
	AdaSpent         int64         `json:"buyer_spent" bson:"buyer_spent"`
	OtherSpent       int64         `json:"other_spent" bson:"other_spent"`
	MarketReceived   int64         `json:"marketplace_fee" bson:"marketplace_fee"`
	OtherReceived    int64         `json:"royalty_payment" bson:"royalty_payment"`
	SellerReceived   int64         `json:"seller_payment" bson:"seller_payment"`
	Fee              int64         `json:"transaction_fee" bson:"transaction_fee"`
}

type AssetMintedBy struct {
	TxHash    string `json:"tx_hash" bson:"tx_hash"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}

/*
type MovedAsset struct {
	AssetAscii       string `json:"asset_name" bson:"asset_name"`
	AssetPolicy      string `json:"asset_policy" bson:"asset_policy"`
	AssetFingerPrint string `json:"fingerprint" bson:"fingerprint"`
	FromAddress      string `json:"source_address" bson:"source_address"`
	ToAddress        string `json:"buyer_address" bson:"buyer_address"`
}

*/

func (s *TransNftSummary) SetMarketPlace(v MktPlace) {
	s.MarketPlace = v
}

func (s *TransNftSummary) SetTimestamp(v int64) {
	s.Timestamp = v
}

func (s *TransNftSummary) GetMarketPlaceName() string {
	return s.MarketPlace.Name
}

func (s *TransNftSummary) SetPurpose(v string) {
	s.Purpose = v
}

func (s *TransNftSummary) SetTotalInput(v int64) {
	s.TotalInput = v
}

func (s *TransNftSummary) SetTxHash(v string) {
	s.TxHash = v
}
func (s *TransNftSummary) SetTotalOutput(v int64) {
	s.TotalOutput = v
}

func (s *TransNftSummary) SetAssetName(v string) {
	s.AssetAscii = v
}

func (s *TransNftSummary) GetAssetName() string {
	return s.AssetAscii
}

func (s *TransNftSummary) SetAssetQuantity(v int64) {
	s.AssetQuantity = v
}

func (s *TransNftSummary) SetAssetMintedBy(v AssetMintedBy) {
	s.AssetMintedBy = v
}

func (s *TransNftSummary) SetAssetPolicy(v string) {
	s.AssetPolicy = v
}
func (s *TransNftSummary) SetAssetFingerPrint(v string) {
	s.AssetFingerPrint = v
}

func (s *TransNftSummary) SetFromAddress(v string) {
	s.FromAddress = v
}

func (s *TransNftSummary) GetFromAddress() string {
	return s.FromAddress
}
func (s *TransNftSummary) GetSourceAddress() string {
	return s.FromAddress
}

func (s *TransNftSummary) SetToAddress(v string) {
	s.ToAddress = v
}

func (s *TransNftSummary) SetAssetPrice(v int64) {
	if s.AssetQuantity > 1 {
		v = v / s.AssetQuantity
	}
	s.AssetPrice = v
}

func (s *TransNftSummary) SetFee(v int64) {
	s.Fee = v
}

func (s *TransNftSummary) SetMarketReceived(v int64) {
	s.MarketReceived = v
}

func (s *TransNftSummary) SetOtherReceived(v int64) {
	s.OtherReceived = v
}

func (s *TransNftSummary) SetSellerReceived(v int64) {
	s.SellerReceived = v
}
func (s *TransNftSummary) SetSellerAddress(v string) {
	s.SellerAddress = v
}

func (s *TransNftSummary) SetAdaSpent(v int64) {
	s.AdaSpent = v
}

func (s *TransNftSummary) SetOtherSpent(v int64) {
	s.OtherSpent = v
}

func (s *TransNftSummary) SetAdaBack(v int64) {
	s.AdaBack = v
}
