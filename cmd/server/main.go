package main

import (
	"database/sql"
	"net/http"

	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/serverconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

var DB *sql.DB

func main() {

	flag := serverconfigs.ParseFlagsAndENV()
	r := handlers.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	go handlers.GetSignal(server, flag)
	server.ListenAndServe()
	defer storage.DB.Close()

}
