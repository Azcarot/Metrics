package main

import (
	"net/http"

	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/serverconfigs"
)

func main() {

	flag := serverconfigs.ParseFlagsAndENV()
	r := handlers.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	go handlers.GetSignal(server, flag)
	server.ListenAndServe()

}
