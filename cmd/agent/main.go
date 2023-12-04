package main

import (
	"time"

	"github.com/Azcarot/Metrics/cmd/agent/measure"
	"github.com/Azcarot/Metrics/cmd/agent/postmetrics"
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
			for _, url := range urls {
				postmetrics.PostMetrics(url)
			}

		}
	}
}
