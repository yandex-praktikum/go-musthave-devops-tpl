package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"metrics/internal/server/storage"
	"net/http"
	"strconv"
	"strings"
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
	//если ключ не найден
	if strings.Contains(fmt.Sprint(err), "MemStat key not found") {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Not found"))
		return
	}

	//если другая ошибка
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
	//если ключ не найден
	if strings.Contains(fmt.Sprint(err), "MemStat key not found") {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Not found"))
		return
	}

	//если другая ошибка
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

func PrintStatsValues(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	htmlTemplate := `
<html>
    <head>
    <title></title>
    </head>
    <body>
        %v
    </body>
</html>`
	keyValuesHtml := ""

	for k, v := range memStatsStorage.GetDbSchema() {
		keyValuesHtml += fmt.Sprintf("<div><b>%v</b>: %v</div>", k, v)
	}

	htmlPage := fmt.Sprintf(htmlTemplate, keyValuesHtml)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(htmlPage))
}
