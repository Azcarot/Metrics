package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azcarot/Metrics/cmd/types"
)

func (st *StorageHandler) HandleJSONPostMetrics() http.Handler {
	var metric types.Metrics

	postMetric := func(res http.ResponseWriter, req *http.Request) {
		var buf bytes.Buffer
		res.Header().Set("Content-Type", types.JSONContentType)
		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
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
			result, err := st.Storage.GetStoredMetrics(metric.ID, strings.ToLower(metric.MType))

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			newvalue, err := strconv.ParseInt(result, 0, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			metric.Delta = &newvalue
			resp, err := json.Marshal(metric)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.WriteHeader(http.StatusOK)
			res.Write(resp)
		case types.GuageType:
			value := strconv.FormatFloat(float64(*metric.Value), 'g', -1, 64)
			err := st.Storage.StoreMetrics(metric.ID, strings.ToLower(metric.MType), value)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			result, err := st.Storage.GetStoredMetrics(metric.ID, strings.ToLower(metric.MType))

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			newvalue, err := strconv.ParseFloat(result, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			metric.Value = &newvalue
			resp, err := json.Marshal(metric)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			res.WriteHeader(http.StatusOK)
			res.Write(resp)
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
		// читаем тело запроса
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
			result, err := st.Storage.GetStoredMetrics(metric.MType, metric.ID)
			res.Header().Add("Content-Type", types.JSONContentType)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			} else {
				println("2q4qrwetewt")
				println(metric.MType)
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
				default:
					res.WriteHeader(http.StatusNotFound)
					return
				}
			}
		} else {
			res.Header().Add("Content-Type", types.JSONContentType)
		}
	}
	return http.HandlerFunc(getMetric)
}
