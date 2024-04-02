package storage

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	var f Flags
	f.FlagDBAddr = "host='localhost' user='postgres' password='12345' sslmode=disable"
	connectToDB(f)
	ST = MakeStore(DB)
	dbName := "testdb"
	ctx := context.Background()
	_, err = DB.Exec(ctx, "DROP DATABASE IF EXISTS "+dbName)
	if err != nil {
		//handle the error
		log.Fatal(err)
	}
	_, err = DB.Exec(ctx, "create database "+dbName)
	if err != nil {
		//handle the error
		log.Fatal(err)
	}
	ST.CreateTablesForMetrics()
	// truncates all test data after the tests are run
	defer func() {

		_, _ = DB.Exec(ctx, fmt.Sprintf("DELETE FROM %s", "metrics"))

		DB.Close(ctx)
	}()

	return m.Run(), nil
}

func TestSQLStore_WriteMetricsToPstgrs(t *testing.T) {
	type args struct {
		data Metrics
	}
	value := rand.Float64()
	counter := rand.Int63n(100)
	tests := []struct {
		name string
		args args
	}{
		{name: "NormalGauge", args: args{data: Metrics{ID: "111", MType: GuageType, Value: &value}}},
		{name: "NormalCounter", args: args{data: Metrics{ID: "222", MType: CounterType, Delta: &counter}}},
		{name: "wrongData", args: args{data: Metrics{ID: "333", MType: "ho"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ST.WriteMetricsToPstgrs(tt.args.data)
		})
	}
}

func TestSQLStore_BatchWriteToPstgrs(t *testing.T) {
	type args struct {
		data []Metrics
	}
	value := rand.Float64()
	counter := rand.Int63n(100)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "normal data",
			args: args{data: []Metrics{{ID: "111", MType: GuageType, Value: &value},
				{ID: "222", MType: CounterType, Delta: &counter}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ST.BatchWriteToPstgrs(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SQLStore.BatchWriteToPstgrs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
