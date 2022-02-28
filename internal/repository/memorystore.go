package repository

import (
	"fmt"
	"sync"

	"github.com/itd27m01/go-metrics-service/internal/pkg/metrics"
)

type InMemoryStore struct {
	metricsCache map[string]*metrics.Metric
	mu           sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	var m InMemoryStore

	m.metricsCache = make(map[string]*metrics.Metric)

	return &m
}

func (m *InMemoryStore) UpdateCounterMetric(metricName string, metricData metrics.Counter) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	currentMetric, ok := m.metricsCache[metricName]
	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) += metricData
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("%w %s:%s", ErrMetricTypeMismatch, metricName, currentMetric.MType)
	default:
		m.metricsCache[metricName] = &metrics.Metric{
			ID:    metricName,
			MType: metrics.CounterMetricTypeName,
			Delta: &metricData,
		}
	}

	return nil
}

func (m *InMemoryStore) ResetCounterMetric(metricName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zero metrics.Counter
	currentMetric, ok := m.metricsCache[metricName]
	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) = zero
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("%w %s:%s", ErrMetricTypeMismatch, metricName, currentMetric.MType)
	default:
		m.metricsCache[metricName] = &metrics.Metric{
			ID:    metricName,
			MType: metrics.CounterMetricTypeName,
			Delta: &zero,
		}
	}

	return nil
}

func (m *InMemoryStore) UpdateGaugeMetric(metricName string, metricData metrics.Gauge) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	currentMetric, ok := m.metricsCache[metricName]
	switch {
	case ok && currentMetric.Value != nil:
		*(currentMetric.Value) = metricData
	case ok && currentMetric.Value == nil:
		return fmt.Errorf("%w %s:%s", ErrMetricTypeMismatch, metricName, currentMetric.MType)
	default:
		m.metricsCache[metricName] = &metrics.Metric{
			ID:    metricName,
			MType: metrics.GaugeMetricTypeName,
			Value: &metricData,
		}
	}

	return nil
}

func (m *InMemoryStore) GetMetric(metricName string) (*metrics.Metric, bool) {
	metric, ok := m.metricsCache[metricName]

	return metric, ok
}

func (m *InMemoryStore) GetMetrics() map[string]*metrics.Metric {
	return m.metricsCache
}

func (m *InMemoryStore) SaveMetrics() error { return nil }
func (m *InMemoryStore) LoadMetrics() error { return nil }
func (m *InMemoryStore) Close() error       { return nil }
