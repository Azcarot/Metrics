package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azcarot/Metrics/cmd/types"
)

func (st *StorageHandler) HandleJSONPostMetrics() http.Handler {
	var metricData types.Metrics
	var metricResult types.Metrics

	postMetric := func(res http.ResponseWriter, req *http.Request) {
		var buf bytes.Buffer
		res.Header().Set("Content-Type", types.JSONContentType)
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
		case types.CounterType:
			value := strconv.Itoa(int(*metricData.Delta))
			err := st.Storage.StoreMetrics(metricData.ID, strings.ToLower(metricData.MType), value)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
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

		case types.GuageType:
			value := strconv.FormatFloat(float64(*metricData.Value), 'g', -1, 64)
			err := st.Storage.StoreMetrics(metricData.ID, strings.ToLower(metricData.MType), value)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
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
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
	}

	return http.HandlerFunc(postMetric)
}

func (st *StorageHandler) HandleJSONGetMetrics() http.Handler {

	getMetric := func(res http.ResponseWriter, req *http.Request) {
		var metric types.Metrics
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
			result, err := st.Storage.GetStoredMetrics(metric.ID, strings.ToLower(metric.MType))
			res.Header().Add("Content-Type", types.JSONContentType)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			} else {
				fmt.Println("Type ", metric.MType, " Name ", metric.ID)
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
