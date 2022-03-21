// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// LikeHandler is an autogenerated mock type for the LikeHandler type
type LikeHandler struct {
	mock.Mock
}

// DeleteCommentLike provides a mock function with given fields: _a0, _a1
func (_m *LikeHandler) DeleteCommentLike(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// DeleteLikePost provides a mock function with given fields: _a0, _a1
func (_m *LikeHandler) DeleteLikePost(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// PostCommentLike provides a mock function with given fields: _a0, _a1
func (_m *LikeHandler) PostCommentLike(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// PostLikePost provides a mock function with given fields: _a0, _a1
func (_m *LikeHandler) PostLikePost(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}
