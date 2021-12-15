package metrics

import (
	"fmt"
	"sync"
)

type Gauge struct {
	gauge float64
	name  string
	mu    sync.Mutex
}

func (c *Gauge) Set(v float64) {
	c.mu.Lock()
	c.gauge = v
	c.mu.Unlock()
}

func (c *Gauge) Get() float64 {
	var tmp float64
	c.mu.Lock()
	tmp = c.gauge
	c.mu.Unlock()
	return tmp
}

func (c *Gauge) Name() string {
	return c.name
}

func (c *Gauge) Type() MetricType {
	return GaugeType
}

func NewGauge(n string) *Gauge {
	return &Gauge{name: n}
}

func (c *Gauge) String() string {
	return fmt.Sprintf("%v", c.Get())
}
