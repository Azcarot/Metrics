package main

import (
	"agent/measure"
	"net/http"
	"strconv"
	"strings"
)

type MemInteractions interface {
	GetStoredMetrics() measure.MemStorage
	PostMetrics() bool
}

func GetStoredMetrics(m measure.MemStorage) measure.MemStorage {
	return m
}

func PostMetrics(m measure.MemStorage) bool {
	return true
}

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func PostReq(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	// requiredcontenttype := []string{"text/plain"}
	if method == http.MethodPost {
		// ctype := req.Header.Get("Content-type")
		// if ctype != requiredcontenttype[0] {
		//  	res.WriteHeader(http.StatusUnsupportedMediaType)
		//   	return
		//  }
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
		var mem measure.MemStorage
		switch urlmap["mettype"] {

		case "gauge":
			if value, err := strconv.Atoi(urlmap["metvalue"]); err == nil {
				mem.Gaugemem = make(map[string]measure.Gauge)
				mem.Gaugemem["metname"] = measure.Gauge(value)
				res.WriteHeader(http.StatusOK)
				return
			} else {
				res.WriteHeader(http.StatusBadRequest)
			}
		case "counter":
			if value, err := strconv.Atoi(urlmap["metvalue"]); err == nil {
				mem.Countermem = make(map[string]measure.Counter, 1)
				mem.Countermem["metname"] += measure.Counter(value)
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

	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(PostReq))
}

func main() {
	if runerr := run(); runerr != nil {
		panic(runerr)
	}

}
