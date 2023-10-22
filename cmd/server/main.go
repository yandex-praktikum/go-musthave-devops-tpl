package main

import (
	//"fmt"
	//"os"
	"encoding/binary"
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	//"bufio"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type metricValue struct {
	val       [8]byte
	isCounter bool
}

type metric struct {
	mtype string
	value string
}

var (
	metricMap = make(map[string]metricValue)
	//metricStringMap = make(map[string]metric)
	lastCounterData int64
)

func int64ToBytes(value int64) [8]byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(value))
	return buf
}

func int64FromBytes(bytes []byte) int64 {
	i := binary.LittleEndian.Uint64(bytes)
	return int64(i)
}

func float64ToBytes(value float64) [8]byte {
	bits := math.Float64bits(value)
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], bits)
	return bytes
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func storeMetric(m map[string]metricValue) {
	options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	file, err := os.OpenFile("metrics.data", options, os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, m)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.RequestURI()
	//log.Println(q)
	reqMethod := r.Method
	log.Println("New request with method " + reqMethod + ", uri: " + q)
	if reqMethod != "POST" {
		outputMessage := "Only POST method is alload"
		log.Println("Wrong method for the handler. " + outputMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(outputMessage))
	}
	var m metricValue
	switch chi.URLParam(r, "metricType") {
	case "gauge":
		f, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64)
		if err != nil {
			log.Fatal(err)
		}
		m.val = float64ToBytes(f)
		//log.Println(float64FromBytes(m.val[:]))
		m.isCounter = false
		metricMap[chi.URLParam(r, "metricName")] = m
		w.WriteHeader(http.StatusOK)
		r.Body.Close()

	case "counter":
		c, err := strconv.ParseInt(chi.URLParam(r, "metricValue"), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		lastCounterData = lastCounterData + c // Change naming...
		m.val = int64ToBytes(lastCounterData)
		m.isCounter = true
		//log.Println(int64FromBytes(m.val[:]))
		metricMap[chi.URLParam(r, "metricName")] = m
		w.WriteHeader(http.StatusOK)
		r.Body.Close()

	default:
		log.Println("Type", chi.URLParam(r, "metricType"), "wrong")
		outputMessage := "Type " + chi.URLParam(r, "metricType") + " not supported, only [counter/gauge]"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(outputMessage))
		r.Body.Close()
	}

	//log.Println(metricMap)
	storeMetric(metricMap)
}

func handlerGetMetrics(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metricName")
	q := r.URL.RequestURI()
	log.Println(q)
	reqMethod := r.Method
	log.Println(reqMethod)

	if reqMethod != "GET" {

		log.Println("Method is wrong")
		outputMessage := "Only POST method is alload"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(outputMessage))

	}
	switch chi.URLParam(r, "metricType") {
	case "gauge":
		v, ok := metricMap[metricName]

		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v\n", float64FromBytes(v.val[:]))

		r.Body.Close()

	case "counter":
		v, ok := metricMap[metricName]

		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v\n", int64FromBytes(v.val[:]))

		r.Body.Close()

	default:
		log.Println("Type", chi.URLParam(r, "metricType"), "wrong")
		outputMessage := "Type " + chi.URLParam(r, "metricType") + " not supported, only [counter/gauge]"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(outputMessage))
		r.Body.Close()
	}

	log.Println(metricMap)
	options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	file, err := os.OpenFile("metrics.data", options, os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, metricMap)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func httpPrintAllMetrics(w http.ResponseWriter, _ *http.Request) {
	var stringMetricMap metric
	metricStringMap := make(map[string]metric)
	for k, v := range metricMap {
		if !v.isCounter {
			stringMetricMap.mtype = "gauge"
			stringMetricMap.value = strconv.FormatFloat((float64FromBytes(v.val[:])), 'f', -1, 64)
			metricStringMap[k] = stringMetricMap
		} else {
			stringMetricMap.mtype = "counter"
			stringMetricMap.value = strconv.FormatInt((int64FromBytes(v.val[:])), 10)
			metricStringMap[k] = stringMetricMap
		}

	}

	Templ, _ := template.New("printMetricsHTML").Parse(`
	<html>
	  <head>
	    <title>METRICS</title>
	   <meta http-equiv="refresh" content="3" />
	  </head>
	  <h1><center>METRICS</center></h1>
	  {{ range $key, $v := . }}
	  {{ $key }} : {{ $v }} <br>
	  {{ end }}
	</html>
	`)

	w.WriteHeader(http.StatusOK)
	Templ.Execute(w, metricStringMap)
}

func main() {
	port := ":8080"
	r := chi.NewRouter()
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handler)
	r.Get("/update/{metricType}/{metricName}/{metricValue}", handler)
	r.Get("/value/{metricType}/{metricName}", handlerGetMetrics)
	r.Get("/", httpPrintAllMetrics)

	// http.HandleFunc("/", handler)
	// err := http.ListenAndServe("localhost:8080", nil)
	log.Println("http server started on " + port)
	err := http.ListenAndServe("127.0.0.1"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
