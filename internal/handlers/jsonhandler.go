package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azcarot/Metrics/internal/storage"
)

func (st *StorageHandler) HandleJSONPostMetrics(flag storage.Flags) http.Handler {
	var metricData storage.Metrics
	var metricResult storage.Metrics

	postMetric := func(res http.ResponseWriter, req *http.Request) {
		var buf bytes.Buffer
		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metricData); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		switch metricData.MType {
		case storage.CounterType:
			value := strconv.Itoa(int(*metricData.Delta))
			err := st.Storage.StoreMetrics(metricData.ID, strings.ToLower(metricData.MType), value)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(flag.FlagFileStorage) != 0 && flag.FlagStoreInterval == 0 {
				fileName := flag.FlagFileStorage
				storage.WriteToFile(fileName, metricData)
			}
			result, err := st.Storage.GetStoredMetrics(metricData.ID, strings.ToLower(metricData.MType))

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			newvalue, err := strconv.ParseInt(result, 0, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			metricData.Delta = &newvalue
			metricResult = metricData

		case storage.GuageType:
			value := strconv.FormatFloat(float64(*metricData.Value), 'g', -1, 64)

			err := st.Storage.StoreMetrics(metricData.ID, strings.ToLower(metricData.MType), value)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(flag.FlagFileStorage) != 0 && flag.FlagStoreInterval == 0 {
				fileName := flag.FlagFileStorage
				storage.WriteToFile(fileName, metricData)
			}
			result, err := st.Storage.GetStoredMetrics(metricData.ID, strings.ToLower(metricData.MType))

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			newvalue, err := strconv.ParseFloat(result, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			metricData.Value = &newvalue
			metricResult = metricData

		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(metricResult)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", storage.JSONContentType)
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
	}

	return http.HandlerFunc(postMetric)
}

func (st *StorageHandler) HandleJSONGetMetrics() http.Handler {

	getMetric := func(res http.ResponseWriter, req *http.Request) {
		var metric storage.Metrics
		var buf bytes.Buffer
		// переменная reader будет равна r.Body или *gzip.Reader

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(metric.MType) > 0 && len(metric.ID) > 0 {
			result, err := st.Storage.GetStoredMetrics(metric.ID, strings.ToLower(metric.MType))
			res.Header().Add("Content-Type", storage.JSONContentType)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			} else {
				switch metric.MType {
				case storage.CounterType:
					value, err := strconv.ParseInt(result, 0, 64)
					if err != nil {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
					metric.Delta = &value
					resp, err := json.Marshal(metric)
					if err != nil {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
					res.Write(resp)
				case storage.GuageType:
					value, err := strconv.ParseFloat(result, 64)
					if err != nil {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
					metric.Value = &value
					resp, err := json.Marshal(metric)
					if err != nil {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
					res.Write(resp)
				default:
					res.WriteHeader(http.StatusNotFound)
					return
				}
			}
		} else {
			res.Header().Add("Content-Type", storage.JSONContentType)
		}
	}
	return http.HandlerFunc(getMetric)
}
