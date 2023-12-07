package main

import (
	"flag"
	"log"
	"time"

	"github.com/Azcarot/Metrics/cmd/agent/measure"
	"github.com/Azcarot/Metrics/cmd/server/handlers"
	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/caarlos0/env/v6"
)

var agentData struct {
	pollint   time.Duration
	reportint int
	addr      string
}

var agentFlags struct {
	pollinterval   int
	reportInterval int
	flagAddr       string
}

type agentENV struct {
	Address string        `env:"ADDRESS"`
	PollInt time.Duration `env: "POLL_INTERVAL"`
	RepInt  int           `env: "REPORT_INTERVAL"`
}

func parseFlags() {
	flag.StringVar(&agentFlags.flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&agentFlags.pollinterval, "p", 2, "PollInterval")
	flag.IntVar(&agentFlags.reportInterval, "r", 10, "PollInterval")
	flag.Parse()
}

// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
func setValues() {
	parseFlags()
	var envcfg agentENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}
	agentData.pollint = time.Duration(agentFlags.pollinterval)
	agentData.reportint = agentFlags.reportInterval
	agentData.addr = agentFlags.flagAddr
	if envcfg.Address != "" {
		agentData.addr = envcfg.Address
	}
	if envcfg.PollInt > 0 {
		agentData.pollint = envcfg.PollInt
	}
	if envcfg.RepInt > 0 {
		agentData.reportint = envcfg.RepInt
	}
}

func main() {
	setValues()
	var metric types.MemStorage
	counter := 0
	for {
		metric = measure.CollectMetrics(metric)
		time.Sleep(time.Duration(agentData.pollint) * time.Second)
		counter += 2
		if counter%agentData.reportint == 0 {
			urls := handlers.Makepath(metric, agentData.addr)
			for _, url := range urls {
				handlers.PostMetrics(url)
			}

		}
	}
}
