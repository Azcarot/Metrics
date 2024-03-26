// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azcarot/Metrics/internal/storage (interfaces: PgxStorage)

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	http "net/http"
	reflect "reflect"

	storage "github.com/Azcarot/Metrics/internal/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockPgxStorage is a mock of PgxStorage interface.
type MockPgxStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPgxStorageMockRecorder
}

// MockPgxStorageMockRecorder is the mock recorder for MockPgxStorage.
type MockPgxStorageMockRecorder struct {
	mock *MockPgxStorage
}

// NewMockPgxStorage creates a new mock instance.
func NewMockPgxStorage(ctrl *gomock.Controller) *MockPgxStorage {
	mock := &MockPgxStorage{ctrl: ctrl}
	mock.recorder = &MockPgxStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPgxStorage) EXPECT() *MockPgxStorageMockRecorder {
	return m.recorder
}

// BatchWriteToPstgrs mocks base method.
func (m *MockPgxStorage) BatchWriteToPstgrs(arg0 []storage.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchWriteToPstgrs", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchWriteToPstgrs indicates an expected call of BatchWriteToPstgrs.
func (mr *MockPgxStorageMockRecorder) BatchWriteToPstgrs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchWriteToPstgrs", reflect.TypeOf((*MockPgxStorage)(nil).BatchWriteToPstgrs), arg0)
}

// CheckDBConnection mocks base method.
func (m *MockPgxStorage) CheckDBConnection() http.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDBConnection")
	ret0, _ := ret[0].(http.Handler)
	return ret0
}

// CheckDBConnection indicates an expected call of CheckDBConnection.
func (mr *MockPgxStorageMockRecorder) CheckDBConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDBConnection", reflect.TypeOf((*MockPgxStorage)(nil).CheckDBConnection))
}

// CreateTablesForMetrics mocks base method.
func (m *MockPgxStorage) CreateTablesForMetrics() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateTablesForMetrics")
}

// CreateTablesForMetrics indicates an expected call of CreateTablesForMetrics.
func (mr *MockPgxStorageMockRecorder) CreateTablesForMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTablesForMetrics", reflect.TypeOf((*MockPgxStorage)(nil).CreateTablesForMetrics))
}

// WriteMetricsToPstgrs mocks base method.
func (m *MockPgxStorage) WriteMetricsToPstgrs(arg0 storage.Metrics, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteMetricsToPstgrs", arg0, arg1)
}

// WriteMetricsToPstgrs indicates an expected call of WriteMetricsToPstgrs.
func (mr *MockPgxStorageMockRecorder) WriteMetricsToPstgrs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMetricsToPstgrs", reflect.TypeOf((*MockPgxStorage)(nil).WriteMetricsToPstgrs), arg0, arg1)
}
