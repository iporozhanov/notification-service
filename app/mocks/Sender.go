// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	notification "notification-service/notification"

	mock "github.com/stretchr/testify/mock"
)

// Sender is an autogenerated mock type for the Sender type
type Sender struct {
	mock.Mock
}

type Sender_Expecter struct {
	mock *mock.Mock
}

func (_m *Sender) EXPECT() *Sender_Expecter {
	return &Sender_Expecter{mock: &_m.Mock}
}

// SendNotification provides a mock function with given fields: _a0
func (_m *Sender) SendNotification(_a0 notification.Notification) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendNotification")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(notification.Notification) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Sender_SendNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendNotification'
type Sender_SendNotification_Call struct {
	*mock.Call
}

// SendNotification is a helper method to define mock.On call
//   - _a0 notification.Notification
func (_e *Sender_Expecter) SendNotification(_a0 interface{}) *Sender_SendNotification_Call {
	return &Sender_SendNotification_Call{Call: _e.mock.On("SendNotification", _a0)}
}

func (_c *Sender_SendNotification_Call) Run(run func(_a0 notification.Notification)) *Sender_SendNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(notification.Notification))
	})
	return _c
}

func (_c *Sender_SendNotification_Call) Return(_a0 error) *Sender_SendNotification_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Sender_SendNotification_Call) RunAndReturn(run func(notification.Notification) error) *Sender_SendNotification_Call {
	_c.Call.Return(run)
	return _c
}

// NewSender creates a new instance of Sender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *Sender {
	mock := &Sender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
