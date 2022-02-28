package workers

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/itd27m01/go-metrics-service/internal/pkg/metrics"
	"github.com/itd27m01/go-metrics-service/internal/repository"
)

const (
	counterIncrement = 1
)

type PollerConfig struct {
	PollInterval time.Duration `env:"POLL_INTERVAL"`
}

type PollerWorker struct {
	Cfg PollerConfig
}

func (pw *PollerWorker) Run(ctx context.Context, mtr repository.Store) {
	pollTicker := time.NewTicker(pw.Cfg.PollInterval)
	defer pollTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			UpdateMemStatsMetrics(mtr)
		}
	}
}

func UpdateMemStatsMetrics(mtr repository.Store) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	mtr.UpdateCounterMetric("PollCount", counterIncrement)

	mtr.UpdateGaugeMetric("Alloc", metrics.Gauge(memStats.Alloc))
	mtr.UpdateGaugeMetric("BuckHashSys", metrics.Gauge(memStats.BuckHashSys))

	mtr.UpdateGaugeMetric("BuckHashSys", metrics.Gauge(memStats.BuckHashSys))
	mtr.UpdateGaugeMetric("Frees", metrics.Gauge(memStats.Frees))
	mtr.UpdateGaugeMetric("GCCPUFraction", metrics.Gauge(memStats.GCCPUFraction))
	mtr.UpdateGaugeMetric("GCSys", metrics.Gauge(memStats.GCSys))
	mtr.UpdateGaugeMetric("HeapAlloc", metrics.Gauge(memStats.HeapAlloc))
	mtr.UpdateGaugeMetric("HeapIdle", metrics.Gauge(memStats.HeapIdle))
	mtr.UpdateGaugeMetric("HeapInuse", metrics.Gauge(memStats.HeapInuse))
	mtr.UpdateGaugeMetric("HeapObjects", metrics.Gauge(memStats.HeapObjects))
	mtr.UpdateGaugeMetric("HeapReleased", metrics.Gauge(memStats.HeapReleased))
	mtr.UpdateGaugeMetric("HeapSys", metrics.Gauge(memStats.HeapSys))
	mtr.UpdateGaugeMetric("LastGC", metrics.Gauge(memStats.LastGC))
	mtr.UpdateGaugeMetric("Lookups", metrics.Gauge(memStats.Lookups))
	mtr.UpdateGaugeMetric("MCacheInuse", metrics.Gauge(memStats.MCacheInuse))
	mtr.UpdateGaugeMetric("MCacheSys", metrics.Gauge(memStats.MCacheSys))
	mtr.UpdateGaugeMetric("MSpanInuse", metrics.Gauge(memStats.MSpanInuse))
	mtr.UpdateGaugeMetric("MSpanSys", metrics.Gauge(memStats.MSpanSys))
	mtr.UpdateGaugeMetric("Mallocs", metrics.Gauge(memStats.Mallocs))
	mtr.UpdateGaugeMetric("NextGC", metrics.Gauge(memStats.NextGC))
	mtr.UpdateGaugeMetric("NumForcedGC", metrics.Gauge(memStats.NumForcedGC))
	mtr.UpdateGaugeMetric("NumGC", metrics.Gauge(memStats.NumGC))
	mtr.UpdateGaugeMetric("OtherSys", metrics.Gauge(memStats.OtherSys))
	mtr.UpdateGaugeMetric("PauseTotalNs", metrics.Gauge(memStats.PauseTotalNs))
	mtr.UpdateGaugeMetric("StackInuse", metrics.Gauge(memStats.StackInuse))
	mtr.UpdateGaugeMetric("StackSys", metrics.Gauge(memStats.StackSys))
	mtr.UpdateGaugeMetric("Sys", metrics.Gauge(memStats.Sys))
	mtr.UpdateGaugeMetric("TotalAlloc", metrics.Gauge(memStats.TotalAlloc))

	mtr.UpdateGaugeMetric("RandomValue", metrics.Gauge(rand.Int63()))
}
