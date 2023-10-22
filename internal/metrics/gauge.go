package metrics

import (
	"fmt"
	"sync"
)

type Gauge struct {
	Gauge float64
	ID    string
	mu    sync.Mutex
}

func (c *Gauge) Set(v float64) {
	c.mu.Lock()
	c.Gauge = v
	c.mu.Unlock()
}

func (c *Gauge) Get() float64 {
	var tmp float64
	c.mu.Lock()
	tmp = c.Gauge
	c.mu.Unlock()
	return tmp
}

func (c *Gauge) Name() string {
	return c.ID
}

func (c *Gauge) Type() MetricType {
	return GaugeType
}

func NewGauge(n string) *Gauge {
	return &Gauge{ID: n}
}

func (c *Gauge) String() string {
	return fmt.Sprintf("%v", c.Get())
}
