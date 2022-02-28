package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/itd27m01/go-metrics-service/internal/pkg/metrics"
)

const (
	fileMode = 0640
)

type FileStore struct {
	file         *os.File
	syncChannel  chan struct{}
	metricsCache map[string]*metrics.Metric
	mu           sync.Mutex
}

func NewFileStore(filePath string, syncChannel chan struct{}) (*FileStore, error) {
	var fs FileStore

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, fileMode)
	if err != nil {
		return nil, err
	}

	metricsCache := make(map[string]*metrics.Metric)
	fs = FileStore{
		file:         file,
		syncChannel:  syncChannel,
		metricsCache: metricsCache,
	}

	return &fs, nil
}

func (fs *FileStore) UpdateCounterMetric(metricName string, metricData metrics.Counter) error {
	fs.mu.Lock()
	defer fs.sync()
	defer fs.mu.Unlock()

	currentMetric, ok := fs.metricsCache[metricName]
	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) += metricData
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("%w %s:%s", ErrMetricTypeMismatch, metricName, currentMetric.MType)
	default:
		fs.metricsCache[metricName] = &metrics.Metric{
			ID:    metricName,
			MType: metrics.CounterMetricTypeName,
			Delta: &metricData,
		}
	}

	return nil
}

func (fs *FileStore) ResetCounterMetric(metricName string) error {
	fs.mu.Lock()
	defer fs.sync()
	defer fs.mu.Unlock()

	var zero metrics.Counter
	currentMetric, ok := fs.metricsCache[metricName]
	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) = zero
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("%w %s:%s", ErrMetricTypeMismatch, metricName, currentMetric.MType)
	default:
		fs.metricsCache[metricName] = &metrics.Metric{
			ID:    metricName,
			MType: metrics.CounterMetricTypeName,
			Delta: &zero,
		}
	}

	return nil
}

func (fs *FileStore) UpdateGaugeMetric(metricName string, metricData metrics.Gauge) error {
	fs.mu.Lock()
	defer fs.sync()
	defer fs.mu.Unlock()

	currentMetric, ok := fs.metricsCache[metricName]
	switch {
	case ok && currentMetric.Value != nil:
		*(currentMetric.Value) = metricData
	case ok && currentMetric.Value == nil:
		return fmt.Errorf("%w %s:%s", ErrMetricTypeMismatch, metricName, currentMetric.MType)
	default:
		fs.metricsCache[metricName] = &metrics.Metric{
			ID:    metricName,
			MType: metrics.GaugeMetricTypeName,
			Value: &metricData,
		}
	}

	return nil
}

func (fs *FileStore) GetMetric(metricName string) (*metrics.Metric, bool) {
	metric, ok := fs.metricsCache[metricName]

	return metric, ok
}

func (fs *FileStore) GetMetrics() map[string]*metrics.Metric {
	return fs.metricsCache
}

func (fs *FileStore) sync() {
	fs.syncChannel <- struct{}{}
}

func (fs *FileStore) Close() error {
	if err := fs.SaveMetrics(); err != nil {
		log.Printf("Something went wrong durin metrics preserve %q", err)
	}

	if err := fs.file.Sync(); err != nil {
		log.Printf("Failed to sync metrics: %q", err)
	}

	return fs.file.Close()
}

func (fs *FileStore) LoadMetrics() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	jsonDecoder := json.NewDecoder(fs.file)

	log.Printf("Load metrics from %s", fs.file.Name())

	return jsonDecoder.Decode(&(fs.metricsCache))
}

func (fs *FileStore) SaveMetrics() (err error) {
	log.Printf("Dump metrics to %s", fs.file.Name())

	fs.mu.Lock()
	defer fs.mu.Unlock()

	const (
		offset     = 0
		whence     = 0
		truncateTo = 0
	)
	_, err = fs.file.Seek(offset, whence)
	if err != nil {
		return err
	}

	if err := fs.file.Truncate(truncateTo); err != nil {
		return err
	}

	encoder := json.NewEncoder(fs.file)

	return encoder.Encode(&fs.metricsCache)
}
