package types

type Epoch struct {
	Epoch  int    `bson:"epoch" json:"epoch"`
	Blocks []Blck `bson:"blocks" json:"blocs"`
}
