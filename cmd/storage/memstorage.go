package storage

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

const GuageType = "gauge"
const CounterType = "counter"

type Gauge float64
type Counter int64
type MemStorage struct {
	Gaugemem   map[string]Gauge
	Countermem map[string]Counter
}

type Flags struct {
	FlagAddr          string
	FlagStoreInterval int
	FlagFileStorage   string
	FlagRestore       bool
}

type ServerENV struct {
	Address       string `env:"ADDRESS"`
	StoreInterval string `env:"STORE_INTERVAL"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
}

type AllowedMetrics struct {
	Name string
}
type MemInteractions interface {
	GetStoredMetrics(string, string) (string, error)
	StoreMetrics(string, string, string) error
	GetAllMetrics() string
	ReadMetricsFromFile(string)
	GetAllMetricsAsMetricType() []Metrics
}

func (m *MemStorage) StoreMetrics(n string, t string, v string) error {
	switch t {
	case GuageType:
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		m.Gaugemem[n] = Gauge(value)
	case CounterType:
		value, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		m.Countermem[n] += Counter(value)
	}
	return nil
}

func (m *MemStorage) ReadMetricsFromFile(filename string) {
	Consumer, err := NewConsumer(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer Consumer.Close()
	metrics, err := Consumer.ReadEvent()
	if err != nil {
		log.Fatal(err)
	}
	if len(metrics.MType) > 0 {
		switch strings.ToLower(metrics.MType) {
		case "gauge":
			m.Gaugemem[metrics.ID] = Gauge(*metrics.Value)
		case "counter":
			m.Countermem[metrics.ID] = Counter(*metrics.Delta)

		}
	}

}

func (m *MemStorage) GetAllMetricsAsMetricType() []Metrics {
	var FinalData []Metrics
	for n, v := range m.Gaugemem {
		var data Metrics
		data.ID = n
		*data.Value = float64(v)
		data.MType = "gauge"
		FinalData = append(FinalData, data)
	}
	for n, v := range m.Countermem {
		var data Metrics
		data.ID = n
		*data.Delta = int64(v)
		data.MType = "counter"
		FinalData = append(FinalData, data)
	}
	return FinalData
}

func (m *MemStorage) GetAllMetrics() string {
	var result string
	for n, v := range m.Gaugemem {
		value := strconv.FormatFloat(float64(v), 'g', -1, 64)
		result += "Metrics name: " + n + "\n" + "Metrics value: " + value
	}
	for n, v := range m.Countermem {
		value := strconv.Itoa(int(v))
		result += "Metrics name: " + n + "\n" + "Metrics value: " + value
	}

	return result
}

func (m *MemStorage) GetStoredMetrics(n string, t string) (string, error) {
	var result string
	var err error
	switch t {
	case GuageType:
		val, ok := m.Gaugemem[n]

		if ok {
			result = strconv.FormatFloat(float64(val), 'g', -1, 64)
		} else {
			err = errors.New("no metric stored")
		}

	case CounterType:
		val, ok := m.Countermem[n]
		if ok {
			result = strconv.Itoa(int(val))
		} else {
			err = errors.New("no metric stored")
		}
	default:
		return "", errors.New("wrong type")
	}
	return result, err
}

var MetricNameTypes = map[string]string{
	"PollCount":     CounterType,
	"Frees":         GuageType,
	"GCSys":         GuageType,
	"HeapAlloc":     GuageType,
	"HeapIdle":      GuageType,
	"HeapInuse":     GuageType,
	"HeapObjects":   GuageType,
	"HeapReleased":  GuageType,
	"HeapSys":       GuageType,
	"LastGC":        GuageType,
	"Lookups":       GuageType,
	"MCacheInuse":   GuageType,
	"MCacheSys":     GuageType,
	"MSpanInuse":    GuageType,
	"MSpanSys":      GuageType,
	"Alloc":         GuageType,
	"GCCPUFraction": GuageType,
	"NextGC":        GuageType,
	"Mallocs":       GuageType,
	"NumForcedGC":   GuageType,
	"NumGC":         GuageType,
	"OtherSys":      GuageType,
	"PauseTotalNs":  GuageType,
	"StackInuse":    GuageType,
	"StackSys":      GuageType,
	"Sys":           GuageType,
	"TotalAlloc":    GuageType,
	"RandomValue":   GuageType,
}
