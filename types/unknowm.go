package types

import (
	"encoding/json"
	"github.com/paulmatencio/oura-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Unknown struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Context     Context            `json:"context" bson:"context"`
	Fingerprint string             `json:"fingerprint,omitempty" bson:"fingerprint,omitempty"`
	Unknown     interface{}        `json:"unknown",bson:"unknown"`
}

func (u *Unknown) GetType(v string) (typ string, err error) {
	err = json.Unmarshal([]byte(v), &u)
	if err == nil {
		if u.Fingerprint != "" {
			typ, err = utils.GetType(u.Fingerprint)
		}
	}
	return
}
