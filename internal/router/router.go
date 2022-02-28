package router

import (
	"net/http"
	"strconv"
	"test_go/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Data struct {
	key   string
	value int64
	types string
}

type DataStore struct {
	datastore map[string]Data
}

func newDataStore() DataStore {
	return DataStore{make(map[string]Data)}
}

//var m = map[string]float64{}
var jsonData = []byte(`{
	"name": "morpheus",
	"job": "leader"
}`)

var m = map[string]int64{}
var ds = newDataStore()

/*
Puts the key and value in the DataStore.
If the key (k) already exists will replace with the provided value (v).
*/
func (ds *DataStore) put(k, v, t string) {
	//b2, _ := strconv.ParseFloat(v, 32)
	b2, _ := strconv.ParseInt(v, 10, 64)
	if t == "gauge" {

		dx := Data{key: k, value: b2}
		ds.datastore[k] = dx
	}
	if t == "counter" {
		//m := make(map[string]float64)
		m[k] = b2

	}
}

func New() chi.Router {
	//ds := newDataStore()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			for k, v := range ds.datastore {
				//fmt.Println(k, " ", v)

				//rw.Write([]byte(carFunc(carID)))
				//"<h1>Hello, World</h1>"
				//rw.Write([]byte(fmt.Sprintf("%s %v\n", k, v.value)))
				rw.Write([]byte("<h1>" + k + " " + strconv.FormatUint(uint64(v.value), 10) + "</h1>"))
			}

		})
		r.Post("/update/{name}/{tupe}/{stats}", func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")
			tupe := chi.URLParam(r, "tupe")
			stats := chi.URLParam(r, "stats")
			w.Write([]byte(name + "-" + tupe + "-" + stats))
			ds.put(name, stats, tupe)
			//w.Write([]byte(fmt.Sprintf("Query string values: %s", ds)))
			//return
		})
	})
	r.MethodFunc("GET", "/update/", handlers.StatusHandler)
	//r.MethodFunc("POST", "/update/", handlers.StatusHandler)
	return r
}
