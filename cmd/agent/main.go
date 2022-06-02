package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Metrics struct {
	Alloc,
	TotalAlloc,
	LiveObjects,
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
	OtherSys,
	StackInuse,
	StackSys,
	Sys,

	PauseTotalNs uint64
	NumGC        uint32
	NumGoroutine int
	PollCount    int
	RandomValue  int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func (m *Metrics) UpdateMetrics(duration int) {
	//var m Metrics
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	var PollCount = 0
	m.PollCount = PollCount
	rand.Seed(time.Now().Unix())
	m.RandomValue = rand.Intn(100) + 1
	for {
		<-time.After(interval)
		PollCount++
		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		m.PollCount = PollCount
		rand.Seed(time.Now().Unix())
		m.RandomValue = rand.Intn(10000) + 1

		// Just encode to json and print

		// b, _ := json.Marshal(m)
		// fmt.Println(string(b))

	}
}

func (m *Metrics) PostMetrics(serverAddr string, duration int) {

	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)
		b, _ := json.Marshal(m)
		//fmt.Println(string(b))
		var inInterface map[string]float64
		json.Unmarshal(b, &inInterface)

		for field, val := range inInterface {
			var uri string
			if field != "PollCount" {
				uri = "/update/gauge/" + field + "/" + strconv.FormatFloat(val, 'f', -1, 64)
			} else {
				uri = "/update/counter/" + field + "/" + strconv.FormatFloat(val, 'f', -1, 64)
			}
			//request, err := http.Post(serverAddr+uri, "text/plain", bytes.NewReader(Float64bytes(val)))
			request, err := http.Post(serverAddr+uri, "text/plain", bytes.NewReader([]byte(strconv.FormatFloat(val, 'f', -1, 64))))
			if err != nil {
				log.Fatal(err)
			}
			request.Body.Close()
			//request.Header.Add("Content-Type", "text/plain")
			//fmt.Println(string(request.ContentLength))
			fmt.Println(serverAddr + uri)
		}

	}
}

func main() {
	var metric1 Metrics
	var pollInterval = 2
	var reportInterval = 10
	var url = "http://127.0.0.1"
	var port = "8080"
	var serverAddr = url + ":" + port

	//ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	cmd := exec.Command("ifconfig")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	check(err)
	if err != nil {
		log.Fatal(err)
	}
	go metric1.UpdateMetrics(pollInterval)
	// for ; true; <-ticker.C {
	// 	metric1.PostMetrics(serverAddr, reportInterval)
	// }
	go metric1.PostMetrics(serverAddr, reportInterval)
	for {
		time.Sleep(time.Second)
	}

}
