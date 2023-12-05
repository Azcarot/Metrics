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
	urlmap := make(map[string]string)
	reqpaths := []string{"subaction", "mettype", "metname", "metvalue"}
	url := strings.Split(req.URL.Path, "/")
	for i := range reqpaths {
		urlmap[reqpaths[i]] = url[i+1]
		if reqpaths[i] == "subaction" && url[i+1] != "update" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(url[i+1]) == 0 && reqpaths[i] == "metname" {
			res.WriteHeader(http.StatusNotFound)
			return
		}
	}
	var mem types.MemStorage
	switch urlmap["mettype"] {

	case "gauge":
		if value, err := strconv.Atoi(urlmap["metvalue"]); err == nil {
			mem.Gaugemem = make(map[string]types.Gauge)
			mem.Gaugemem["metname"] = types.Gauge(value)
			res.WriteHeader(http.StatusOK)
			return
		} else {
			res.WriteHeader(http.StatusBadRequest)
		}
	case "counter":
		if value, err := strconv.Atoi(urlmap["metvalue"]); err == nil {
			mem.Countermem = make(map[string]types.Counter, 1)
			mem.Countermem["metname"] += types.Counter(value)
			res.WriteHeader(http.StatusOK)
			return
		} else {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

}

func HandleGetMetrics(res http.ResponseWriter, req *http.Request) {
	metric := strings.ToLower(chi.URLParam(req, "name") + ` ` + chi.URLParam(req, "type"))
	for _, c := range types.MetricNameTypes {
		if strings.ToLower(c) == metric {
			res.WriteHeader(http.StatusOK)
			res.Header().Add("Content-Type", "text/plain")
			io.WriteString(res, strconv.FormatFloat(rand.Float64(), 'g', -1, 64))
			return
		}
	}
	http.Error(res, "unknown metric: "+metric, http.StatusNotFound)
}
func HandleGetAllMetrics(res http.ResponseWriter, req *http.Request) {
	for _, c := range types.MetricNameTypes {
		io.WriteString(res, strings.ToLower(c)+strconv.FormatFloat(rand.Float64(), 'g', -1, 64))
		res.WriteHeader(http.StatusOK)
		res.Header().Add("Content-Type", "text/html")
		return
	}
}
