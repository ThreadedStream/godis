package main

type KeyValue struct {
	Key    string      `json:"key"`
	Value  interface{} `json:"value"`
	Field  string      `json:"field"`
	HasTtl bool        `json:"has_ttl"`
	Exp    int64       `json:"exp"`
}

type ResponseEntities struct {
	Status  string   `json:"status"`
	Keys    []string `json:"keys"`
	Value   string   `json:"value"`
	Values  []string `json:"values"`
	Created int      `json:"created"`
}

type KeysRequestModel struct {
	Pattern string   `json:"pattern"`
	Keys    []string `json:"keys"`
}

type DelRequestModel struct {
	Key string `json:"key"`
}
