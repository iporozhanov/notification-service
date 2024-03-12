// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// RateLimiter is an autogenerated mock type for the RateLimiter type
type RateLimiter struct {
	mock.Mock
}

type RateLimiter_Expecter struct {
	mock *mock.Mock
}

func (_m *RateLimiter) EXPECT() *RateLimiter_Expecter {
	return &RateLimiter_Expecter{mock: &_m.Mock}
}

// Allow provides a mock function with given fields: key
func (_m *RateLimiter) Allow(key string) bool {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Allow")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RateLimiter_Allow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Allow'
type RateLimiter_Allow_Call struct {
	*mock.Call
}

// Allow is a helper method to define mock.On call
//   - key string
func (_e *RateLimiter_Expecter) Allow(key interface{}) *RateLimiter_Allow_Call {
	return &RateLimiter_Allow_Call{Call: _e.mock.On("Allow", key)}
}

func (_c *RateLimiter_Allow_Call) Run(run func(key string)) *RateLimiter_Allow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *RateLimiter_Allow_Call) Return(_a0 bool) *RateLimiter_Allow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RateLimiter_Allow_Call) RunAndReturn(run func(string) bool) *RateLimiter_Allow_Call {
	_c.Call.Return(run)
	return _c
}

// NewRateLimiter creates a new instance of RateLimiter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRateLimiter(t interface {
	mock.TestingT
	Cleanup(func())
}) *RateLimiter {
	mock := &RateLimiter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}