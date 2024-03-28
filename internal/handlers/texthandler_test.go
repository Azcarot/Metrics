package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/Azcarot/Metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	r := chi.NewRouter()
	Storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}

	r.Route("/", func(r chi.Router) {
		r.Post("/update/{type}/{name}/{value}", Storagehandler.HandlePostMetrics().ServeHTTP)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	var testTable = []struct {
		url    string
		want   string
		status int
	}{
		{"/update/counter/testcounter/2", "", http.StatusOK},
		{"/update/gauge/testgauge/44", "", http.StatusOK},
		// проверим на ошибочный запрос
		{"/update/fail/fail/3", "", http.StatusBadRequest},
		{"/update//testgauge/3", "", http.StatusBadRequest},
		{"/update/gauge//3", "", http.StatusBadRequest},
		{"/update/gauge/testgauge/data", "", http.StatusBadRequest},
	}
	for _, v := range testTable {
		resp, get := testRequest(t, ts, "POST", v.url)
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
		defer resp.Body.Close()
	}
}
