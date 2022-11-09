package types

import (
	"encoding/json"
	"github.com/paulmatencio/oura-go/utils"
	"github.com/rs/zerolog/log"
)

type Events struct {
	Ev  map[string][]string
	Kp  map[string]bool
	Kp1 Set[string]
}

func (m *Events) Init() {
	m.Ev = make(map[string][]string)
	m.Ev["tx"] = []string{}
	m.Ev["utxo"] = []string{}
	m.Ev["asst"] = []string{}
	m.Ev["mint"] = []string{}
	m.Ev["stxi"] = []string{}
	m.Ev["coll"] = []string{}
	m.Ev["meta"] = []string{}
	m.Ev["dtum"] = []string{}
	m.Ev["blck"] = []string{}
	m.Ev["rdmr"] = []string{}
	m.Ev["witn"] = []string{}
	m.Ev["witp"] = []string{}
	m.Ev["pool"] = []string{}
	m.Ev["dele"] = []string{}
	m.Ev["reti"] = []string{}
	m.Ev["skde"] = []string{}
	m.Ev["skre"] = []string{}
	m.Ev["cip25"] = []string{}
	m.Ev["cip15"] = []string{}
	m.Ev["scpt"] = []string{}
	m.Ev["other"] = []string{}
}

func (m *Events) Set(typ string, v string) {
	switch typ {
	case "tx":
		m.Ev["tx"] = append(m.Ev["tx"], v)
	case "utxo":
		m.Ev["utxo"] = append(m.Ev["utxo"], v)
	case "stxi":
		m.Ev["stxi"] = append(m.Ev["stxi"], v)
	case "asst":
		m.Ev["asst"] = append(m.Ev["asst"], v)
	case "meta":
		m.Ev["meta"] = append(m.Ev["meta"], v)
	case "coll":
		m.Ev["coll"] = append(m.Ev["coll"], v)
	case "mint":
		m.Ev["mint"] = append(m.Ev["mint"], v)
	case "dtum":
		m.Ev["dtum"] = append(m.Ev["dtum"], v)
	case "blck":
		m.Ev["blck"] = append(m.Ev["blck"], v)
	case "rdmr":
		m.Ev["rdmr"] = append(m.Ev["rdmr"], v)
	case "witn":
		m.Ev["witn"] = append(m.Ev["witn"], v)
	case "witp":
		m.Ev["witp"] = append(m.Ev["witp"], v)
	case "pool":
		m.Ev["pool"] = append(m.Ev["pool"], v)
	case "reti":
		m.Ev["reti"] = append(m.Ev["reti"], v)
	case "dele":
		m.Ev["dele"] = append(m.Ev["dele"], v)
	case "skde":
		m.Ev["skde"] = append(m.Ev["skde"], v)
	case "skre":
		m.Ev["skre"] = append(m.Ev["skre"], v)
	case "cip25":
		m.Ev["cip25"] = append(m.Ev["cip25"], v)
	case "cip15":
		m.Ev["cip15"] = append(m.Ev["cip15"], v)
	case "scpt":
		m.Ev["scpt"] = append(m.Ev["scpt"], v)
	default:
		log.Trace().Msgf("Other type: %s", typ)
		m.Ev["other"] = append(m.Ev["other"], v)
	}
}

func (m *Events) SplitTypes(t []string) {

	var unknown Unknown
	for _, v := range t {
		err := json.Unmarshal([]byte(v), &unknown)
		if err == nil {
			if unknown.Fingerprint != "" {
				typ, _ := utils.GetType(unknown.Fingerprint)
				m.Set(typ, v)
			}
		}
	}
}

func (m *Events) Keep(list []string) {
	m.Kp = make(map[string]bool)
	for _, v := range list {
		m.Kp[v] = true
	}
}

func (m *Events) Has(key string) bool {
	if _, ok := m.Kp[key]; ok {
		return true
	}
	return false
}

func (m *Events) Selected(events []string) {
	m.Kp1 = NewSet("")
	m.Kp1.Add(events)
}
