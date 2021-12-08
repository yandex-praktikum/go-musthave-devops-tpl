package metrics

import "time"

type MetricType string

const (
	GaugeType MetricType 	= "gauge"
	CounterType MetricType 	= "counter"
)

type Timeout time.Ticker

var (
	PollTicker    = time.NewTicker(2 * time.Second)
	ReportTicker  = time.NewTicker(10 * time.Second)
)

