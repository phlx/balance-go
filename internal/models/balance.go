package models

import (
	"time"
)

type Balance struct {
	tableName struct{} `pg:"balances,alias:b"` //nolint
	Id        int64
	UserId    int64   `pg:",notnull,unique"`
	Balance   float64 `pg:",notnull,use_zero"`
	UpdatedAt time.Time
}
