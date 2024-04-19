// Основной серверный пакет. Ицидиирует связь с бд, создает роутер и слушает назначенный порт
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/serverconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version=%s\nBuild date =%s\nBuild commit =%s\n", buildVersion, buildDate, buildCommit)
	flag := serverconfigs.ParseFlagsAndENV()
	if flag.FlagDBAddr != "" {
		err := storage.NewConn(flag)
		if err != nil {
			panic(err)
		}
		storage.ST.CheckDBConnection()
		storage.ST.CreateTablesForMetrics()
		defer storage.DB.Close(context.Background())
	}
	r := handlers.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	//Сервер для pprof
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	go handlers.GetSignal(server, flag)
	server.ListenAndServe()

}
