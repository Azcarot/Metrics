package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var DB *pgx.Conn
var ST PgxStorage

type PgxStorage interface {
	WriteMetricsToPstgrs(data Metrics, t string)
	BatchWriteToPstgrs(data []Metrics) error
	CheckDBConnection() http.Handler
	CreateTablesForMetrics()
}
type pgxConnTime struct {
	attempts          int
	timeBeforeAttempt int
}

type SQLStore struct {
	DB *pgx.Conn
}

func MakeStore(db *pgx.Conn) PgxStorage {
	return &SQLStore{
		DB: db,
	}
}

func NewConn(f Flags) error {
	var err error
	var attempts pgxConnTime
	attempts.attempts = 3
	attempts.timeBeforeAttempt = 1
	err = connectToDB(f)
	for err != nil {
		//если ошибка связи с бд, то это не эскпортируемый тип, отличный от PgError
		var pqErr *pgconn.PgError
		if errors.Is(err, pqErr) {
			return err

		}
		if attempts.attempts == 0 {
			return err
		}
		times := time.Duration(attempts.timeBeforeAttempt)
		time.Sleep(times * time.Second)
		attempts.attempts -= 1
		attempts.timeBeforeAttempt += 2
		err = connectToDB(f)

	}
	return nil
}

func connectToDB(f Flags) error {
	var err error
	ps := fmt.Sprintf(f.FlagDBAddr)
	DB, err = pgx.Connect(context.Background(), ps)
	ST = MakeStore(DB)
	return err
}

func (db *SQLStore) CheckDBConnection() http.Handler {
	checkConnection := func(res http.ResponseWriter, req *http.Request) {

		err := db.DB.Ping(context.Background())
		result := (err == nil)
		if result {
			res.WriteHeader(http.StatusOK)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

	}
	return http.HandlerFunc(checkConnection)
}

func (db *SQLStore) CreateTablesForMetrics() {
	query := `CREATE TABLE IF NOT EXISTS metrics (name text, type text, gauge_value double precision default NULL, counter_value int default NULL )`
	queryForFun := `DROP TABLE IF EXISTS metrics CASCADE`
	ctx := context.Background()
	_, err := db.DB.Exec(ctx, queryForFun)
	if err != nil {

		log.Printf("Error %s when Droping product table", err)

	}
	_, err = db.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating product table", err)

	}

}

func (db *SQLStore) WriteMetricsToPstgrs(data Metrics, t string) {
	ctx := context.Background()
	switch t {
	case "gauge":
		db.DB.Exec(ctx, `insert into metrics (name, type, gauge_value) values ($1, $2, $3);`, data.ID, data.MType, data.Value)
	case "counter":
		db.DB.Exec(ctx, `INSERT INTO metrics (name, type, counter_value) VALUES ($1, $2, $3);`, data.ID, data.MType, data.Delta)
	default:
		return
	}

}

func (db *SQLStore) BatchWriteToPstgrs(data []Metrics) error {
	copyCount, queryErr := db.DB.CopyFrom(
		context.Background(),
		pgx.Identifier{"metrics"},
		[]string{"name", "type", "counter_value", "gauge_value"},
		pgx.CopyFromSlice(len(data), func(i int) ([]interface{}, error) {
			return []interface{}{data[i].ID, data[i].MType, data[i].Delta, data[i].Value}, nil
		}),
	)
	if queryErr != nil {
		return queryErr
	}
	if int(copyCount) < len(data) {
		return fmt.Errorf("expected more rows in insert")
	}
	return nil
}
