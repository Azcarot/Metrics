package main

import (
	"flag"
	"time"

	"github.com/Azcarot/Metrics/cmd/agent/measure"
	"github.com/Azcarot/Metrics/cmd/server/handlers"
	"github.com/Azcarot/Metrics/cmd/types"
)

var agentFlags struct {
	pollinterval   int
	reportInterval int
	flagAddr       string
}

func parseFlags() {
	flag.StringVar(&agentFlags.flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&agentFlags.pollinterval, "p", 2, "PollInterval")
	flag.IntVar(&agentFlags.reportInterval, "r", 10, "PollInterval")
	flag.Parse()
}

func main() {
	parseFlags()
	var metric types.MemStorage
	counter := 0
	for {
		metric = measure.CollectMetrics(metric)
		time.Sleep(time.Duration(agentFlags.pollinterval) * time.Second)
		counter += 2
		if counter%agentFlags.reportInterval == 0 {
			urls := handlers.Makepath(metric, agentFlags.flagAddr)
			for _, url := range urls {
				handlers.PostMetrics(url)
			}

		}
	}
}
