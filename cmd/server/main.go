package main

import (
	"net/http"

	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/serverconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

func main() {

	flag := serverconfigs.ParseFlagsAndENV()
	storage.ConnectToDB(flag)
	storage.CreateTablesForMetrics(storage.DB)
	r := handlers.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	go handlers.GetSignal(server, flag)
	server.ListenAndServe()
	defer storage.DB.Close()

}
