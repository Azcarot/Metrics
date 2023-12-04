package main

import (
	"net/http"
	"strconv"
	"strings"
)

type gauge int
type counter int
type MemStorage struct {
	floatmem   map[string]float64
	gaugemem   map[string]gauge
	countermem map[string][]counter
	int64mem   map[string][]int64
}

type MemInteractions interface {
	GetStoredMetrics() MemStorage
	PostMetrics() bool
}

func GetStoredMetrics(m MemStorage) MemStorage {
	return m
}

func PostMetrics(m MemStorage) bool {
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
		var mem MemStorage
		switch urlmap["mettype"] {
		case "float64":
			if value, err := strconv.ParseFloat(urlmap["metvalue"], 64); err == nil {
				mem.floatmem = make(map[string]float64)
				mem.floatmem["metname"] = value
				res.WriteHeader(http.StatusOK)
				return
			} else {
				res.WriteHeader(http.StatusBadRequest)
			}
		case "gauge":
			if value, err := strconv.Atoi(urlmap["metvalue"]); err == nil {
				mem.gaugemem = make(map[string]gauge)
				mem.gaugemem["metname"] = gauge(value)
				res.WriteHeader(http.StatusOK)
				return
			} else {
				res.WriteHeader(http.StatusBadRequest)
			}
		case "counter":
			if value, err := strconv.Atoi(urlmap["metvalue"]); err == nil {
				mem.countermem = make(map[string][]counter, 1)
				mem.countermem["metname"] = append(mem.countermem["metname"], counter(value))
				res.WriteHeader(http.StatusOK)
				return
			} else {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		case "int64":
			if value, err := strconv.ParseInt(urlmap["metvalue"], 0, 64); err == nil {
				mem.int64mem = make(map[string][]int64, 1)
				mem.int64mem["metname"] = append(mem.int64mem["metname"], value)
				res.WriteHeader(http.StatusOK)
				return
			} else {
				res.WriteHeader(http.StatusBadRequest)
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
