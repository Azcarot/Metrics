package main

import (
	"agent/measure"
	"agent/postmetrics"
	"fmt"
	"net/http"
	"time"
)

func main() {

	var metric measure.MemStorage
	counter := 0
	for {
		metric = measure.GetMetrics(metric)
		time.Sleep(2 * time.Second)
		counter += 2
		if counter%10 == 0 {
			urls := postmetrics.Makepath(metric)
			var resp *http.Response
			for _, url := range urls {
				resp = postmetrics.PostMetrics(url)
			}
			fmt.Sprintf("metric no longer supported", resp)
		}
	}
}
