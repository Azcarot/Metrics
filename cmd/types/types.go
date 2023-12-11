package types

import (
	"errors"
	"strconv"
)

const GuageType = "gauge"
const CounterType = "counter"

type Gauge float64
type Counter int64
type MemStorage struct {
	Gaugemem   map[string]Gauge
	Countermem map[string]Counter
}

type AllowedMetrics struct {
	Name string
}
type MemInteractions interface {
	GetStoredMetrics(string, string) (string, error)
	StoreMetrics(string, string, string) error
	GetAllMetrics() string
}

func (m *MemStorage) StoreMetrics(n string, t string, v string) error {
	switch t {
	case GuageType:
		value, err := strconv.ParseFloat(v, 64)
		m.Gaugemem[n] = Gauge(value)
		return err
	case CounterType:
		value, err := strconv.Atoi(v)
		m.Countermem[n] += Counter(value)
		return err
	}
	return nil
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
	"PollCount":    CounterType,
	"Frees":        GuageType,
	"GCSys":        GuageType,
	"HeapAlloc":    GuageType,
	"HeapIdle":     GuageType,
	"HeapInuse":    GuageType,
	"HeapObjects":  GuageType,
	"HeapReleased": GuageType,
	"HeapSys":      GuageType,
	"LastGC":       GuageType,
	"Lookups":      GuageType,
	"MCacheInuse":  GuageType,
	"MCacheSys":    GuageType,
	"NextGC":       GuageType,
	"Mallocs":      GuageType,
	"NumForcedGC":  GuageType,
	"NumGC":        GuageType,
	"OtherSys":     GuageType,
	"PauseTotalNs": GuageType,
	"StackInuse":   GuageType,
	"StackSys":     GuageType,
	"Sys":          GuageType,
	"TotalAlloc":   GuageType,
	"RandomValue":  GuageType,
}
