package handlers

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_storage "github.com/Azcarot/Metrics/internal/mock"
	"github.com/Azcarot/Metrics/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestStorageHandler_HandleJSONPostMetrics(t *testing.T) {
	type requestData struct {
		data storage.Metrics
	}
	type testing struct {
		name         string
		metricData   storage.Metrics
		flag         storage.Flags
		expStatus    int
		wrongRequest bool
	}
	tests := make([]testing, 4)
	Storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}
	delta := rand.Int63n(100)
	gauge := rand.Float64()
	tests[0] = testing{name: "BadRequest(Wrong Metric Data)", expStatus: 400,
		wrongRequest: true, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: storage.Metrics{ID: "222", MType: "Gauge", Delta: &delta},
	}
	tests[1] = testing{name: "Goodrequest(Gauge Data)", expStatus: 200,
		wrongRequest: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: storage.Metrics{ID: "222", MType: storage.GuageType, Value: &gauge},
	}
	tests[2] = testing{name: "GoodRequest(Counter Data)", expStatus: 200,
		wrongRequest: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: storage.Metrics{ID: "222", MType: storage.CounterType, Delta: &delta},
	}
	tests[3] = testing{name: "BadRequest(No Data)", expStatus: 400,
		wrongRequest: true, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: storage.Metrics{ID: "222"},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_storage.NewMockPgxStorage(ctrl)
	for _, tt := range tests {
		var req requestData
		req.data = tt.metricData

		storage.ST = mock
		if !tt.wrongRequest {
			mock.EXPECT().WriteMetricsToPstgrs(gomock.Eq(tt.metricData)).Times(1)
		}
		handler := Storagehandler.HandleJSONPostMetrics(tt.flag)
		recorder := httptest.NewRecorder()
		url := "/update/"
		body, err := json.Marshal(req.data)
		require.NoError(t, err)
		reader := bytes.NewReader(body)
		reqst, err := http.NewRequest(http.MethodPost, url, reader)
		require.NoError(t, err)
		handler.ServeHTTP(recorder, reqst)
		require.Equal(t, tt.expStatus, recorder.Code)
		if tt.expStatus != 400 {
			var respMetrics storage.Metrics
			err = json.Unmarshal(recorder.Body.Bytes(), &respMetrics)
			require.NoError(t, err)
			type metricWithValue struct {
				ID    string  `json:"id"`
				MType string  `json:"type"`
				Delta int64   `json:"delta,omitempty"`
				Value float64 `json:"value,omitempty"`
			}
			var readValue float64
			var respMetricsWithValue metricWithValue
			var metricsWithValue metricWithValue
			if respMetrics.MType == storage.GuageType {
				readValue = *respMetrics.Value
				respMetricsWithValue.ID = respMetrics.ID
				respMetricsWithValue.MType = respMetrics.MType
				respMetricsWithValue.Value = readValue
				metricValue := *tt.metricData.Value
				metricsWithValue.ID = tt.metricData.ID
				metricsWithValue.MType = tt.metricData.MType
				metricsWithValue.Value = metricValue

				require.Equal(t, respMetricsWithValue, metricsWithValue)

			}
		}
	}
}

func TestStorageHandler_HandleMultipleJSONPostMetrics(t *testing.T) {
	type requestData struct {
		data []storage.Metrics
	}
	type testing struct {
		name       string
		metricData []storage.Metrics
		flag       storage.Flags
		expStatus  int
		writeToDB  bool
	}
	tests := make([]testing, 5)
	Storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}
	delta := rand.Int63n(100)
	gauge := rand.Float64()
	tests[0] = testing{name: "BadRequest(Wrong Metric Data)", expStatus: 400,
		writeToDB: true, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: []storage.Metrics{{ID: "222", MType: "Wrong", Delta: &delta}, {ID: "222", MType: "Wrong", Delta: &delta}},
	}
	tests[1] = testing{name: "Goodrequest(Gauge Data)", expStatus: 200,
		writeToDB: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: []storage.Metrics{{ID: "222", MType: storage.GuageType, Value: &gauge}, {ID: "333", MType: storage.GuageType, Value: &gauge}},
	}
	tests[2] = testing{name: "GoodRequest(Counter Data)", expStatus: 200,
		writeToDB: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: []storage.Metrics{{ID: "222", MType: storage.CounterType, Delta: &delta}, {ID: "333", MType: storage.CounterType, Delta: &delta}},
	}
	tests[3] = testing{name: "BadRequest(No Data)", expStatus: 400,
		writeToDB: true, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: []storage.Metrics{{ID: "222"}},
	}
	tests[4] = testing{name: "GoodRequest(Counter Data)", expStatus: 200,
		writeToDB: true, flag: storage.Flags{},
		metricData: []storage.Metrics{{ID: "222", MType: storage.CounterType, Delta: &delta}, {ID: "333", MType: storage.CounterType, Delta: &delta}},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_storage.NewMockPgxStorage(ctrl)
	for _, tt := range tests {
		var req requestData
		req.data = tt.metricData

		storage.ST = mock
		if !tt.writeToDB {
			mock.EXPECT().BatchWriteToPstgrs(gomock.Eq(tt.metricData)).Times(1)
		}
		handler := Storagehandler.HandleMultipleJSONPostMetrics(tt.flag)
		recorder := httptest.NewRecorder()
		url := "/updates/"
		body, err := json.Marshal(req.data)
		require.NoError(t, err)
		reader := bytes.NewReader(body)
		reqst, err := http.NewRequest(http.MethodPost, url, reader)
		require.NoError(t, err)
		handler.ServeHTTP(recorder, reqst)
		require.Equal(t, tt.expStatus, recorder.Code)
		if tt.expStatus != 400 {
			storedData := Storagehandler.Storage.GetAllMetricsAsMetricType()
			type metricWithValue struct {
				ID    string  `json:"id"`
				MType string  `json:"type"`
				Delta int64   `json:"delta,omitempty"`
				Value float64 `json:"value,omitempty"`
			}
			storedDataWithValue := make([]metricWithValue, len(storedData))
			metricDataWithValue := make([]metricWithValue, len(tt.metricData))
			for i, metric := range storedData {
				if metric.MType == storage.GuageType {
					storedDataWithValue[i].Value = *metric.Value
				}
				if metric.MType == storage.CounterType {
					storedDataWithValue[i].Delta = *metric.Delta
				}
				storedDataWithValue[i].ID = metric.ID
				storedDataWithValue[i].MType = metric.MType
			}
			for i, metric := range tt.metricData {
				if metric.MType == storage.GuageType {
					metricDataWithValue[i].Value = *metric.Value

					metricDataWithValue[i].ID = metric.ID
					metricDataWithValue[i].MType = metric.MType
					require.Contains(t, storedDataWithValue, metricDataWithValue[i])
				}
			}
		}
	}
}

