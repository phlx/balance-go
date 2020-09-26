package health

type Response struct {
	Postgres bool     `json:"postgres"`
	Redis    bool     `json:"redis"`
	Errors   []string `json:"errors"`
	Time     string   `json:"time"`
	Latency  int64    `json:"latency_ms"`
}
