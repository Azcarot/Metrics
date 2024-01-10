package handlers

import (
	"time"

	"github.com/Azcarot/Metrics/internal/middleware"
	"github.com/Azcarot/Metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type StorageHandler struct {
	Storage storage.MemInteractions
}

var Flag storage.Flags
var Storagehandler StorageHandler

func MakeRouter(flag storage.Flags) *chi.Mux {
	Storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}

	if flag.FlagRestore && len(flag.FlagFileStorage) > 0 {

		Storagehandler.Storage.ReadMetricsFromFile(flag.FlagFileStorage)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	if flag.FlagStoreInterval > 0 && len(flag.FlagFileStorage) > 0 {
		go func(name string) {
			reporttime := time.Duration(flag.FlagStoreInterval) * time.Second
			reporttimer := time.After(reporttime)
			for {
				<-reporttimer
				fullMetrics := Storagehandler.Storage.GetAllMetricsAsMetricType()
				for _, data := range fullMetrics {
					storage.WriteToFile(name, data)
				}
			}
		}(flag.FlagFileStorage)
	}
	defer logger.Sync()
	// делаем регистратор SugaredLogger
	middleware.Sugar = *logger.Sugar()
	r := chi.NewRouter()
	r.Use(middleware.WithLogging, middleware.GzipHandler)
	r.Route("/", func(r chi.Router) {
		r.Get("/", Storagehandler.HandleGetAllMetrics().ServeHTTP)
		r.Post("/update/", Storagehandler.HandleJSONPostMetrics(flag).ServeHTTP)
		r.Post("/value/", Storagehandler.HandleJSONGetMetrics().ServeHTTP)
		r.Post("/update/{type}/{name}/{value}", Storagehandler.HandlePostMetrics().ServeHTTP)
		r.Get("/value/{name}/{type}", Storagehandler.HandleGetMetrics().ServeHTTP)
	})
	return r
}
