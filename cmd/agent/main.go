package main

import (
	"flag"
	"log"
	"time"

	"github.com/Azcarot/Metrics/cmd/agent/measure"
	"github.com/Azcarot/Metrics/cmd/server/handlers"
	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/caarlos0/env/v10"
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

type AgentENV struct {
	Address string `env:"ADDRESS"`
	PollInt int    `env:"POLL_INTERVAL"`
	RepInt  int    `env:"REPORT_INTERVAL"`
}

func parseFlags() {
	flag.StringVar(&agentFlags.flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&agentFlags.pollinterval, "p", 2, "PollInterval")
	flag.IntVar(&agentFlags.reportInterval, "r", 10, "ReportInterval")
	flag.Parse()
}

// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
func setValues() {
	parseFlags()
	envcfg := AgentENV{}
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

	if int(envcfg.PollInt) > 0 {
		agentData.pollint = time.Duration(envcfg.PollInt)
	}
	if envcfg.RepInt > 0 {
		agentData.reportint = envcfg.RepInt
	}
}

func main() {
	setValues()
	var metric types.MemStorage
	sleeptime := time.Duration(agentData.pollint) * time.Second
	reporttime := time.Duration(agentData.reportint) * time.Second
	reporttimer := time.After(reporttime)
	for {
		select {
		case <-reporttimer:
			urls := handlers.Makepath(metric, agentData.addr)
			for _, url := range urls {
				handlers.PostMetrics(url)
			}
			reporttimer = time.After(reporttime)
		default:
			metric = measure.CollectMetrics(metric)
			time.Sleep(sleeptime)
		}
	}
}
