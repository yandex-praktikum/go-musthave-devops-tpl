package metrics

import (
	"testing"
)

func TestGauge(t *testing.T) {
	g := Gauge{}
	g.Set(1.0)
	if g.Get() != 1.0 {
		t.Error("Result not equal 1.0")
	}
}
