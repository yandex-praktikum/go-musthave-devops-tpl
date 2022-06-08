package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"metrics/internal/server/storage"
)

func UpdateGaugePost(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	statName := chi.URLParam(request, "statName")
	statValue := chi.URLParam(request, "statValue")
	statValueInt, err := strconv.ParseFloat(statValue, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Bad request"))
		return
	}

	err = memStatsStorage.UpdateGaugeValue(statName, statValueInt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Server error"))
		return
	}

	log.Println("Update gauge:")
	log.Printf("%v: %v\n", statName, statValue)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Ok"))
}

func UpdateCounterPost(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	statName := chi.URLParam(request, "statName")
	statValue := chi.URLParam(request, "statValue")
	statCounterValue, err := strconv.ParseInt(statValue, 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	err = memStatsStorage.UpdateCounterValue(statName, statCounterValue)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	log.Println("Inc counter:")
	log.Printf("%v: %v\n", statName, statValue)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Ok"))
}

func UpdateNotImplementedPost(rw http.ResponseWriter, request *http.Request) {
	log.Println("Update not implemented statType")
	rw.WriteHeader(http.StatusNotImplemented)
	rw.Write([]byte("Not implemented"))
}

func PrintStatsValues(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	htmlTemplate := `
<html>
    <head>
    <title></title>
    </head>
    <body>
		<h3 class="keyvalues-header">All values: </h3>
        %v
    </body>
</html>`
	keyValuesHTML := ""

	for k, v := range memStatsStorage.GetDBSchema() {
		keyValuesHTML += fmt.Sprintf("<div><b>%v</b>: %v</div>", k, v)
	}

	htmlPage := fmt.Sprintf(htmlTemplate, keyValuesHTML)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(htmlPage))
}

func PrintStatValue(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	statName := chi.URLParam(request, "statName")
	statValue, err := memStatsStorage.ReadValue(statName)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Unknown statName"))
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(statValue))
}
