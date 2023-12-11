package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/go-chi/chi/v5"
)

func (st *StorageHandler) HandlePostMetrics(res http.ResponseWriter, req *http.Request) {
	if len(chi.URLParam(req, "name")) == 0 || len(chi.URLParam(req, "value")) == 0 || len(chi.URLParam(req, "type")) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	switch strings.ToLower(chi.URLParam(req, "type")) {
	case types.CounterType, types.GuageType:
		err := st.Storage.StoreMetrics(chi.URLParam(req, "name"), strings.ToLower(chi.URLParam(req, "type")), chi.URLParam(req, "value"))
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func MakeRouter() *chi.Mux {
	storagehandler := StorageHandler{
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
	return r
}

type StorageHandler struct {
	Storage types.MemInteractions
}

func (st *StorageHandler) HandleGetMetrics(res http.ResponseWriter, req *http.Request) {
	result, err := st.Storage.GetStoredMetrics(chi.URLParam(req, "type"), chi.URLParam(req, "name"))
	res.Header().Add("Content-Type", "text/plain")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
	} else {
		io.WriteString(res, result)
	}
}

func (st *StorageHandler) HandleGetAllMetrics(res http.ResponseWriter, req *http.Request) {
	result := st.Storage.GetAllMetrics()
	io.WriteString(res, result)
	res.Header().Add("Content-Type", "text/html")
}
