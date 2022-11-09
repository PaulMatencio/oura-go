package types

type ChMapJSON struct {
	ChRoyaltyAddr string `json:"chRoyaltyAddr" bson:"chRoyaltyAddr""`
	ChFee         int    `json:"chFee" bson:"chFee"`
	ChListingMode string `json:"chListingMode" bson:"chListingMode"`
	ChPrice       int    `json:"chPrice" bson:"chPrice"`
	ChRoyalty     int    `json:"chRoyalty" bson:"chRoyalty"`
	ChSellerAddr  string `json:"chSellerAddr" bson:"chSellerAddr"`
	Name          string `json:"name" bson:"name"`
	Wallet        string `json:"wallet" bson:"wallet"`
	Website       string `json:"website" bson:"website" `
	ChID          string `json:"chId" bson:"chId"`
	ChPolicyID    string `json:"chPolicyId" bson:"chPolicyId"`
	Image         string `json:"image" bson:"image"`
	MediaType     string `json:"mediaType" bson:"mediaType" `
	ChAssetName   string `json:"chAssetName" bson:"chAssetName"`
	Description   string `json:"description" bson:"description"`
	MintedBy      string `json:"mintedBy" bson:"mintedBy"`
	Collection    string `json:"collection" bson:"collection"`
}
