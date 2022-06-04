package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"metrics/internal/server/config"
	"metrics/internal/server/handlers"
	"metrics/internal/server/storage"
	"net/http"
	"time"
)

type Server struct {
	startTime time.Time
	chiRouter chi.Router
}

func initRouter(memStatsStorage storage.MemStatsMemoryRepo) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	//Маршруты
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.PrintStatsValues(writer, request, memStatsStorage)
	})

	router.Route("/update", func(router chi.Router) {
		router.Post("/gauge/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateGaugePost(writer, request, memStatsStorage)
		})
		router.Post("/counter/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateCounterPost(writer, request, memStatsStorage)
		})
	})

	return router
}

func (server *Server) Run() {
	memStatsStorage := storage.NewMemStatsMemoryRepo()
	server.chiRouter = initRouter(memStatsStorage)

	fullHostAddr := fmt.Sprintf("%v:%v", config.ConfigHostname, config.ConfigPort)
	log.Fatal(http.ListenAndServe(fullHostAddr, server.chiRouter))
}
