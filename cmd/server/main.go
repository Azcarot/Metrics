package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Azcarot/Metrics/cmd/server/handlers"
	"github.com/caarlos0/env/v6"
)

var flagAddr string

type serverENV struct {
	Address string `env:"ADDRESS"`
}

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}

func main() {
	parseFlags()
	var envcfg serverENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}
	if envcfg.Address != "" {
		flagAddr = envcfg.Address
	}
	r := handlers.MakeRouter()
	http.ListenAndServe(flagAddr, r)

}
