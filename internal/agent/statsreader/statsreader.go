package statsreader

import (
	"math/rand"
	"runtime"
)

type gauge float64
type counter int64

type MemoryStatsDump struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge

	HeapAlloc    gauge
	HeapIdle     gauge
	HeapInuse    gauge
	HeapObjects  gauge
	HeapReleased gauge

	HeapSys     gauge
	LastGC      gauge
	Lookups     gauge
	MCacheInuse gauge
	MCacheSys   gauge

	MSpanInuse  gauge
	MSpanSys    gauge
	Mallocs     gauge
	NextGC      gauge
	NumForcedGC gauge

	NumGC        gauge
	OtherSys     gauge
	PauseTotalNs gauge
	StackInuse   gauge
	StackSys     gauge

	Sys         gauge
	TotalAlloc  gauge
	PollCount   counter
	RandomValue gauge
}

func (memoryStats *MemoryStatsDump) Refresh() {
	var MemStatistics runtime.MemStats
	runtime.ReadMemStats(&MemStatistics)

	memoryStats.BuckHashSys = gauge(MemStatistics.BuckHashSys)
	memoryStats.Frees = gauge(MemStatistics.Frees)
	memoryStats.GCCPUFraction = gauge(MemStatistics.GCCPUFraction)
	memoryStats.GCSys = gauge(MemStatistics.GCSys)
	memoryStats.HeapAlloc = gauge(MemStatistics.HeapAlloc)

	memoryStats.HeapIdle = gauge(MemStatistics.HeapIdle)
	memoryStats.HeapInuse = gauge(MemStatistics.HeapInuse)
	memoryStats.HeapObjects = gauge(MemStatistics.HeapObjects)
	memoryStats.HeapReleased = gauge(MemStatistics.HeapReleased)
	memoryStats.HeapSys = gauge(MemStatistics.HeapSys)

	memoryStats.LastGC = gauge(MemStatistics.LastGC)
	memoryStats.Lookups = gauge(MemStatistics.Lookups)
	memoryStats.MCacheInuse = gauge(MemStatistics.MCacheInuse)
	memoryStats.MCacheSys = gauge(MemStatistics.MCacheSys)
	memoryStats.MSpanInuse = gauge(MemStatistics.MSpanInuse)

	memoryStats.MSpanSys = gauge(MemStatistics.MSpanSys)
	memoryStats.Mallocs = gauge(MemStatistics.Mallocs)
	memoryStats.NextGC = gauge(MemStatistics.NextGC)
	memoryStats.NumForcedGC = gauge(MemStatistics.NumForcedGC)
	memoryStats.NumGC = gauge(MemStatistics.NumGC)

	memoryStats.OtherSys = gauge(MemStatistics.OtherSys)
	memoryStats.PauseTotalNs = gauge(MemStatistics.PauseTotalNs)
	memoryStats.StackInuse = gauge(MemStatistics.StackInuse)
	memoryStats.StackSys = gauge(MemStatistics.StackSys)

	memoryStats.Alloc = gauge(MemStatistics.Alloc)
	memoryStats.Sys = gauge(MemStatistics.Sys)
	memoryStats.TotalAlloc = gauge(MemStatistics.TotalAlloc)
	memoryStats.PollCount = counter(memoryStats.PollCount + 1)
	memoryStats.RandomValue = gauge(rand.Float64())
}
