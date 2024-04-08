package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azcarot/Metrics/internal/storage"
)

// HandleJSONPostMetrics обрабатывает запросы на запись единичной метрики, принятой в виде JSON
// Метрики всегда пишутся во внутренную память, а запись их в файл
// или бд определяется соответствующими флагами
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
		data := buf.Bytes()

		if err = json.Unmarshal(data, &metricData); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		err = st.Storage.StoreMetrics(metricData)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(flag.FlagDBAddr) != 0 {
			storage.PgxStorage.WriteMetricsToPstgrs(storage.ST, metricData)
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
		switch metricData.MType {
		case storage.CounterType:
			var newvalue int64
			newvalue, err = strconv.ParseInt(result, 0, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			metricData.Delta = &newvalue
			metricResult = metricData

		case storage.GuageType:
			var newvalue float64
			newvalue, err = strconv.ParseFloat(result, 64)
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
		if len(flag.FlagKey) > 0 {
			result, _ := st.Storage.GetStoredMetrics(metricData.ID, strings.ToLower(metricData.MType))
			result = storage.ShaMetrics(result, flag.FlagKey)
			res.Header().Set("HashSHA256", result)
		}
		res.Header().Set("Content-Type", storage.JSONContentType)
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
	}

	return http.HandlerFunc(postMetric)
}

// HandleMultipleJSONPostMetrics обрабатывает запросы
// на запись множества метрик, принятых в виде JSON
// Метрики всегда пишутся во внутренную память, а запись их в файл
// или бд определяется соответствующими флагами
func (st *StorageHandler) HandleMultipleJSONPostMetrics(flag storage.Flags) http.Handler {
	getMetrics := func(res http.ResponseWriter, req *http.Request) {
		var metrics []storage.Metrics
		var buf bytes.Buffer
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, metric := range metrics {
			if strings.ToLower(metric.MType) != storage.CounterType && strings.ToLower(metric.MType) != storage.GuageType {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		if flag.FlagDBAddr != "" {
			err := storage.ST.BatchWriteToPstgrs(metrics)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		var storeerr error
		for _, metricData := range metrics {
			storeerr = st.Storage.StoreMetrics(metricData)
			if storeerr != nil {
				break
			}
			if len(flag.FlagFileStorage) != 0 && flag.FlagStoreInterval == 0 && flag.FlagDBAddr == "" {
				fileName := flag.FlagFileStorage
				storage.WriteToFile(fileName, metricData)
			}
		}
		if storeerr != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(getMetrics)
}

// HandleJSONGetMetrics обрабатывает запросы на чтение метрик
func (st *StorageHandler) HandleJSONGetMetrics(flag storage.Flags) http.Handler {

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
					if len(flag.FlagKey) > 0 {
						result, _ = st.Storage.GetStoredMetrics(metric.ID, strings.ToLower(metric.MType))
						result = storage.ShaMetrics(result, flag.FlagKey)
						res.Header().Set("HashSHA256", result)
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
					if len(flag.FlagKey) > 0 {
						result, _ := st.Storage.GetStoredMetrics(metric.ID, strings.ToLower(metric.MType))
						result = storage.ShaMetrics(result, flag.FlagKey)
						res.Header().Set("HashSHA256", result)
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
