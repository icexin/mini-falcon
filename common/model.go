package common

type Metric struct {
	Metric    string   `json:"metric"`
	Endpoint  string   `json:"endpoint"`
	Tag       []string `json:"tag"`
	Value     float64  `json:"value"`
	Timestamp int64    `json:"timestamp"`
}

type Alarm struct {
	Mail     []string
	Expr     string
	Value    float64
	Endpoint string
	Tag      []string
	Recovery bool
}
