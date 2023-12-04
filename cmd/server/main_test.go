package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostReq(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name      string
		request   *http.Request
		responder httptest.ResponseRecorder
		want      want
	}{
		{name: "OkTest",
			request:   httptest.NewRequest(http.MethodPost, "/update/counter/someMetric/527", nil),
			responder: *httptest.NewRecorder(),
			want: want{
				code:     200,
				response: "200 OK",
			},
		},
		{name: "NotOkTest",
			request:   httptest.NewRequest(http.MethodPost, "/update/sdfsdf/someMetric/527", nil),
			responder: *httptest.NewRecorder(),
			want: want{
				code:     400,
				response: "400 Bad Request",
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostReq(&tt.responder, tt.request)
			res := tt.responder.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, res.Status)
			defer res.Body.Close()
		})
	}
}
