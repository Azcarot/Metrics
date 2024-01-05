package main

import (
	"fmt"
	"net/http"

	"github.com/Azcarot/Metrics/cmd/server/handlers"
)

func main() {
	flag := handlers.ParseFlagsAndENV()
	fmt.Println(flag)
	r := handlers.MakeRouter(flag)
	http.ListenAndServe(flag.FlagAddr, r)

}
