package handlers

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/go-chi/chi/v5"
)

func HandlePostMetrics(res http.ResponseWriter, req *http.Request) {
	metric := strings.ToLower(chi.URLParam(req, "type"))
	for _, c := range types.MetricNameTypes {
		if (strings.ToLower(c) == metric) && len(chi.URLParam(req, "name")) != 0 {
			switch strings.ToLower(chi.URLParam(req, "type")) {
			case "gauge":
				if value, err := strconv.Atoi(chi.URLParam(req, "value")); err == nil {
					types.Storage.Gaugemem[chi.URLParam(req, "name")] = types.Gauge(value)
					res.WriteHeader(http.StatusOK)
					return
				} else {
					res.WriteHeader(http.StatusBadRequest)
				}
			case "counter":
				if value, err := strconv.Atoi(chi.URLParam(req, "value")); err == nil {
					types.Storage.Countermem[chi.URLParam(req, "name")] = types.Counter(value)
					res.WriteHeader(http.StatusOK)
					return
				} else {
					res.WriteHeader(http.StatusBadRequest)
				}

			default:
				res.WriteHeader(http.StatusBadRequest)
			}
		}
	}
	res.WriteHeader(http.StatusBadRequest)
}

func HandleGetMetrics(res http.ResponseWriter, req *http.Request) {
	metric := strings.ToLower(chi.URLParam(req, "name") + ` ` + chi.URLParam(req, "type"))
	for _, c := range types.MetricNameTypes {
		if strings.ToLower(c) == metric {
			res.Header().Add("Content-Type", "text/plain")
			switch chi.URLParam(req, "type") {
			case "gauge":
				io.WriteString(res, strconv.FormatFloat(float64(types.Storage.Gaugemem[chi.URLParam(req, "name")]), 'g', -1, 64))
			case "counter":
				io.WriteString(res, strconv.Itoa(int(types.Storage.Countermem[chi.URLParam(req, "name")])))
			}

			return
		}
	}
	http.Error(res, "unknown metric: "+metric, http.StatusNotFound)
}

func HandleGetAllMetrics(res http.ResponseWriter, req *http.Request) {
	for i := range types.MetricNameTypes {
		io.WriteString(res, strings.ToLower(i)+` `+strconv.FormatFloat(rand.Float64(), 'g', -1, 64))

	}
	res.Header().Add("Content-Type", "text/html")
}
