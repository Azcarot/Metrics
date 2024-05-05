// Package agentconfigs - Пакет обработки флагов и переменных окружения
package agentconfigs

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
)

type AgentData struct {
	Pollint    time.Duration
	Reportint  time.Duration
	Addr       string
	HashKey    string
	RateLimit  int
	CryptoKey  string
	ConfigPath string
}

var agentFlags struct {
	pollinterval   int
	reportInterval int
	rateLimit      int
	flagAddr       string
	hashKey        string
	cryptoKey      string
	config         string
}

type AgentENV struct {
	Address    string `json:"address" env:"ADDRESS"`
	Key        string `json:"key" env:"KEY"`
	CryptoKey  string `json:"crypto_key" env:"CRYPTO_KEY"`
	PollInt    string `json:"poll_interval" env:"POLL_INTERVAL"`
	RepInt     string `json:"report_interval" env:"REPORT_INTERVAL"`
	RateLimit  int    `json:"rate_limit" env:"RATE_LIMIT"`
	ConfigPath string `env:"CONFIG"`
}

// parseFlags() получает флаги и сохраняет полученные в них данные в структуру AgentData
func parseFlags() *AgentData {
	var flagData AgentData
	flag.StringVar(&agentFlags.flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&agentFlags.hashKey, "k", "", "key to hash sha")
	flag.StringVar(&agentFlags.cryptoKey, "crypto-key", "", "path to file with public key")
	flag.StringVar(&agentFlags.config, "config", "", "path to config file")
	flag.IntVar(&agentFlags.pollinterval, "p", 2, "PollInterval")
	flag.IntVar(&agentFlags.reportInterval, "r", 10, "ReportInterval")
	flag.IntVar(&agentFlags.rateLimit, "l", 1, "amount of requests sended at one time")
	flag.Parse()
	envcfg := AgentENV{}
	isFlagSet := make(map[string]bool)
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}
	flagData.ConfigPath = agentFlags.config
	//Обработка флагов из файла и пересечения их с флагами командной строки и переменных окружения
	if flagData.ConfigPath != "" {
		if envcfg.ConfigPath != "" {
			flagData.ConfigPath = envcfg.ConfigPath
		}
		fileData, err := parseFile(flagData.ConfigPath)
		if err != nil {
			log.Fatal(err)
		}
		flagData = *fileData
		flag.Visit(func(f *flag.Flag) {
			isFlagSet[f.Name] = true
		})
		if len(flagData.Addr) == 0 || isFlagSet["a"] {
			flagData.Addr = agentFlags.flagAddr
			if envcfg.Address != "" {
				flagData.Addr = envcfg.Address
			}
		}
		if len(flagData.HashKey) == 0 || isFlagSet["k"] {
			flagData.HashKey = agentFlags.hashKey
			if envcfg.Key != "" {
				flagData.HashKey = envcfg.Key
			}
		}
		if len(flagData.CryptoKey) == 0 || isFlagSet["crypto-key"] {
			flagData.CryptoKey = agentFlags.cryptoKey
			if envcfg.CryptoKey != "" {
				flagData.CryptoKey = envcfg.CryptoKey
			}

		}
		if flagData.Pollint == 0 || isFlagSet["p"] {
			flagData.Pollint = time.Duration(agentFlags.pollinterval) * time.Second
			if envcfg.PollInt != "" {
				flagData.Pollint, _ = time.ParseDuration(envcfg.PollInt)
			}
		}

		if flagData.Reportint == 0 || isFlagSet["r"] {
			flagData.Reportint = time.Duration(agentFlags.reportInterval) * time.Second
			if envcfg.RepInt != "" {
				flagData.Reportint, _ = time.ParseDuration(envcfg.RepInt)
			}

		}
		if flagData.RateLimit == 0 || isFlagSet["l"] {
			flagData.RateLimit = agentFlags.rateLimit
			if envcfg.RateLimit > 0 {
				flagData.RateLimit = envcfg.RateLimit
			}
		}
		return &flagData
	}
	//Обработка флагов косандной строки и переменных окружения
	flagData.Addr = agentFlags.flagAddr
	flagData.Pollint = time.Duration(agentFlags.pollinterval)
	flagData.Reportint = time.Duration(agentFlags.reportInterval)
	flagData.HashKey = agentFlags.hashKey
	flagData.CryptoKey = agentFlags.cryptoKey
	flagData.ConfigPath = agentFlags.config

	if envcfg.CryptoKey != "" {
		flagData.CryptoKey = envcfg.CryptoKey
	}
	if len(envcfg.PollInt) > 0 {
		flagData.Pollint, _ = time.ParseDuration(envcfg.PollInt)
	}
	if len(envcfg.RepInt) > 0 {
		flagData.Reportint, _ = time.ParseDuration(envcfg.RepInt)
	}
	if len(envcfg.Key) > 0 {
		flagData.HashKey = envcfg.Key
	}

	if envcfg.RateLimit > 0 {
		flagData.RateLimit = envcfg.RateLimit
	}

	return &flagData
}

func parseFile(path string) (*AgentData, error) {
	file, err := os.Open(path)
	if err != nil {

		return nil, err
	}
	defer file.Close()
	var flagData AgentData
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON
	var envData AgentENV
	err = json.Unmarshal(data, &envData)
	if err != nil {
		return nil, err
	}
	flagData.Pollint, _ = time.ParseDuration(envData.PollInt)
	flagData.Reportint, _ = time.ParseDuration(envData.RepInt)
	flagData.Addr = envData.Address
	flagData.HashKey = envData.Key
	flagData.CryptoKey = envData.CryptoKey
	return &flagData, nil
}

// SetValues обрабатывает как флаги, так и переменные окружение
// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
func SetValues() *AgentData {
	flagData := parseFlags()

	return flagData
}
