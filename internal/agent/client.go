package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Monitor struct {
	name  string
	stats uint64
	types string
}

type Metrics struct {
	ID    string `json:"id"`              // имя метрики
	MType string `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta uint64 `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value string `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMonitor(duration int64, ch chan int64, counter int64) {
	//var m Monitor
	//counter++
	var rtm runtime.MemStats

	//counter := 0
	var interval = time.Duration(duration) * time.Second
	//for {
	//<-time.After(interval) //place of problem
	time.AfterFunc(interval, func() {
		ch <- 1
		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		//m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		//m.Alloc = rtm.Alloc
		//m.TotalAlloc = rtm.TotalAlloc
		//m.Sys = rtm.Sys
		//m.Mallocs = rtm.Mallocs
		//m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		//m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		//m.PauseTotalNs = rtm.PauseTotalNs
		//m.NumGC = rtm.NumGC

		// Just encode to json and print
		//b, _ := json.Marshal(m)
		//fmt.Println(string(b))

		endpoint := "http://localhost:8080/update/"
		// контейнер данных для запроса
		data := url.Values{}
		// cnt :=counter++
		long := ""
		Alloc := Monitor{"Alloc", rtm.Alloc, "gauge"}
		TotalAlloc := Monitor{"TotalAlloc", rtm.Alloc, "gauge"}
		Sys := Monitor{"Sys", rtm.Sys, "Gauge"}
		Mallocs := Monitor{"Mallocs", rtm.Mallocs, "gauge"}
		Frees := Monitor{"NumGC", rtm.Frees, "gauge"}
		PollCount := Monitor{"PollCount", uint64(counter), "counter"}

		v := []Monitor{Alloc, TotalAlloc, Sys, Mallocs, Frees, PollCount}
		for _, service := range v {
			longurl := service.name + "/" + service.types + "/" + strconv.FormatUint(service.stats, 10)
			endpoints := endpoint + longurl
			long = strings.TrimSuffix(long, "\n")
			//metric := []Metrics{}
			if service.types == "gauge" {

				metric := Metrics{
					ID:    service.name,
					MType: service.types,
					Value: strconv.FormatUint(service.stats, 10),
				}

			} else {
				metric := Metrics{
					ID:    service.name,
					MType: service.types,
					Delta: service.stats,
				}
			}

			json, err := json.Marshal(metric)
			if err != nil {
				panic(err)
			}

			// заполняем контейнер данными
			//data.Set("url", long)
			// конструируем HTTP-клиент
			client := &http.Client{}
			// конструируем запрос
			// запрос методом POST должен, кроме заголовков, содержать тело
			// тело должно быть источником потокового чтения io.Reader
			// в большинстве случаев отлично подходит bytes.Buffer
			request, err := http.NewRequest(http.MethodPost, endpoints, bytes.NewBuffer(json))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
			// отправляем запрос и получаем ответ
			response, err := client.Do(request)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// печатаем код ответа
			fmt.Println("Статус-код ", response.Status)
			defer response.Body.Close()
			// читаем поток из тела ответа
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// и печатаем его
			fmt.Println(string(body))
		}
	})
}
