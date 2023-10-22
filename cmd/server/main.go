package main

import (
	"log"
	"net/http"

	// важно: путь до пакета handlers в вашем случае может быть другим

	"test_go/internal/router"
)

func main() {
	//ds := router.newDataStore()

	appRouter := router.New()

	serv := &http.Server{
		Addr:    ":8080",
		Handler: appRouter,
	}
	//http.HandleFunc("/update/", handlers.StatusHandler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(serv.ListenAndServe())
}
