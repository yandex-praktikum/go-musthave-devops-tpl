package metrics

import (
	"bytes"
	"encoding/json"
)

const (
	GaugeMetricTypeName   = "gauge"
	CounterMetricTypeName = "counter"
)

type Gauge float64
type Counter int64

type Metric struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр, принимающий значение gauge или counter
	Delta *Counter `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

func (m *Metric) EncodeMetric() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	jsonEncoder := json.NewEncoder(&buf)

	if err := jsonEncoder.Encode(m); err != nil {
		return nil, err
	}

	return &buf, nil
}
