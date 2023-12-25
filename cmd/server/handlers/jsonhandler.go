package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"bytes"
	"encoding/json"

	"github.com/Azcarot/Metrics/cmd/types"
)

func (st *StorageHandler) HandleJSONPostMetrics() http.Handler {
	var metric types.Metrics

	postMetric := func(res http.ResponseWriter, req *http.Request) {
		var buf bytes.Buffer
		jsonDecoder := json.NewDecoder(&buf)
		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = jsonDecoder.Decode(&metric); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(metric.ID) == 0 || len(metric.MType) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case types.CounterType:
			value := strconv.Itoa(int(*metric.Delta))
			err := st.Storage.StoreMetrics(metric.ID, strings.ToLower(metric.MType), value)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Header().Set("Content-Type", types.JSONContentType)
			resp, err := json.Marshal(metric)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Write(resp)
			res.WriteHeader(http.StatusOK)
		case types.GuageType:
			value := strconv.FormatFloat(float64(*metric.Value), 'g', -1, 64)
			err := st.Storage.StoreMetrics(metric.ID, strings.ToLower(metric.MType), value)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Header().Set("Content-Type", types.JSONContentType)
			resp, err := json.Marshal(metric)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.Write(resp)
			res.WriteHeader(http.StatusOK)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	return http.HandlerFunc(postMetric)
}

func (st *StorageHandler) HandleJSONGetMetrics() http.Handler {
	var metric types.Metrics
	getMetric := func(res http.ResponseWriter, req *http.Request) {

		var buf bytes.Buffer
		jsonDecoder := json.NewDecoder(&buf)
		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = jsonDecoder.Decode(&metric); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := st.Storage.GetStoredMetrics(metric.MType, metric.ID)
		res.Header().Add("Content-Type", types.JSONContentType)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			switch metric.MType {
			case types.CounterType:
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
			case types.GuageType:
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
			}
		}
	}
	return http.HandlerFunc(getMetric)
}
