package main

import (
	"net/http"

	"github.com/Azcarot/Metrics/cmd/server/handlers"
	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/go-chi/chi/v5"
)

func main() {
	storagehandler := &handlers.StorageHandler{
		Storage: &types.MemStorage{
			Gaugemem: make(map[string]types.Gauge), Countermem: make(map[string]types.Counter)},
	}
	r := chi.NewRouter()
	r.Use()
	r.Route("/", func(r chi.Router) {
		r.Get("/", http.HandlerFunc(storagehandler.HandleGetAllMetrics))
		r.Post("/update/{type}/{name}/{value}", http.HandlerFunc(storagehandler.HandlePostMetrics))
		r.Get("/value/{name}/{type}", http.HandlerFunc(storagehandler.HandleGetMetrics))
	})
	runerr := http.ListenAndServe(":8080", r)
	if runerr != nil {
		panic(runerr)
	}

}
