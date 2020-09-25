package metrics

import (
	"gopkg.in/alexcesaro/statsd.v2"
)

type StatsDDecorator struct {
	statsd *statsd.Client
}

func New(statsd *statsd.Client) Client {
	return &StatsDDecorator{statsd: statsd}
}

func (c *StatsDDecorator) Count(bucket string, n interface{}) {
	c.statsd.Count(bucket, n)
}

func (c *StatsDDecorator) Increment(bucket string) {
	c.statsd.Increment(bucket)
}

func (c *StatsDDecorator) Gauge(bucket string, value interface{}) {
	c.statsd.Gauge(bucket, value)
}

func (c *StatsDDecorator) Timing(bucket string, value interface{}) {
	c.statsd.Timing(bucket, value)
}

func (c *StatsDDecorator) Histogram(bucket string, value interface{}) {
	c.statsd.Histogram(bucket, value)
}

func (c *StatsDDecorator) NewTiming() Timing {
	return c.statsd.NewTiming()
}
