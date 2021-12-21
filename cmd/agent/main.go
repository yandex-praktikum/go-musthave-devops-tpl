package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/efrikin/go-musthave-devops-tpl/internal/metrics"
	"github.com/efrikin/go-musthave-devops-tpl/internal/models"
)

func main() {
	var cfg models.Config
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "Send metrics to address:port")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "Report of interval")
	flag.DurationVar(&cfg.PoolInterval, "p", 2*time.Second, "Pool of interval")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	var (
		pollTicker    = time.NewTicker(cfg.PoolInterval)
		reportTicker  = time.NewTicker(cfg.ReportInterval)
		Alloc         = metrics.NewGauge("Alloc")
		TotalAlloc    = metrics.NewGauge("TotalAlloc")
		BuckHashSys   = metrics.NewGauge("BuckHashSys")
		Frees         = metrics.NewGauge("Frees")
		GCCPUFraction = metrics.NewGauge("GCCPUFraction")
		GCSys         = metrics.NewGauge("GCSys")
		HeapAlloc     = metrics.NewGauge("HeapAlloc")
		HeapIdle      = metrics.NewGauge("HeapIdle")
		HeapInuse     = metrics.NewGauge("HeapInuse")
		HeapObjects   = metrics.NewGauge("HeapObjects")
		HeapReleased  = metrics.NewGauge("HeapReleased")
		HeapSys       = metrics.NewGauge("HeapSys")
		LastGC        = metrics.NewGauge("LastGC")
		Lookups       = metrics.NewGauge("Lookups")
		MCacheInuse   = metrics.NewGauge("MCacheInuse")
		MCacheSys     = metrics.NewGauge("MCacheSys")
		MSpanInuse    = metrics.NewGauge("MSpanInuse")
		MSpanSys      = metrics.NewGauge("MSpanSys")
		Mallocs       = metrics.NewGauge("Mallocs")
		NextGC        = metrics.NewGauge("NextGC")
		NumForcedGC   = metrics.NewGauge("NumForcedGC")
		NumGC         = metrics.NewGauge("NumGC")
		OtherSys      = metrics.NewGauge("OtherSys")
		PauseTotalNs  = metrics.NewGauge("PauseTotalNs")
		StackInuse    = metrics.NewGauge("StackInuse")
		StackSys      = metrics.NewGauge("StackSys")
		Sys           = metrics.NewGauge("Sys")
		RandomValue   = metrics.NewGauge("RandomValue")
		PollCount     = metrics.NewCounter("PollCount")
	)

	go func() {
		rand.Seed(time.Now().UnixNano())
		m := runtime.MemStats{}
		for range pollTicker.C {
			runtime.ReadMemStats(&m)
			Alloc.Set(float64(m.Alloc))
			TotalAlloc.Set(float64(m.TotalAlloc))
			BuckHashSys.Set(float64(m.BuckHashSys))
			Frees.Set(float64(m.Frees))
			GCCPUFraction.Set(float64(m.GCCPUFraction))
			GCSys.Set(float64(m.GCSys))
			HeapAlloc.Set(float64(m.HeapAlloc))
			HeapIdle.Set(float64(m.HeapIdle))
			HeapInuse.Set(float64(m.HeapInuse))
			HeapObjects.Set(float64(m.HeapObjects))
			HeapReleased.Set(float64(m.HeapReleased))
			HeapSys.Set(float64(m.HeapSys))
			LastGC.Set(float64(m.LastGC))
			Lookups.Set(float64(m.Lookups))
			MCacheInuse.Set(float64(m.MCacheInuse))
			MCacheSys.Set(float64(m.MCacheSys))
			MSpanInuse.Set(float64(m.MSpanInuse))
			MSpanSys.Set(float64(m.MSpanSys))
			Mallocs.Set(float64(m.Mallocs))
			NextGC.Set(float64(m.NextGC))
			NumForcedGC.Set(float64(m.NumForcedGC))
			NumGC.Set(float64(m.NumGC))
			OtherSys.Set(float64(m.OtherSys))
			PauseTotalNs.Set(float64(m.PauseTotalNs))
			StackInuse.Set(float64(m.StackInuse))
			StackSys.Set(float64(m.StackSys))
			Sys.Set(float64(m.Sys))
			RandomValue.Set(rand.Float64())
			PollCount.Increment(1)
		}
	}()
	go func() {
		b := bytes.NewBuffer([]byte{})
		var m = []interface{}{
			Alloc,
			TotalAlloc,
			BuckHashSys,
			Frees,
			GCCPUFraction,
			GCSys,
			HeapAlloc,
			HeapIdle,
			HeapInuse,
			HeapObjects,
			HeapReleased,
			HeapSys,
			LastGC,
			Lookups,
			MCacheInuse,
			MCacheSys,
			MSpanInuse,
			MSpanSys,
			Mallocs,
			NextGC,
			NumForcedGC,
			NumGC,
			OtherSys,
			PauseTotalNs,
			StackInuse,
			StackSys,
			Sys,
			RandomValue,
			PollCount,
		}

		for range reportTicker.C {
			for _, v := range m {
				var body []byte
				typedV, ok := v.(*metrics.Gauge)
				if ok {
					tmpV := typedV.Get()
					body, _ = json.Marshal(models.Metrics{
						ID:    typedV.Name(),
						MType: string(typedV.Type()),
						Value: &tmpV,
					})
				} else {
					typedV := v.(*metrics.Counter)
					tmpV := typedV.Get()
					body, _ = json.Marshal(models.Metrics{
						ID:    typedV.Name(),
						MType: string(typedV.Type()),
						Delta: &tmpV,
					})
				}
				b.Reset()
				b.Write(body)
				fmt.Printf("%s\n", b.String())
				r, err := http.Post(fmt.Sprintf("http://%s/update", cfg.Address), "application/json", b)
				if err == nil {
					r.Body.Close()
				}
				if err != nil {
					fmt.Printf("Unable to send metric %s to server: %v\n", typedV.Name(), err)
					continue
				}
			}
		}
	}()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
