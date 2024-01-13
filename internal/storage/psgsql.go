package storage

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectToDB(f Flags) {
	var err error
	ps := fmt.Sprintf(f.FlagDBAddr)
	DB, err = sql.Open("pgx", ps)
	if err != nil {
		panic(err)
	}

}

func CheckDBConnection(db *sql.DB) http.Handler {
	checkConnection := func(res http.ResponseWriter, req *http.Request) {
		err := db.Ping()
		result := (err == nil)
		if result {
			res.WriteHeader(http.StatusOK)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

	}
	return http.HandlerFunc(checkConnection)
}
