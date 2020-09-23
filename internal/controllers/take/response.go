package take

import (
	"time"
)

type Response struct {
	Transaction int64     `json:"transaction"`
	Time        time.Time `json:"time"`
}
