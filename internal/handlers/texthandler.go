package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Azcarot/Metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

func (st *StorageHandler) HandlePostMetrics() http.Handler {
	postMetric := func(res http.ResponseWriter, req *http.Request) {
		if len(chi.URLParam(req, "name")) == 0 || len(chi.URLParam(req, "value")) == 0 || len(chi.URLParam(req, "type")) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		switch strings.ToLower(chi.URLParam(req, "type")) {
		case storage.CounterType, storage.GuageType:
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
	return http.HandlerFunc(postMetric)
}

func (st *StorageHandler) HandleGetMetrics() http.Handler {
	getMetric := func(res http.ResponseWriter, req *http.Request) {

		result, err := st.Storage.GetStoredMetrics(chi.URLParam(req, "type"), chi.URLParam(req, "name"))
		res.Header().Set("Content-Type", "text/html")
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			io.WriteString(res, result)
		}
	}
	return http.HandlerFunc(getMetric)
}

func (st *StorageHandler) HandleGetAllMetrics() http.Handler {
	getMetrics := func(res http.ResponseWriter, req *http.Request) {
		result := st.Storage.GetAllMetrics()
		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, result)
	}
	return http.HandlerFunc(getMetrics)
}
