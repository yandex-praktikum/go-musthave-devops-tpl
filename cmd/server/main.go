package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
        "html/template"

	"github.com/efrikin/go-musthave-devops-tpl/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	gaugeStorage = map[string]*metrics.Gauge{}
	counterStorage = map[string]*metrics.Counter{}
	mu      = sync.Mutex{}
)

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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{metricType}/{metricName}/{metricValue}/", httpPrint)
	// r.Get("/value/gauge/{metricName}", httpPrintGaugeMetrics)
	r.Get("/value/" + string(metrics.GaugeType) + "/{metricName}", httpPrintGaugeMetrics)
	// fmt.Println(string(metrics.GaugeType))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", httpPrint)
	r.Get("/value/" + string(metrics.GaugeType) + "/{metricName}", httpPrintCounterMetrics)
	r.Get("/", httpPrintMetricsHTML)
	http.ListenAndServe(":8080", r)
}

