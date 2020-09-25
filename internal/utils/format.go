package utils

import (
	"time"

	"balance/internal/times"
)

func Format(t time.Time) string {
	return t.In(Location()).Format(times.TimeFormatWithMilliseconds)
}
