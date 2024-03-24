package handlers

import (
	"net/http/pprof"
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

	r.Use(middleware.WithLogging, middleware.GetCheck(flag), middleware.GzipHandler)
	attachPprof(r)
	r.Route("/", func(r chi.Router) {
		r.Get("/", Storagehandler.HandleGetAllMetrics().ServeHTTP)
		r.Get("/ping", storage.CheckDBConnection(storage.DB).ServeHTTP)
		r.Post("/update/", Storagehandler.HandleJSONPostMetrics(flag).ServeHTTP)
		r.Post("/updates/", Storagehandler.HandleMultipleJSONPostMetrics(flag).ServeHTTP)
		r.Post("/value/", Storagehandler.HandleJSONGetMetrics(flag).ServeHTTP)
		r.Post("/update/{type}/{name}/{value}", Storagehandler.HandlePostMetrics().ServeHTTP)
		r.Get("/value/{name}/{type}", Storagehandler.HandleGetMetrics().ServeHTTP)
	})

	return r
}

func attachPprof(mux *chi.Mux) {

	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

}
