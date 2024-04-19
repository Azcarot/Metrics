package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"
)

// GuageType Константа определяющая тип gauge
const GuageType = "gauge"

// CounterType Константа определяющая тип counter
const CounterType = "counter"

type Gauge float64
type Counter int64
type MemStorage struct {
	Gaugemem   map[string]Gauge
	Countermem map[string]Counter
}

var storedData MemStorage

// Flags - Все возможные флаги
type Flags struct {
	FlagAddr          string
	FlagStoreInterval int
	FlagFileStorage   string
	FlagRestore       bool
	FlagDBAddr        string
	FlagKey           string
	FlagCrypto        string
}

// ServerENV - Переменные окружения
type ServerENV struct {
	Address       string `env:"ADDRESS"`
	StoreInterval string `env:"STORE_INTERVAL"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
	DBAddress     string `env:"DATABASE_DSN"`
	Key           string `env:"KEY"`
	CryptoKey     string `env:"CRYPTO_KEY"`
}

type AllowedMetrics struct {
	Name string
}

// MemInteractions - Интерфейс для работы с хранилищем
type MemInteractions interface {
	GetStoredMetrics(string, string) (string, error)
	StoreMetrics(data Metrics) error
	GetAllMetrics() string
	ReadMetricsFromFile(string)
	GetAllMetricsAsMetricType() []Metrics
}

// StoreMetrics сохраняет метрику во внутренней памяти
func (m *MemStorage) StoreMetrics(data Metrics) error {

	switch data.MType {
	case GuageType:
		if storedData.Gaugemem == nil {
			storedData.Gaugemem = make(map[string]Gauge)
		}

		m.Gaugemem[data.ID] = Gauge(*data.Value)
		storedData.Gaugemem[data.ID] = m.Gaugemem[data.ID]
	case CounterType:
		if storedData.Countermem == nil {
			storedData.Countermem = make(map[string]Counter)
		}
		delta := Counter(*data.Delta)
		m.Countermem[data.ID] += delta
		storedData.Countermem[data.ID] += delta
	default:
		return errors.New("wrong type")
	}

	return nil
}

// ReadMetricsFromFile читает метрики из файла во внутреннюю память,
// принимает строку с именем файла с метриками
func (m *MemStorage) ReadMetricsFromFile(filename string) {
	storedData.Gaugemem = make(map[string]Gauge)
	storedData.Countermem = make(map[string]Counter)
	Consumer, err := NewConsumer(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer Consumer.Close()
	metrics, err := Consumer.ReadEvent()
	if err != nil {
		log.Fatal(err)
	}

	for _, metric := range *metrics {
		if len(metric.MType) > 0 {
			switch strings.ToLower(metric.MType) {
			case "gauge":
				m.Gaugemem[metric.ID] = Gauge(*metric.Value)
				storedData.Gaugemem[metric.ID] = Gauge(*metric.Value)
			case "counter":
				m.Countermem[metric.ID] += Counter(*metric.Delta)
				storedData.Countermem[metric.ID] += Counter(*metric.Delta)
			}

		}
	}

}

func WriteToFile(f string, mdata Metrics) {
	Producer, err := NewProducer(f)
	if err != nil {
		log.Fatal(err)
	}
	defer Producer.Close()
	if err := Producer.WriteEvent(&mdata); err != nil {
		log.Fatal(err)
	}
}

// GetAllMetricsAsMetricType читает все метрики из внутренней памяти,
// функция выдает все сохраненные метрики в виде слайса типа Metrics
func (m *MemStorage) GetAllMetricsAsMetricType() []Metrics {
	var FinalData []Metrics
	for n, v := range m.Gaugemem {
		var data Metrics
		data.ID = n
		value := float64(v)
		data.Value = &value
		data.MType = "gauge"
		FinalData = append(FinalData, data)
	}
	for n, v := range m.Countermem {
		var data Metrics
		data.ID = n
		value := int64(v)
		data.Delta = &value
		data.MType = "counter"
		FinalData = append(FinalData, data)
	}
	return FinalData
}

// GetAllMetrics читает все метрики из внутренней памяти,
// функция возвращает строку с данными метрик
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

// GetStoredMetrics производит чтение конкретной метрики из внутренней памяти,
// метрика ищется по имени и типу,
// если метрика найдена, возвращается ее значение в виду строки
func (m *MemStorage) GetStoredMetrics(n string, t string) (string, error) {
	var result string
	var err error
	switch t {
	case GuageType:
		val, ok := storedData.Gaugemem[n]

		if ok {
			result = strconv.FormatFloat(float64(val), 'g', -1, 64)
		} else {
			err = errors.New("no metric stored")
		}

	case CounterType:
		val, ok := storedData.Countermem[n]
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

// ShaMetrics кодирует метрики в sh256 по переданному ключу
func ShaMetrics(result string, key string) string {
	b := []byte(result)
	shakey := []byte(key)
	// создаём новый hash.Hash, вычисляющий контрольную сумму SHA-256
	h := hmac.New(sha256.New, shakey)
	// передаём байты для хеширования
	h.Write(b)
	// вычисляем хеш
	hash := h.Sum(nil)
	sha := base64.URLEncoding.EncodeToString(hash)
	return string(sha)
}
