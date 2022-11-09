package cmd

import (
	"context"
	"eagain.net/go/bech32"
	"encoding/hex"
	"encoding/json"
	"github.com/paulmatencio/oura-go/db"
	"github.com/rs/zerolog"
	"os"
	"path"
	"path/filepath"
	"strings"

	//"encoding/json"
	"errors"
	"fmt"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"time"
)

/*
Open badger db
*/

func OpenBdb(dataDir string) (myBdb *db.BadgerDB, err error) {

	if len(dataDir) == 0 || !strings.Contains(dataDir, "/") {
		if home, err := os.UserHomeDir(); err != nil {
			LogError(err, "get user home directory")
		} else {
			dataDir = filepath.Join(home, dataDir)
		}
	}
	return db.NewBadgerDB(dataDir, nil)
}

func InitCopy(mongoUrl string) (*mongodb.MongoCopy, error) {
	var (
		err    error
		option = options.ClientOptions{
			// SocketTimeout: &timeOut,
		}
		req = mongodb.MongoCopy{
			Option: &option,
			//FromDatabase: database,
		}
		// oplog map[string]string
	)

	req.Uri = "mongodb://" + mongoUrl + "/?ssl=" + mongoSSL
	if req.Client, err = req.Connect(); err != nil {
		LogError(err, "mongodb connect")
		// log.Error().Err(err).Msg("mongodb connection error")
	}
	return &req, err
}

func ListAll[T any](result []T, filter string, req *mongodb.MongoDB) {

	var (
		// result1 []T
		cur     *mongo.Cursor
		err     error
		nerr    int64
		Filter  interface{}
		filters types.Filters
	)

	if filter != "" {
		fmt.Println(filter)
		opVal := filters.ValidOp()
		if filter != "" {
			if filters, err = filters.ParseOp(filter, opVal); err != nil {
				LogError(err, "parse operators")
				//log.Error().Err(err).Msg("parse operators")
				return
			}
		}
		Filter = filters.BuildFilter()
	} else {
		Filter = bson.D{{}}
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	coll := req.Client.Database(req.Database).Collection(req.Collection)
	cur, err = coll.Find(ctx, Filter)
	if cur != nil {
		for cur.Next(ctx) {
			var res T
			if err = cur.Decode(&res); err == nil {
				result = append(result, res)
			} else {
				nerr++
			}
		}
	}
	fmt.Printf("Total number of document:%d\nNumber of error %d\n", len(result), nerr)
	return

}
func getPoolId(x interface{}) (poolId string, err error) {

	switch x.(type) {
	case types.Pool:
		t := x.(types.Pool)
		poolId, err = t.GetPoolId()
	default:
	}
	return
}

func GetPoolMeta(client *http.Client, url string, retryNumber int, waitTime time.Duration) (*types.PoolMeta, error) {
	var (
		poolMeta types.PoolMeta
		err      error
		response *http.Response
	)
	for i := 1; i <= retryNumber; i++ {
		if response, err = client.Get(url); err == nil {
			if response.StatusCode == 200 {
				defer response.Body.Close()
				if contents, err := ioutil.ReadAll(response.Body); err == nil {
					err = json.Unmarshal(contents, &poolMeta)
				}
			} else {
				err = errors.New(fmt.Sprintf("Status: %d %s", response.StatusCode, response.Status))
			}
			break
		} else {
			LogError(err, fmt.Sprintf("number of retries: %d", i))
			//log.Error().Err(err).Msgf("number of retries: %d", i)
			time.Sleep(waitTime * time.Millisecond)
		}
	}
	return &poolMeta, err
}

func PrimitiveDtoMap(in interface{}) (Map map[string]interface{}, err error) {
	var b []byte
	b, err = bson.MarshalExtJSON(in, true, true)
	if err == nil {
		err = json.Unmarshal(b, &Map)
	} else {
		LogError(err, "unmarshal mongodb primitive to map")
		// log.Error().Err(err).Msg("unmarshal mongodb primitive to map")
	}
	return
}

func SetPoolId(pool string) (poolId string, err error) {

	v, err := hex.DecodeString(pool)
	if err == nil {
		poolId, err = bech32.Encode("pool", v)
	}
	return
}

func SetStakeAddress(addrHash string, network string) (stakeAddr string, err error) {
	ex := "e1"
	if network != "mainnet" {
		ex = "e0"
	}
	v, err := hex.DecodeString(ex + addrHash)
	if err == nil {
		stakeAddr, err = bech32.Encode("stake", v)
	}
	return
}

func SetFlags() (flags types.Options) {
	flags.Reg = reg
	flags.Concurrent = concurrent
	flags.Limit = limit
	flags.MaxRetry = MaxRetry
	flags.Upload = upload
	flags.MongoDatabase = database
	flags.MongoSSL = mongoSSL
	flags.PrintIt = printIt
	flags.CheckDup = checkDup
	flags.MongoOptions = &clientOption
	return
}
func SetLogFile(filename string) (log zerolog.Logger, err error) {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().Timestamp().Caller().Logger()
	if filename != "" {
		logDir := filepath.Dir(logger)
		if _, err = os.Stat(logDir); os.IsNotExist(err) {
			if err = os.MkdirAll(logDir, 0755); err != nil {
				return log, err
			}
		}
		if file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664); err == nil {
			log = zerolog.New(file).With().Timestamp().Caller().Logger()
			return log, err
		}
	}

	return log, err
}

func LogError(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func GetBadgerPath(what string) (pathName string) {
	switch what {
	case "nft":
		pathName = path.Join(badgerRoot, badgerNft)
	case "trans":
		pathName = path.Join(badgerRoot, badgerTran)
	case "block":
		pathName = path.Join(badgerRoot, badgerBlock)
	case "pool":
		pathName = path.Join(badgerRoot, badgerPool)
	default:
		pathName = path.Join(badgerRoot, badgerDefault)
	}
	return
}
