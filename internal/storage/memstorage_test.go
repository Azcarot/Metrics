package storage

import (
	"reflect"
	"testing"
)

func TestMemStorage_GetAllMetrics(t *testing.T) {
	gaugeValue := float64(333)
	want := "Metrics name: " + "333" + "\n" + "Metrics value: " + "333"
	tests := []struct {
		name      string
		m         *MemStorage
		writeData bool
		want      string
		data      Metrics
	}{
		{name: "Нет записанных данных", want: "", writeData: false},
		{name: "Есть данные", want: want, writeData: true, data: Metrics{ID: "333", MType: GuageType, Value: &gaugeValue}},
	}
	type StorageHandler struct {
		Storage MemInteractions
	}
	Storagehandler := StorageHandler{
		Storage: &MemStorage{
			Gaugemem: make(map[string]Gauge), Countermem: make(map[string]Counter)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.writeData {
				Storagehandler.Storage.StoreMetrics(tt.data)
			}
			if got := Storagehandler.Storage.GetAllMetrics(); got != tt.want {
				t.Errorf("MemStorage.GetAllMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_StoreMetrics(t *testing.T) {
	type args struct {
		data Metrics
	}
	gaugeValue := float64(333)
	counterValue := int64(333)
	tests := []struct {
		name    string
		m       *MemStorage
		args    args
		wantErr bool
	}{
		{name: "no data", wantErr: true}, {name: "NormalData", args: args{data: Metrics{ID: "333", MType: GuageType, Value: &gaugeValue}}, wantErr: false},
		{name: "NormalDataDelta", args: args{data: Metrics{ID: "444", MType: CounterType, Delta: &counterValue}}, wantErr: false}, {name: "FalseType", args: args{data: Metrics{ID: "444", MType: "FALSETYPE", Delta: &counterValue}}, wantErr: true},
	}
	type StorageHandler struct {
		Storage MemInteractions
	}
	Storagehandler := StorageHandler{
		Storage: &MemStorage{
			Gaugemem: make(map[string]Gauge), Countermem: make(map[string]Counter)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Storagehandler.Storage.StoreMetrics(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.StoreMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_GetStoredMetrics(t *testing.T) {
	type args struct {
		n string
		t string
	}
	type StorageHandler struct {
		Storage MemInteractions
	}
	gaugeValue := float64(333)
	counterValue := int64(333)
	Storagehandler := StorageHandler{
		Storage: &MemStorage{
			Gaugemem: make(map[string]Gauge), Countermem: make(map[string]Counter)},
	}
	tests := []struct {
		name    string
		m       *MemStorage
		args    args
		data    Metrics
		want    string
		wantErr bool
	}{
		{name: "No data", args: args{n: "", t: ""}, wantErr: true},
		{name: "gauge", args: args{n: "222", t: GuageType},
			data: Metrics{ID: "222", MType: GuageType, Value: &gaugeValue},
			want: "333", wantErr: false},
		{name: "counter", args: args{n: "888", t: CounterType},
			data: Metrics{ID: "888", MType: CounterType, Delta: &counterValue},
			want: "333", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				Storagehandler.Storage.StoreMetrics(tt.data)
			}
			got, err := Storagehandler.Storage.GetStoredMetrics(tt.args.n, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetStoredMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.GetStoredMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetAllMetricsAsMetricType(t *testing.T) {
	type StorageHandler struct {
		Storage MemInteractions
	}
	gaugeValue := float64(333)
	counterValue := int64(333)
	var empty []Metrics
	tests := []struct {
		name       string
		m          *MemStorage
		dataToSave Metrics
		want       []Metrics
	}{
		{name: "no saved data", want: empty},
		{name: "saved gauge",
			dataToSave: Metrics{ID: "222", MType: GuageType, Value: &gaugeValue},
			want:       []Metrics{{ID: "222", MType: GuageType, Value: &gaugeValue}}},
		{name: "saved gauge and counter",
			dataToSave: Metrics{ID: "333", MType: CounterType, Delta: &counterValue},
			want: []Metrics{{ID: "222", MType: GuageType, Value: &gaugeValue},
				{ID: "333", MType: CounterType, Delta: &counterValue}}}}
	Storagehandler := StorageHandler{
		Storage: &MemStorage{
			Gaugemem: make(map[string]Gauge), Countermem: make(map[string]Counter)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name != "no saved data" {
				Storagehandler.Storage.StoreMetrics(tt.dataToSave)
			}
			if got := Storagehandler.Storage.GetAllMetricsAsMetricType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetAllMetricsAsMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}
