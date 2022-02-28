package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *MetricsServer) startListener() {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)

	RegisterHandlers(mux, s.Cfg.MetricsStore)

	httpServer := &http.Server{
		Addr:    s.Cfg.ServerAddress,
		Handler: mux,
	}

	s.listener = httpServer

	log.Println(s.listener.ListenAndServe())
}

func (s *MetricsServer) stopListener() {
	err := s.listener.Shutdown(s.context)
	if err != nil {
		log.Printf("HTTP server ListenAndServe shut down: %q", err)
	}
}
