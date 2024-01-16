package main

import (
	"context"
	"net/http"

	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/serverconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

func main() {

	flag := serverconfigs.ParseFlagsAndENV()
	if flag.FlagDBAddr != "" {
		err := storage.NewConn(flag)
		if err != nil {
			panic(err)
		}
		storage.CheckDBConnection(storage.DB)
		storage.CreateTablesForMetrics(storage.DB)
		defer storage.DB.Close(context.Background())
	}
	r := handlers.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	go handlers.GetSignal(server, flag)
	server.ListenAndServe()

}
