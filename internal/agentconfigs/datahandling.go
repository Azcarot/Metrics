package agentconfigs

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Azcarot/Metrics/internal/storage"
)

func GzipForAgent(b []byte) []byte {
	var w bytes.Buffer
	gz, err := gzip.NewWriterLevel(&w, gzip.BestSpeed)
	if err != nil {
		panic(fmt.Sprintf("cannot make gzip writer %s ", b))

	}
	defer gz.Close()
	if _, err := gz.Write(b); err != nil {
		panic(fmt.Sprintf("cannot make gzip %s ", b))
	}
	return b
}

func Makepath(m storage.MemStorage, a string) []string {
	var path []string
	pathscount := 0
	for name, value := range m.Gaugemem {
		path = append(path, "http://"+a+"/update/gauge/"+name+"/"+strconv.FormatFloat(float64(value), 'g', -1, 64))
		pathscount++
	}
	for name, value := range m.Countermem {
		path = append(path, "http://"+a+"/update/counter/"+name+"/"+strconv.Itoa(int(value)))
		pathscount++
	}
	return path
}
func MakeSHA(b []byte, k string) []byte {

	key := []byte(k)
	// создаём новый hash.Hash, вычисляющий контрольную сумму SHA-256
	h := hmac.New(sha256.New, key)
	// передаём байты для хеширования
	h.Write(b)
	// вычисляем хеш
	hash := h.Sum(nil)
	return hash
}

func MakeJSON(m storage.MemStorage) ([][]byte, []byte) {
	var body [][]byte
	var metric storage.Metrics
	var metrics []storage.Metrics
	for name, value := range m.Gaugemem {
		newvalue := float64(value)
		metric.ID = name
		metric.MType = "gauge"
		metric.Value = &newvalue
		resp, err := json.Marshal(metric)
		if err != nil {
			panic(fmt.Sprintf("cannot make json %s ", body))
		}
		body = append(body, resp)
		metrics = append(metrics, metric)
	}
	for name, value := range m.Countermem {
		newvalue := int64(value)
		metric.ID = name
		metric.MType = "counter"
		metric.Delta = &newvalue
		resp, err := json.Marshal(metric)
		if err != nil {
			panic(fmt.Sprintf("cannot make json %s ", body))
		}
		body = append(body, resp)
		metrics = append(metrics, metric)
	}
	fullJSON, err := json.Marshal(metrics)
	if err != nil {
		panic(fmt.Sprintf("cannot make []json %s ", err))
	}
	return body, fullJSON
}
