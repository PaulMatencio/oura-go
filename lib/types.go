package lib

/*

	Working type
	intermediate storage
*/
type Assets struct {
	Policy string
	// Fingerprint  string
	Description  interface{}
	PolicyLink   string
	PolicyScript interface{}
	Publisher    interface{}
	CopyRight    interface{}
	Image        interface{}
	Name         interface{}
	Artist       string
	MediaType    string
	Version      string
	RawJson      interface{}
	Asset        []Asset
}

type Asset struct {
	Source       string
	Policy       string
	Asset        string
	Fingerprint  string
	RawJSON      interface{}
	Name         interface{}
	Artist       string
	Image        interface{}
	Description  interface{}
	PolicyLink   string
	PolicyScript interface{}
	Publisher    interface{}
	MediaType    string
	CopyRight    interface{}
	Version      string
	Website      interface{}
	Project      interface{}
}
