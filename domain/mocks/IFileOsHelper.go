// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	image "image"
	fs "io/fs"

	io "io"

	mock "github.com/stretchr/testify/mock"

	os "os"
)

// IFileOsHelper is an autogenerated mock type for the IFileOsHelper type
type IFileOsHelper struct {
	mock.Mock
}

// Copy provides a mock function with given fields: _a0, _a1
func (_m *IFileOsHelper) Copy(_a0 io.Writer, _a1 io.Reader) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(io.Writer, io.Reader) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Writer, io.Reader) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: _a0
func (_m *IFileOsHelper) Create(_a0 string) (*os.File, error) {
	ret := _m.Called(_a0)

	var r0 *os.File
	if rf, ok := ret.Get(0).(func(string) *os.File); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*os.File)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecodeImage provides a mock function with given fields: _a0
func (_m *IFileOsHelper) DecodeImage(_a0 io.Reader) (image.Image, string, error) {
	ret := _m.Called(_a0)

	var r0 image.Image
	if rf, ok := ret.Get(0).(func(io.Reader) image.Image); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(image.Image)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(io.Reader) string); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(io.Reader) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MkDirAll provides a mock function with given fields: _a0, _a1
func (_m *IFileOsHelper) MkDirAll(_a0 string, _a1 fs.FileMode) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, fs.FileMode) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResizeAndSaveFileToLocale provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *IFileOsHelper) ResizeAndSaveFileToLocale(_a0 string, _a1 image.Image, _a2 string, _a3 string) (string, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, image.Image, string, string) string); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, image.Image, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}