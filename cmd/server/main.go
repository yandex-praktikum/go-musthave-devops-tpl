package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/efrikin/go-musthave-devops-tpl/internal/metrics"
)

var (
	gauge   = metrics.Gauge{}
	counter = metrics.Counter{}
	storage = map[string]interface{}{}
	mu      = sync.Mutex{}
)

func httpPrint(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "I need method POST", http.StatusNotFound)
		return
	}
	u := strings.Split(strings.TrimRight(r.URL.Path, "/"), "/")

	if len(u) != 5 {
		http.Error(w, "Error", http.StatusNotFound)
		return
	}
	if u[1] != "update" {
		http.Error(w, "Error", http.StatusNotFound)
		return
	}
	if u[2] == gauge.Type() {
		_, err := strconv.ParseFloat(u[4], 64)
		if err != nil {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		w.WriteHeader(http.StatusOK)
		storage[u[3]] = u[4]
		return
	}
	if u[2] == counter.Type() {
		mu.Lock()
		defer mu.Unlock()
		var tmp, ok = storage[u[3]]
		if !ok {
			tmp = int64(0)
		}
		tmp2, err := strconv.ParseInt(u[4], 10, 64)
		if err != nil {
			http.Error(w, "Bad convenrt int to string", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		storage[u[3]] = tmp.(int64) + tmp2
		return
	}
	http.Error(w, "Error", http.StatusNotImplemented)
}

func httpPrintMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v\n", storage)
}

func main() {
	http.HandleFunc("/", httpPrint)
	http.HandleFunc("/metrics", httpPrintMetrics)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
