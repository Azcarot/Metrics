package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MetricRouter() chi.Router {
	storagehandler := &StorageHandler{
		Storage: &types.MemStorage{
			Gaugemem: make(map[string]types.Gauge), Countermem: make(map[string]types.Counter)},
	}
	r := chi.NewRouter()
	r.Get("/update/{type}/{name}/{value}", storagehandler.HandlePostMetrics)
	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandlePostMetrics(t *testing.T) {
	ts := httptest.NewServer(MetricRouter())
	var testTable = []struct {
		url    string
		want   string
		status int
	}{
		{"/update/counter/testcounter/2", "", http.StatusOK},
		{"/update/gauge/testgauge/44", "", http.StatusOK},
		// проверим на ошибочный запрос
		{"/update/fail/fail/3", "", http.StatusBadRequest},
	}
	for _, v := range testTable {
		resp, get := testRequest(t, ts, "GET", v.url)
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
		defer resp.Body.Close()
	}
}
