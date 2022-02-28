package repository

import (
	"errors"

	"github.com/itd27m01/go-metrics-service/internal/pkg/metrics"
)

var ErrMetricTypeMismatch = errors.New("possible metric type mismatch")

type Store interface {
	UpdateCounterMetric(name string, value metrics.Counter) error
	ResetCounterMetric(name string) error
	UpdateGaugeMetric(name string, value metrics.Gauge) error

	GetMetric(name string) (*metrics.Metric, bool)
	GetMetrics() map[string]*metrics.Metric

	SaveMetrics() error
	LoadMetrics() error
	Close() error
}
