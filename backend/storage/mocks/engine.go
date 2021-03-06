// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	storage "github.com/deface90/def-feelings/storage"
	mock "github.com/stretchr/testify/mock"
)

// Engine is an autogenerated mock type for the Engine type
type Engine struct {
	mock.Mock
}

// CheckUserSession provides a mock function with given fields: ctx, token
func (_m *Engine) CheckUserSession(ctx context.Context, token string) (bool, *storage.User) {
	ret := _m.Called(ctx, token)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 *storage.User
	if rf, ok := ret.Get(1).(func(context.Context, string) *storage.User); ok {
		r1 = rf(ctx, token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*storage.User)
		}
	}

	return r0, r1
}

// CreateStatus provides a mock function with given fields: ctx, status
func (_m *Engine) CreateStatus(ctx context.Context, status *storage.Status) (int64, error) {
	ret := _m.Called(ctx, status)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *storage.Status) int64); ok {
		r0 = rf(ctx, status)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *storage.Status) error); ok {
		r1 = rf(ctx, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *Engine) CreateUser(ctx context.Context, user *storage.User) (int64, error) {
	ret := _m.Called(ctx, user)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *storage.User) int64); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *storage.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, userID
func (_m *Engine) DeleteUser(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EditUser provides a mock function with given fields: ctx, user
func (_m *Engine) EditUser(ctx context.Context, user *storage.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *storage.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetStatus provides a mock function with given fields: ctx, statusID
func (_m *Engine) GetStatus(ctx context.Context, statusID int64) (*storage.Status, error) {
	ret := _m.Called(ctx, statusID)

	var r0 *storage.Status
	if rf, ok := ret.Get(0).(func(context.Context, int64) *storage.Status); ok {
		r0 = rf(ctx, statusID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.Status)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, statusID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, userID
func (_m *Engine) GetUser(ctx context.Context, userID int64) (*storage.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *storage.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) *storage.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByToken provides a mock function with given fields: ctx, token
func (_m *Engine) GetUserByToken(ctx context.Context, token string) (*storage.User, error) {
	ret := _m.Called(ctx, token)

	var r0 *storage.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *storage.User); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUsername provides a mock function with given fields: ctx, username
func (_m *Engine) GetUserByUsername(ctx context.Context, username string) (*storage.User, error) {
	ret := _m.Called(ctx, username)

	var r0 *storage.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *storage.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFeelings provides a mock function with given fields: ctx, request
func (_m *Engine) ListFeelings(ctx context.Context, request storage.ListFeelingsRequest) ([]*storage.Feeling, error) {
	ret := _m.Called(ctx, request)

	var r0 []*storage.Feeling
	if rf, ok := ret.Get(0).(func(context.Context, storage.ListFeelingsRequest) []*storage.Feeling); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*storage.Feeling)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, storage.ListFeelingsRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListStatuses provides a mock function with given fields: ctx, request
func (_m *Engine) ListStatuses(ctx context.Context, request storage.ListStatusesRequest) ([]*storage.Status, int64, error) {
	ret := _m.Called(ctx, request)

	var r0 []*storage.Status
	if rf, ok := ret.Get(0).(func(context.Context, storage.ListStatusesRequest) []*storage.Status); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*storage.Status)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, storage.ListStatusesRequest) int64); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, storage.ListStatusesRequest) error); ok {
		r2 = rf(ctx, request)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ListUsers provides a mock function with given fields: ctx, request
func (_m *Engine) ListUsers(ctx context.Context, request storage.ListUsersRequest) ([]*storage.User, int64, error) {
	ret := _m.Called(ctx, request)

	var r0 []*storage.User
	if rf, ok := ret.Get(0).(func(context.Context, storage.ListUsersRequest) []*storage.User); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*storage.User)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, storage.ListUsersRequest) int64); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, storage.ListUsersRequest) error); ok {
		r2 = rf(ctx, request)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Shutdown provides a mock function with given fields:
func (_m *Engine) Shutdown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
