package main

import (
	"net/http"

	"github.com/Azcarot/Metrics/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
)

func run() error {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", http.HandlerFunc(handlers.HandleGetAllMetrics))
		r.Post("/update/{type}/{name}/{value}", http.HandlerFunc(handlers.HandlePostMetrics))
		r.Get("/value/{name}/{type}", http.HandlerFunc(handlers.HandleGetMetrics))
	})
	return http.ListenAndServe(":8080", r)
}

func main() {
	if runerr := run(); runerr != nil {
		panic(runerr)
	}

}
