package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/go-chi/chi/v5"
)

func (st *StorageHandler) HandlePostMetrics(res http.ResponseWriter, req *http.Request) {
	if len(chi.URLParam(req, "name")) != 0 && len(chi.URLParam(req, "value")) != 0 && len(chi.URLParam(req, "type")) != 0 {
		if strings.ToLower(chi.URLParam(req, "type")) == "gauge" || strings.ToLower(chi.URLParam(req, "type")) == "counter" {
			err := st.Storage.StoreMetrics(chi.URLParam(req, "name"), strings.ToLower(chi.URLParam(req, "type")), chi.URLParam(req, "value"))
			if err == nil {
				res.WriteHeader(http.StatusOK)
				return
			}
		}
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusBadRequest)
	return
}

type StorageHandler struct {
	Storage types.MemInteractions
}

func (st *StorageHandler) HandleGetMetrics(res http.ResponseWriter, req *http.Request) {
	result, err := st.Storage.GetStoredMetrics(chi.URLParam(req, "type"), chi.URLParam(req, "name"))
	res.Header().Add("Content-Type", "text/plain")
	io.WriteString(res, result)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
	}
}

func (st *StorageHandler) HandleGetAllMetrics(res http.ResponseWriter, req *http.Request) {
	result := st.Storage.GetAllMetrics()
	io.WriteString(res, result)
	res.Header().Add("Content-Type", "text/html")
}
