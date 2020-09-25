package metrics

import (
	"time"
)

type Client interface {
	Count(bucket string, n interface{})
	Increment(bucket string)
	Gauge(bucket string, value interface{})
	Timing(bucket string, value interface{})
	Histogram(bucket string, value interface{})
	NewTiming() Timing
}

type Timing interface {
	Send(bucket string)
	Duration() time.Duration
}
