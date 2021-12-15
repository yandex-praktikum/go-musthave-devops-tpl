package metrics

type MetricType string

const (
	GaugeType    MetricType = "gauge"
	CounterType  MetricType = "counter"
	PullTicker              = 2
	ReportTicker            = 10
)
