package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"metrics/internal/server/storage"
	"net/http"
	"strconv"
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

	fmt.Println("Update gauge:")
	fmt.Printf("%v: %v\n\n", statName, statValue)
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

	fmt.Println("Inc counter:")
	fmt.Printf("%v: %v", statName, statValue)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Ok"))
}
