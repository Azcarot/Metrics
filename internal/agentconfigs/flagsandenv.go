package agentconfigs

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
)

type AgentData struct {
	Pollint   time.Duration
	Reportint int
	Addr      string
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

func parseFlags() *AgentData {
	var flagData AgentData
	flag.StringVar(&agentFlags.flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&agentFlags.pollinterval, "p", 2, "PollInterval")
	flag.IntVar(&agentFlags.reportInterval, "r", 10, "ReportInterval")
	flag.Parse()
	flagData.Pollint = time.Duration(agentFlags.pollinterval)
	flagData.Reportint = agentFlags.reportInterval
	flagData.Addr = agentFlags.flagAddr
	return &flagData
}

// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
func SetValues() *AgentData {
	flagData := parseFlags()
	envcfg := AgentENV{}
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	if envcfg.Address != "" {
		flagData.Addr = envcfg.Address
	}

	if int(envcfg.PollInt) > 0 {
		flagData.Pollint = time.Duration(envcfg.PollInt)
	}
	if envcfg.RepInt > 0 {
		flagData.Reportint = envcfg.RepInt
	}
	return flagData
}
