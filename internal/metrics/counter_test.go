package metrics

import (
	"testing"
)

func TestCounter(t *testing.T) {
	c := Counter{}
	c.Increment(1)
	if c.Get() != 1 {
		t.Error("Result not equal 1")
	}
}
