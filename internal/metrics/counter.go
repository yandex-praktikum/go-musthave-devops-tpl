package metrics

import (
	"fmt"
	"sync"
)

type Counter struct {
	count int64
	name  string
	mu    sync.Mutex
}

func (c *Counter) Increment(i int64) {
	c.mu.Lock()
	// TODO: There is problem. Handle overflow
	c.count += i
	c.mu.Unlock()
}

func (c *Counter) Get() int64 {
	var tmp int64
	c.mu.Lock()
	tmp = c.count
	c.mu.Unlock()
	return tmp
}

func (c *Counter) Name() string {
	return c.name
}

func (c *Counter) Type() MetricType {
	return CounterType
}

func NewCounter(n string) *Counter {
	return &Counter{name: n}
}

func (c *Counter) String() string {
	return fmt.Sprintf("%v", c.Get())
}
