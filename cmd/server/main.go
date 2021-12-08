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
	storage = map[string]interface{}{}
	mu      = sync.Mutex{}
)

func httpPrint(w http.ResponseWriter, r *http.Request) {

	metricType := metrics.MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType == metrics.GaugeType {
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

	if metricType == metrics.CounterType {
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
	templ, _ := template.New("printMetricsHTML").Parse(`
	<html>
	  <head>
	    <title>METRICS</title>
	   <meta http-equiv="refresh" content="10" />
	  </head>
	  <h1><center>METRICS</center></h1>
	  {{ range $key, $value := . }}
	  <b>{{ $key | printf "%s" }}</b>:{{ $value | printf "\t%v"}} <br>
	  {{ end }}
	</html>
	`)
	w.WriteHeader(http.StatusOK)

	templ.Execute(w, storage)

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
