package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	_ "strings"
	"sync"

	"github.com/efrikin/go-musthave-devops-tpl/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	gauge   = metrics.Gauge{}
	counter = metrics.Counter{}
	storage = map[string]interface{}{}
	mu      = sync.Mutex{}
)

func httpPrint(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType == gauge.Type() {
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		w.WriteHeader(http.StatusOK)
		storage[metricName] = metricValue
		return
	}

	if metricType == counter.Type() {
		mu.Lock()
		defer mu.Unlock()
		var tmp, ok = storage[metricName]
		if !ok {
			tmp = int64(0)
		}
		tmp2, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Bad convenrt int to string", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		storage[metricName] = tmp.(int64) + tmp2
		return
	}
	http.Error(w, "Error", http.StatusNotImplemented)
}

func httpPrintMetrics(w http.ResponseWriter, r *http.Request) {

	metricName := chi.URLParam(r, "metricName")
	val, ok := storage[metricName]

	if !ok {
		http.Error(w,"Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v\n", val)
}

func httpPrintMetricsHTML(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `<html>
					<head>
					<title>METRICS</title>
					<meta http-equiv="refresh" content="10" />
					</head>
					<h1><center>METRICS</center></h1>`)
	keys := make([]string, 0, len(storage))
	for k := range storage {
		keys = append(keys, k)

	}
	sort.Strings(keys)

	for _, v := range keys {
		fmt.Fprintf(w, "<p>%s=%v</p>", v, storage[v])
	}

	fmt.Fprintf(w, "</html>")

}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Post("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })
	r.Post("/update/{metricType}/{metricName}/{metricValue}/", httpPrint)
	r.Get("/value/{metricType}/{metricName}/", httpPrintMetrics)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", httpPrint)
	r.Get("/value/{metricType}/{metricName}", httpPrintMetrics)
	r.Get("/", httpPrintMetricsHTML)
	http.ListenAndServe(":8080", r)
}
