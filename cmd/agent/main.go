package main

import (
	"time"

	"github.com/Azcarot/Metrics/cmd/agent/measure"
	"github.com/Azcarot/Metrics/cmd/server/handlers"
	"github.com/Azcarot/Metrics/cmd/types"
)

func main() {

	var metric types.MemStorage
	counter := 0
	for {
		metric = measure.CollectMetrics(metric)
		time.Sleep(2 * time.Second)
		counter += 2
		if counter%10 == 0 {
			urls := handlers.Makepath(metric)
			for _, url := range urls {
				handlers.PostMetrics(url)
			}

		}
	}
}
