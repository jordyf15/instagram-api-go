// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	domain "instagram-go/domain"

	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// InsertUser provides a mock function with given fields: _a0
func (_m *UserUsecase) InsertUser(_a0 *domain.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUser provides a mock function with given fields: _a0, _a1, _a2
func (_m *UserUsecase) UpdateUser(_a0 *domain.User, _a1 string, _a2 multipart.File) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User, string, multipart.File) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyCredential provides a mock function with given fields: _a0, _a1
func (_m *UserUsecase) VerifyCredential(_a0 string, _a1 string) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
