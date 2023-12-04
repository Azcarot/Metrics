package main

import (
	"agent/measure"
	"agent/postmetrics"
	"fmt"
	"net/http"
)

func main() {

	var metric measure.MemStorage
	metric = measure.GetMetrics(metric)
	urls := postmetrics.Makepath(metric)
	var resp *http.Response
	for _, url := range urls {
		resp = postmetrics.PostMetrics(url)
	}
	fmt.Sprintf("metric no longer supported", resp)

}
