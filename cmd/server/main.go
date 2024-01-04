package main

import (
	"fmt"
	"net/http"

	"github.com/Azcarot/Metrics/cmd/server/handlers"
)

func main() {
	flag := handlers.ParseFlagsAndENV()

	r := handlers.MakeRouter(flag)
	fmt.Println(flag)
	http.ListenAndServe(flag.FlagAddr, r)

}
