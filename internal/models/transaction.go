package models

import (
	"time"
)

type Transaction struct {
	tableName   struct{} `pg:"transactions,alias:t"` //nolint
	Id          int64
	UserId      int64   `pg:",notnull"`
	Amount      float64 `pg:",notnull"`
	InitiatorId *int64
	Reason      string
	CreatedAt   time.Time
}
