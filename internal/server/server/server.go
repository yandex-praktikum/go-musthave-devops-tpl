package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"metrics/internal/server/config"
	"net/http"
	"time"
)

type Server struct {
	startTime time.Time
	chiRouter chi.Router
}

func initRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// создадим суброутер, который будет содержать две функции
	router.Route("/update", func(router chi.Router) {
		router.Post("/gauge/{statName}/{statValue}", func(rw http.ResponseWriter, request *http.Request) {
			statName := chi.URLParam(request, "statName")
			statValue := chi.URLParam(request, "statName")
			fmt.Println("Update gauge:")
			fmt.Printf("%v: %v", statName, statValue)

			rw.Write([]byte(" "))
		})

		router.Post("/counter/{statName}/{statValue}", func(rw http.ResponseWriter, request *http.Request) {
			statName := chi.URLParam(request, "statName")
			statValue := chi.URLParam(request, "statName")
			fmt.Println("Update gauge:")
			fmt.Printf("%v: %v", statName, statValue)

			rw.Write([]byte(" "))
		})
	})

	return router
}

func (server *Server) Run() {
	server.chiRouter = initRouter()

	fullHostAddr := fmt.Sprintf("%v:%v", config.ConfigHostname, config.ConfigPort)
	log.Fatal(http.ListenAndServe(fullHostAddr, server.chiRouter))
}
