package main

import (
	"net/http"

	"github.com/Azcarot/Metrics/cmd/server/handlers"
)

func main() {

	flag := handlers.ParseFlagsAndENV()
	r := handlers.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	go handlers.GetSignal(server, flag)
	server.ListenAndServe()

}