func TestStorageHandler_HandleJSONGetMetrics(t *testing.T) {
	type requestData struct {
		data storage.Metrics
	}
	type testing struct {
		name             string
		metricData       storage.Metrics
		flag             storage.Flags
		expSendStatus    int
		expReadStatus    int
		wrongReadRequest bool
		metricReadData   storage.Metrics
	}
	tests := make([]testing, 5)
	Storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}
	delta := rand.Int63n(100)
	gauge := rand.Float64()
	tests[0] = testing{name: "BadRequest(Wrong Metric Data)", expSendStatus: 400,
		wrongReadRequest: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: storage.Metrics{ID: "111", MType: "Gauge", Delta: &delta},
	}
	tests[1] = testing{name: "Goodrequest(Gauge Data)", expSendStatus: 200,
		wrongReadRequest: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData:    storage.Metrics{ID: "222", MType: storage.GuageType, Value: &gauge},
		expReadStatus: 200,
	}
	tests[2] = testing{name: "GoodRequest(Counter Data)", expSendStatus: 200,
		wrongReadRequest: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData:    storage.Metrics{ID: "333", MType: storage.CounterType, Delta: &delta},
		expReadStatus: 200,
	}
	tests[3] = testing{name: "BadRequest(No Data)", expSendStatus: 400,
		wrongReadRequest: false, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData: storage.Metrics{ID: "444"},
	}
	tests[4] = testing{name: "GoodRequest(No saved data)", expSendStatus: 200,
		wrongReadRequest: true, flag: storage.Flags{FlagDBAddr: "someDB"},
		metricData:    storage.Metrics{ID: "333", MType: storage.CounterType, Delta: &delta},
		expReadStatus: 404, metricReadData: storage.Metrics{ID: "888", MType: "Wrong"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_storage.NewMockPgxStorage(ctrl)
	for _, tt := range tests {
		var req requestData
		req.data = tt.metricData

		storage.ST = mock
		if tt.metricData.MType == storage.CounterType || tt.metricData.MType == storage.GuageType {
			mock.EXPECT().WriteMetricsToPstgrs(gomock.Eq(tt.metricData)).Times(1)
		}
		handler := Storagehandler.HandleJSONPostMetrics(tt.flag)
		recorder := httptest.NewRecorder()
		url := "/update/"
		body, err := json.Marshal(req.data)
		require.NoError(t, err)
		reader := bytes.NewReader(body)
		reqst, err := http.NewRequest(http.MethodPost, url, reader)
		require.NoError(t, err)
		handler.ServeHTTP(recorder, reqst)
		require.Equal(t, tt.expSendStatus, recorder.Code)
		if tt.expSendStatus != 400 {

			handler := Storagehandler.HandleJSONGetMetrics(tt.flag)
			recorder := httptest.NewRecorder()
			url := "/value/"
			reqData, err := json.Marshal(req.data)
			if tt.wrongReadRequest {
				req.data = tt.metricReadData
				reqData, err = json.Marshal(req.data)
			}
			require.NoError(t, err)
			reader := bytes.NewReader(reqData)
			reqst, err := http.NewRequest(http.MethodGet, url, reader)
			require.NoError(t, err)
			handler.ServeHTTP(recorder, reqst)
			var readData storage.Metrics
			require.Equal(t, tt.expReadStatus, recorder.Code)
			err = json.Unmarshal(recorder.Body.Bytes(), &readData)
			if !tt.wrongReadRequest {
				require.NoError(t, err)
				require.Equal(t, tt.expReadStatus, recorder.Code)
				type metricWithValue struct {
					ID    string  `json:"id"`
					MType string  `json:"type"`
					Delta int64   `json:"delta,omitempty"`
					Value float64 `json:"value,omitempty"`
				}
				var readValue float64

				if readData.MType == storage.GuageType {
					readValue = *readData.Value

					readDataWithValue := metricWithValue{
						ID:    readData.ID,
						MType: readData.MType,
						Value: readValue,
					}
					var metricValue float64
					var metricDelta int64
					if tt.metricData.MType == storage.GuageType {
						metricValue = *tt.metricData.Value
					}

					if tt.metricData.MType == storage.CounterType {
						metricDelta = *tt.metricData.Delta
					}

					metricDataWithValue := metricWithValue{
						ID:    tt.metricData.ID,
						MType: tt.metricData.MType,
						Value: metricValue,
						Delta: metricDelta,
					}
					require.Equal(t, metricDataWithValue, readDataWithValue)

				}
			}
		}

	}
}
