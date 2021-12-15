package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
        "html/template"
	"encoding/json"
	"io"
	"github.com/caarlos0/env/v6"
	"github.com/efrikin/go-musthave-devops-tpl/internal/metrics"
	"github.com/efrikin/go-musthave-devops-tpl/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	gaugeStorage = map[string]*metrics.Gauge{}
	counterStorage = map[string]*metrics.Counter{}
	mu      = sync.Mutex{}
)

func httpPrintJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "BadRequests", http.StatusBadRequest)
	}
	v := models.Metrics{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		http.Error(w, "BadRequests", http.StatusBadRequest)
	}

	metricType := metrics.MetricType(v.MType)
	metricName := v.ID

	if metricType == metrics.GaugeType {
		metricValue := v.Value
		mu.Lock()
		defer mu.Unlock()
		w.WriteHeader(http.StatusOK)
		metric := &metrics.Gauge{}
		metric.Set(*metricValue)
		gaugeStorage[metricName] = metric
		return
	}

	if metricType == metrics.CounterType {
		metricValue := v.Delta
		mu.Lock()
		defer mu.Unlock()
		var metric, ok = counterStorage[metricName]
		if !ok {
			metric = &metrics.Counter{}
			counterStorage[metricName] =  metric
		}
		w.WriteHeader(http.StatusOK)
		metric.Increment(*metricValue)
		// I want fix this test =\
		fmt.Fprintf(w, "{}")
		return
	}
	http.Error(w, "Error", http.StatusNotImplemented)
}

func httpPrint(w http.ResponseWriter, r *http.Request) {

	metricType := metrics.MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType == metrics.GaugeType {
		metricValueTyped, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		w.WriteHeader(http.StatusOK)
		metric := &metrics.Gauge{}
		metric.Set(metricValueTyped)
		gaugeStorage[metricName] = metric
		return
	}

	if metricType == metrics.CounterType {
		mu.Lock()
		defer mu.Unlock()
		var metric, ok = counterStorage[metricName]
		if !ok {
			metric = &metrics.Counter{}
			counterStorage[metricName] =  metric
		}
		metricValueTyped, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Bad convenrt int to string", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		metric.Increment(metricValueTyped)
		return
	}
	http.Error(w, "Error", http.StatusNotImplemented)
}

func httpPrintMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "BadRequests", http.StatusBadRequest)
	}
	v := models.Metrics{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		http.Error(w, "BadRequests", http.StatusBadRequest)
	}

	metricType := metrics.MetricType(v.MType)
	metricName := v.ID

	if metricType == metrics.GaugeType {
		mu.Lock()
		defer mu.Unlock()
		metric, ok := gaugeStorage[metricName]
		if !ok {
			http.Error(w,"NotFound", http.StatusNotFound)
			return
		}
		tmpV := metric.Get()
		v.Value = &tmpV
		body, _ = json.Marshal(v)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", body)
		return
	}

	if metricType == metrics.CounterType {
		mu.Lock()
		defer mu.Unlock()
		metric, ok := counterStorage[metricName]
		if !ok {
			http.Error(w,"NotFound", http.StatusNotFound)
			return
		}
		tmpV := metric.Get()
		v.Delta = &tmpV
		body, _ = json.Marshal(v)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", body)
		return
	}
	// return
}


func httpPrintGaugeMetrics(w http.ResponseWriter, r *http.Request) {

	metricName := chi.URLParam(r, "metricName")
	val, ok := gaugeStorage[metricName]

	if !ok {
		http.Error(w,"Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v\n", val)
}

func httpPrintCounterMetrics(w http.ResponseWriter, r *http.Request) {

	metricName := chi.URLParam(r, "metricName")
	val, ok := counterStorage[metricName]

	if !ok {
		http.Error(w,"Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v\n", val)
}

func httpPrintMetricsHTML(w http.ResponseWriter, r *http.Request) {
	gaugeTempl, _ := template.New("printMetricsHTML").Parse(`
	<html>
	  <head>
	    <title>GAUGE</title>
	   <meta http-equiv="refresh" content="10" />
	  </head>
	  <h1><center>TYPE OF GAUGE METRICS</center></h1>
	  {{ range $key, $value := . }}
	  <b>{{ $key }}</b>:{{ $value }} <br>
	  {{ end }}
	</html>
	`)
	counterTempl, _ := template.New("printMetricsHTML").Parse(`
	<html>
	  <head>
	    <title>COUNTER</title>
	   <meta http-equiv="refresh" content="10" />
	  </head>
	  <h1><center>TYPE OF COUNTER METRICS</center></h1>
	  {{ range $key, $value := . }}
	  <b>{{ $key }}</b>:{{ $value }} <br>
	  {{ end }}
	</html>
	`)
	w.WriteHeader(http.StatusOK)
	gaugeTempl.Execute(w, gaugeStorage)
	counterTempl.Execute(w, counterStorage)
}

func main() {
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/update/{metricType}/{metricName}/{metricValue}/", httpPrint)
	r.Post("/update", httpPrintJSON)
	r.Post("/update/", httpPrintJSON)
	r.Get("/value/" + string(metrics.GaugeType) + "/{metricName}", httpPrintGaugeMetrics)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", httpPrint)
	r.Get("/value/" + string(metrics.GaugeType) + "/{metricName}", httpPrintCounterMetrics)
	r.Post("/value", httpPrintMetrics)
	r.Post("/value/", httpPrintMetrics)
	r.Get("/", httpPrintMetricsHTML)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		panic(err)
	}
}

