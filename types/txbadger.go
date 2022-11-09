package types

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"github.com/paulmatencio/oura-go/db"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strings"
)

/*

 */

type TxBadger struct {
	Fingerprint string `json:"fingerprint" bson:"fingerprint"`
	TxHash      string `json:"tx_hash" bson:"tx_hash"`
	BlockHash   string `json:"block_hash" bson:"block_hash"`
}

func (tx *TxBadger) New(dataDir string) (badgerDB *db.BadgerDB, err error) {

	if dataDir == "" {
		if dataDir, err = homedir.Dir(); err != nil {
			return nil, err
		}
	}
	if !strings.Contains(dataDir, "/") {
		if h, err := homedir.Dir(); err == nil {
			dataDir = filepath.Join(h, dataDir)
		}
	}
	log.Info().Msgf("New badger data base directory %s", dataDir)
	return db.NewBadgerDB(dataDir, nil)
}

func (tx *TxBadger) SetValue(Txb *Tx) {
	tx.TxHash = Txb.Context.TxHash
	tx.Fingerprint = Txb.Fingerprint
	tx.BlockHash = Txb.Context.BlockHash
}

func (tx *TxBadger) Unmarshal() (value []byte, err error) {
	err = json.Unmarshal(value, tx)
	return
}

func (tx *TxBadger) Set(badgerDb *db.BadgerDB, ns []byte, key []byte, value []byte) (err error) {
	return badgerDb.Set(ns, key, value)

}

func (tx *TxBadger) Get(badgerDb *db.BadgerDB, ns []byte, key []byte) (value []byte, err error) {
	return badgerDb.Get(ns, key)
}
