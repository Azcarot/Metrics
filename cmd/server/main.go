package main

import (
	"context"
	"log"
	"net/http"

	_ "net/http/pprof"

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
		storage.ST.CheckDBConnection()
		storage.ST.CreateTablesForMetrics()
		defer storage.DB.Close(context.Background())
	}
	r := handlers.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	go handlers.GetSignal(server, flag)
	server.ListenAndServe()

	// // // создаём файл журнала профилирования памяти
	// fmem, err := os.Create(`.\profiles\base2.pprof`)
	// if err != nil {
	// 	panic(err)
	// }
	// defer fmem.Close()
	// runtime.GC() // получаем статистику по использованию памяти
	// if err := pprof.WriteHeapProfile(fmem); err != nil {
	// 	panic(err)
	// }
}
