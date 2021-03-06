// Code generated by MockGen. DO NOT EDIT.
// Source: httprequest.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHttpRequestor is a mock of HttpRequestor interface
type MockHttpRequestor struct {
	ctrl     *gomock.Controller
	recorder *MockHttpRequestorMockRecorder
}

// MockHttpRequestorMockRecorder is the mock recorder for MockHttpRequestor
type MockHttpRequestorMockRecorder struct {
	mock *MockHttpRequestor
}

// NewMockHttpRequestor creates a new mock instance
func NewMockHttpRequestor(ctrl *gomock.Controller) *MockHttpRequestor {
	mock := &MockHttpRequestor{ctrl: ctrl}
	mock.recorder = &MockHttpRequestorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpRequestor) EXPECT() *MockHttpRequestorMockRecorder {
	return m.recorder
}

// SendRequest mocks base method
func (m *MockHttpRequestor) SendRequest(url, action string, payload []byte) (int, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRequest", url, action, payload)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SendRequest indicates an expected call of SendRequest
func (mr *MockHttpRequestorMockRecorder) SendRequest(url, action, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRequest", reflect.TypeOf((*MockHttpRequestor)(nil).SendRequest), url, action, payload)
}
