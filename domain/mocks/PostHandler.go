// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// PostHandler is an autogenerated mock type for the PostHandler type
type PostHandler struct {
	mock.Mock
}

// Post provides a mock function with given fields: _a0, _a1
func (_m *PostHandler) Post(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// Posts provides a mock function with given fields: _a0, _a1
func (_m *PostHandler) Posts(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}
