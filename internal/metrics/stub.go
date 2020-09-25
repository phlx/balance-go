package metrics

import (
	"time"
)

type TimingStub struct{}

func (t TimingStub) Send(_ string) {}
func (t TimingStub) Duration() time.Duration {
	return 0 * time.Second
}

type Stub struct{}

func (s Stub) Count(_ string, _ interface{})     {}
func (s Stub) Increment(_ string)                {}
func (s Stub) Gauge(_ string, _ interface{})     {}
func (s Stub) Timing(_ string, _ interface{})    {}
func (s Stub) Histogram(_ string, _ interface{}) {}
func (s Stub) NewTiming() Timing {
	return TimingStub{}
}

func NewStub() Client {
	return &Stub{}
}
