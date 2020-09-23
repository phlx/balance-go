package move

import (
	"time"
)

type Response struct {
	TransactionFromId   int64     `json:"transaction_from_id"`
	TransactionFromTime time.Time `json:"transaction_from_time"`
	TransactionToId     int64     `json:"transaction_to_id"`
	TransactionToTime   time.Time `json:"transaction_to_time"`
}
