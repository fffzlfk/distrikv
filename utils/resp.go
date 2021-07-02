package utils

type Resp struct {
	Shard    int    `json:"shard"`
	CurShard int    `json:"current-shard"`
	Addr     string `json:"addr"`
	Value    string `json:"value"`
	Err      error  `json:"error"`
}
