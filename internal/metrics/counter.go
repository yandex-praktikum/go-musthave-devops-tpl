package metrics

import (
	"fmt"
	"sync"
)

type Counter struct {
	Count int64
	ID    string
	mu    sync.Mutex
}

func (c *Counter) Increment(i int64) {
	c.mu.Lock()
	// TODO: There is problem. Handle overflow
	c.Count += i
	c.mu.Unlock()
}

func (c *Counter) Get() int64 {
	var tmp int64
	c.mu.Lock()
	tmp = c.Count
	c.mu.Unlock()
	return tmp
}

func (c *Counter) Name() string {
	return c.ID
}

func (c *Counter) Type() MetricType {
	return CounterType
}

func NewCounter(n string) *Counter {
	return &Counter{ID: n}
}

func (c *Counter) String() string {
	return fmt.Sprintf("%v", c.Get())
}
