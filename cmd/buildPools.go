/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"github.com/paulmatencio/oura-go/lib"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"time"

	// "github.com/hashicorp/go-retryablehttp"
	"github.com/paulmatencio/oura-go/mongodb"
	"github.com/paulmatencio/oura-go/types"
	"github.com/rs/zerolog/log"
)

var buildPoolsCmd = &cobra.Command{
	Use:   "buildPool",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   buildPool,
}

const (
	ConnectTimeout = 10 * time.Second
	RequestTimeout = 30 * time.Second
	KeepAlive      = 15 * time.Second
	RetryNumber    = 5
	WaitTime       = 2 * time.Second
)

func init() {
	rootCmd.AddCommand(buildPoolsCmd)
	initBuildPool(buildPoolsCmd)
	if database == "" {
		database = mongoDatabase
	}
	if database == "" {
		log.Fatal().Msg("mongodb database is missing")
	}

}

func initBuildPool(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().BoolVarP(&printIt, "print", "p", true, "print the result")
}

func initBuildPools(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&mongoUrl, "mongo-url", "H", "localhost:27017", "mongodb host:port")
	cmd.Flags().StringVarP(&database, "mongo-db", "d", "cardano", "mongodb database name")
	cmd.Flags().StringVarP(&dataDir, "badger-db", "b", "", "local database directory. If no a full path, it with be prefix by home directory")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter")
	cmd.Flags().Int64VarP(&limit, "limit", "", 10, "max returned")
	cmd.Flags().BoolVarP(&printIt, "print", "p", false, "print the result")
	cmd.Flags().BoolVarP(&concurrent, "concurrent", "C", false, "Concurrent upload")
	cmd.Flags().IntVarP(&bulk, "bulk", "n", 10, "max bulk load")
}

func buildPool(cmd *cobra.Command, args []string) {

	if log.Logger, err = SetLogFile(logger); err != nil {
		log.Warn().Msgf(" logging to file %s - error %v ", logger, err)
	}
	flags := SetFlags()
	if req, Filter, err := lib.InitReq(flags, mongoUrl, filter); err == nil {
		BuildPool(req, Filter)
	} else {
		log.Error().Stack().Msgf("InitBuild %v", err)
	}
}

func BuildPool(req *mongodb.MongoDB, Filter interface{}) {
	var (
		pool       types.Pool
		findOption *options.FindOptions
	)
	req.Collection = "pool"
	rr, nErr := lib.FindAll(findOption, pool, Filter, req)
	if nErr > 0 {
		log.Warn().Msgf("FindAll pool returned %d errors", nErr)
	}
	for _, v := range rr {
		if poolN, err := BuildPool1(v, req, Filter); err == nil {
			if b, err := json.Marshal(poolN); err == nil {
				utils.PrintJson(string(b))
			} else {
				log.Error().Stack().Err(err).Msgf("BuildPool1 returned with error:  %v", err)
			}
		}
	}

	return
}

func BuildPool1(pool types.Pool, req *mongodb.MongoDB, Filter interface{}) (pooln types.PoolN, err error) {

	var (
		url       = pool.PoolRegistration.PoolMetadata
		client    = &http.Client{}
		transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(ConnectTimeout) * time.Millisecond, // connection timeout
				KeepAlive: time.Duration(KeepAlive) * time.Millisecond,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxConnsPerHost:       100,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		poolOwners  []string
		Dele        types.Dele
		Skre        types.Skre
		req1        = req
		findOptions *options.FindOptions
		opt1        *options.FindOneOptions
	)
	client.Transport = transport
	pooln.CopyFrom(&pool)
	poolOwners = pool.PoolRegistration.PoolOwners
	/*
		check if pool is not retired
	*/
	var Reti types.Reti
	req1.Collection = "reti"
	Filter1 := bson.M{"pool_retirement.pool": pool.PoolRegistration.Operator}
	p, err := lib.FindOne(opt1, Reti, Filter1, req1)
	if err == mongo.ErrNoDocuments {
		if len(url) > 0 {
			poolMeta, err := GetPoolMeta(client, url, RetryNumber, WaitTime)
			if err == nil {
				pooln.PoolRegistration.PoolMeta = *poolMeta
			}
		}

	} else {
		var ctxm types.ContextN
		ctxm.CopyFrom(&p.Context)
		pooln.PoolRetirement.Context = ctxm
		pooln.PoolRetirement.Epoch = p.PoolRetirement.Epoch
		pooln.PoolRetirement.Pool = p.PoolRetirement.Pool
		poolId, err := SetPoolId(p.PoolRetirement.Pool)
		if err == nil {
			pooln.PoolRetirement.PoolId = poolId
		}
	}
	/*

	 */
	pooln.SetPoolOwner(network)
	pooln.SetRewardAccount()
	/*
		Get stake pool delegation history
	*/
	req1.Collection = "dele"
	Filter1 = bson.M{"stake_delegation.pool_hash": pooln.PoolRegistration.Operator}
	dd, nErr := lib.FindAll(findOptions, Dele, Filter1, req1)
	if nErr > 0 {
		log.Warn().Msgf("Find stake delegation returned %d errors", nErr)
	}
	var del types.StakeDelegationN
	for _, d := range dd {
		var ctxm types.ContextN
		ctxm.CopyFrom(&d.Context)
		del.Context = ctxm
		del.PoolHash = d.StakeDelegation.PoolHash
		del.PoolId, _ = SetPoolId(del.PoolHash)
		del.Credential.AddrKeyhash = d.StakeDelegation.Credential.AddrKeyhash
		del.Credential.StakeKey, _ = SetStakeAddress(d.StakeDelegation.Credential.AddrKeyhash, network)
		pooln.StakeDelegation = append(pooln.StakeDelegation, del)
	}
	//
	req1.Collection = "skre"
	for _, owner := range poolOwners {
		Filter1 = bson.M{"stake_registration.credential.addrkey_hash": owner}
		oo, nErr := lib.FindAll(findOptions, Skre, Filter1, req1)
		if nErr > 0 {
			log.Warn().Msgf("Find stake registration returned %d errors", nErr)
		} else {
			b, err := json.Marshal(&oo)
			if err == nil {
				utils.PrintJson(string(b))
			}
		}
	}

	req1.Collection = "skde"
	for _, owner := range poolOwners {
		Filter1 = bson.M{"stake_deregistration.credential.addrkey_hash": owner}
		oo, nErr := lib.FindAll(findOptions, Skre, Filter1, req1)
		if nErr > 0 {
			log.Warn().Msgf("Find stake registration returned %d errors", nErr)
		} else {
			b, err := json.Marshal(&oo)
			if err == nil {
				utils.PrintJson(string(b))
			}
		}
	}
	return
}
