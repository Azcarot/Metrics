package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

func CreateTablesForMetrics(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS metrics (name text, type text, gauge_value double precision default NULL, counter_value int default NULL )`
	queryForFun := `DROP TABLE IF EXISTS metrics CASCADE`
	ctx := context.Background()
	_, err := db.ExecContext(ctx, queryForFun)
	if err != nil {

		log.Printf("Error %s when Droping product table", err)

	}
	_, err = db.ExecContext(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating product table", err)

	}

}

func WriteMetricsToPstgrs(db *sql.DB, data Metrics, t string) {
	ctx := context.Background()
	switch t {
	case "gauge":
		db.ExecContext(ctx, `insert into metrics (name, type, gauge_value) values ($1, $2, $3);`, data.ID, data.MType, data.Value)
	case "counter":
		fmt.Println("TEST", data.ID)
		p, err := db.ExecContext(ctx, `INSERT INTO metrics (name, type, counter_value) VALUES ($1, $2, $3);`, data.ID, data.MType, data.Delta)
		fmt.Println("p ", p)
		fmt.Println("err ", err)
	default:
		return
	}

}
