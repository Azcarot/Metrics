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
	query := `CREATE TABLE IF NOT EXISTS metrics (id int primary key, name text, type text, gauge_value double precision default NULL, counter_value int default NULL )`
	ctx := context.Background()
	//Проверяем, есть ли такая БД
	test_query := "SELECT datname FROM pg_catalog.pg_database WHERE datname = 'AzcarotPractics'"
	_, err := db.QueryContext(ctx, test_query)

	if err != nil {
		_, err := db.ExecContext(ctx, "CREATE DATABASE 'AzcarotPractics'")
		if err != nil {

			log.Printf("Error %s when creating product DB", err)

		}
	}
	_, err = db.ExecContext(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating product table", err)

	}

}

func WriteMetricsToPstgrs(db *sql.DB, data Metrics, t string) {
	switch t {
	case "gauge":
		db.ExecContext(context.Background(), "insert into metrics (name, type, gauge_value) values)", data.ID, data.MType, data.Value)
	case "counter":
		db.ExecContext(context.Background(), "insert into metrics (name, type, counter_value) values)", data.ID, data.MType, data.Delta)
	default:
		return
	}

}
