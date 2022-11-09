package types

type Scpt struct {
	Context      Context      `json:"context"`
	Fingerprint  string       `json:"fingerprint"`
	NativeScript NativeScript `json:"native_script"`
}

type Scpts struct {
	KeyHash string `json:"keyHash,omitempty"`
	Type    string `json:"type"`
	Slot    int    `json:"slot,omitempty"`
}

type Script struct {
	Scripts []Scpts `json:"scripts"`
	Type    string  `json:"type"`
}

type NativeScript struct {
	PolicyID string `json:"policy_id"`
	Script   Script `json:"script"`
}
