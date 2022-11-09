package db

import (
	"github.com/dgraph-io/badger/v3"
)

func OpenBadgerDB(dir string, logLevel int) (*badger.DB, error) {

	opts := badger.DefaultOptions(dir)
	opts.ValueLogFileSize = 209715200
	opts.BaseLevelSize = 209715200
	opts.SyncWrites = true
	return badger.Open(opts)
}
