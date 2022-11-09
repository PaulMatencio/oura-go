package types

import (
	"encoding/json"
	"fmt"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
)

type Counter struct {
	Txs    int `json:"number_tx"`
	Utxos  int `json:"number_utxo"`
	Stxis  int `json:"number_stxi"`
	Assets int `json:"number_asset"`
	Metas  int `json:"meta_number"`
	Colls  int `json:"collateral_number"`
	Mints  int `json:"mint_number"`
	Datums int `json:"datum_number"`
	Witns  int `json:"native_witness"`
	Witps  int `json:"plutus_witness"`
	Blocks int `json:"cardano_block"`
	Rdmrs  int `json:"plutus_redeemer"`
	Pools  int `json:"pool_registration"`
	Deles  int `json:"pool_delegation"`
	Retis  int `json:"pool_retirement"`
	Skdes  int `json:"stake_deregistration"`
	Skres  int `json:"stake_registration"`
	Cip25s int `json:"cip_25"`
	Cip15s int `json:"cip_15"`
	Scpts  int `json:"script"`
	Others int `json:"other_number"`
}

func (c *Counter) Increment(typ string) {
	switch typ {
	case "tx":
		c.Txs++
	case "utxo":
		c.Utxos++
	case "stxi":
		c.Stxis++
	case "asst":
		c.Assets++
	case "meta":
		c.Metas++
	case "coll":
		c.Colls++
	case "mint":
		c.Mints++
	case "dtum":
		c.Datums++
	case "witn":
		c.Witns++
	case "witp":
		c.Witps++
	case "rdmr":
		c.Rdmrs++
	case "blck":
		c.Blocks++
	case "pool":
		c.Pools++
	case "dele":
		c.Deles++
	case "reti":
		c.Retis++
	case "skde":
		c.Skdes++
	case "skre":
		c.Skres++
	case "cip25":
		c.Cip25s++
	case "cip15":
		c.Cip15s++
	case "scpt":
		c.Scpts++
	default:
		log.Trace().Msgf("Other type: %s", typ)
		c.Others++
	}
}

func (c *Counter) Add(typ string, n int) {
	switch typ {
	case "tx":
		c.Txs += n
	case "utxo":
		c.Utxos += n
	case "stxi":
		c.Stxis += n
	case "asst":
		c.Assets += n
	case "meta":
		c.Metas += n
	case "coll":
		c.Colls += n
	case "mint":
		c.Mints += n
	case "dtum":
		c.Datums += n
	case "witn":
		c.Witns += n
	case "witp":
		c.Witps += n
	case "rdmr":
		c.Rdmrs += n
	case "blck":
		c.Blocks += n
	case "pool":
		c.Pools += n
	case "dele":
		c.Deles += n
	case "reti":
		c.Retis += n
	case "skde":
		c.Skdes += n
	case "skre":
		c.Skres += n
	case "cip25":
		c.Cip25s += n
	case "cip15":
		c.Cip15s += n
	case "scpt":
		c.Scpts += n
	default:
		log.Trace().Msgf("Other type: %s", typ)
		c.Others += n
	}
}

func (c *Counter) Print() {
	b, err := json.Marshal(c)
	if err == nil {
		j, _ := utils.PrettyJson(string(b))
		fmt.Printf("\nCounters:\n %s\n", j)
	}
	c.PrintTotal()
}

func (c *Counter) GetTotal() int {
	return c.Txs +
		c.Utxos +
		c.Assets +
		c.Stxis +
		c.Mints +
		c.Colls +
		c.Datums +
		c.Metas +
		c.Witns +
		c.Witps +
		c.Rdmrs +
		c.Cip25s +
		c.Retis +
		c.Pools +
		c.Skres +
		c.Skdes +
		c.Deles +
		c.Blocks +
		c.Cip15s +
		c.Scpts +
		c.Others
}

func (c *Counter) PrintTotal() {
	fmt.Printf("Total %d\n", c.GetTotal())
}
