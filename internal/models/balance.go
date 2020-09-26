package models

import (
	"time"
)

type Balance struct {
	tableName struct{} `pg:"balances,alias:b"` //nolint
	UserId    int64    `pg:",pk,notnull,unique"`
	Balance   float64  `pg:"type:numeric(20,2),notnull,use_zero"`
	UpdatedAt time.Time
}
