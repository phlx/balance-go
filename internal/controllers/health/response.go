package health

import (
	"time"
)

type Response struct {
	Postgres bool      `json:"postgres"`
	Redis    bool      `json:"redis"`
	Errors   []string  `json:"errors"`
	Time     time.Time `json:"time"`
	Latency  int64     `json:"latency_ms"`
}
