package models

import (
	"time"
)

type Response struct {
	tableName      struct{} `pg:"responses,alias:r"` //nolint
	IdempotencyKey string   `pg:"type:varchar(255),pk,notnull,unique"`
	Status         int      `pg:",notnull"`
	Headers        string   `pg:",notnull"`
	Response       string   `pg:",notnull"`
	CreatedAt      time.Time
}
