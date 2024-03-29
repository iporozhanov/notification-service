// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	notification "notification-service/notification"

	mock "github.com/stretchr/testify/mock"

	query "notification-service/database/query"

	sqlx "github.com/jmoiron/sqlx"
)

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

type DB_Expecter struct {
	mock *mock.Mock
}

func (_m *DB) EXPECT() *DB_Expecter {
	return &DB_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *DB) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DB_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type DB_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *DB_Expecter) Close() *DB_Close_Call {
	return &DB_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *DB_Close_Call) Run(run func()) *DB_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DB_Close_Call) Return(_a0 error) *DB_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DB_Close_Call) RunAndReturn(run func() error) *DB_Close_Call {
	_c.Call.Return(run)
	return _c
}

// GetPendingNotificationsByType provides a mock function with given fields: nType, attempts
func (_m *DB) GetPendingNotificationsByType(nType notification.NotificationType, attempts int64) (*sqlx.Tx, []*query.Notification, error) {
	ret := _m.Called(nType, attempts)

	if len(ret) == 0 {
		panic("no return value specified for GetPendingNotificationsByType")
	}

	var r0 *sqlx.Tx
	var r1 []*query.Notification
	var r2 error
	if rf, ok := ret.Get(0).(func(notification.NotificationType, int64) (*sqlx.Tx, []*query.Notification, error)); ok {
		return rf(nType, attempts)
	}
	if rf, ok := ret.Get(0).(func(notification.NotificationType, int64) *sqlx.Tx); ok {
		r0 = rf(nType, attempts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func(notification.NotificationType, int64) []*query.Notification); ok {
		r1 = rf(nType, attempts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*query.Notification)
		}
	}

	if rf, ok := ret.Get(2).(func(notification.NotificationType, int64) error); ok {
		r2 = rf(nType, attempts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DB_GetPendingNotificationsByType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPendingNotificationsByType'
type DB_GetPendingNotificationsByType_Call struct {
	*mock.Call
}

// GetPendingNotificationsByType is a helper method to define mock.On call
//   - nType notification.NotificationType
//   - attempts int64
func (_e *DB_Expecter) GetPendingNotificationsByType(nType interface{}, attempts interface{}) *DB_GetPendingNotificationsByType_Call {
	return &DB_GetPendingNotificationsByType_Call{Call: _e.mock.On("GetPendingNotificationsByType", nType, attempts)}
}

func (_c *DB_GetPendingNotificationsByType_Call) Run(run func(nType notification.NotificationType, attempts int64)) *DB_GetPendingNotificationsByType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(notification.NotificationType), args[1].(int64))
	})
	return _c
}

func (_c *DB_GetPendingNotificationsByType_Call) Return(_a0 *sqlx.Tx, _a1 []*query.Notification, _a2 error) *DB_GetPendingNotificationsByType_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *DB_GetPendingNotificationsByType_Call) RunAndReturn(run func(notification.NotificationType, int64) (*sqlx.Tx, []*query.Notification, error)) *DB_GetPendingNotificationsByType_Call {
	_c.Call.Return(run)
	return _c
}

// InsertNotification provides a mock function with given fields: _a0
func (_m *DB) InsertNotification(_a0 query.Notification) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for InsertNotification")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(query.Notification) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DB_InsertNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertNotification'
type DB_InsertNotification_Call struct {
	*mock.Call
}

// InsertNotification is a helper method to define mock.On call
//   - _a0 query.Notification
func (_e *DB_Expecter) InsertNotification(_a0 interface{}) *DB_InsertNotification_Call {
	return &DB_InsertNotification_Call{Call: _e.mock.On("InsertNotification", _a0)}
}

func (_c *DB_InsertNotification_Call) Run(run func(_a0 query.Notification)) *DB_InsertNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(query.Notification))
	})
	return _c
}

func (_c *DB_InsertNotification_Call) Return(err error) *DB_InsertNotification_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *DB_InsertNotification_Call) RunAndReturn(run func(query.Notification) error) *DB_InsertNotification_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateNotifications provides a mock function with given fields: tx, notifications
func (_m *DB) UpdateNotifications(tx *sqlx.Tx, notifications []*query.Notification) error {
	ret := _m.Called(tx, notifications)

	if len(ret) == 0 {
		panic("no return value specified for UpdateNotifications")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*sqlx.Tx, []*query.Notification) error); ok {
		r0 = rf(tx, notifications)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DB_UpdateNotifications_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateNotifications'
type DB_UpdateNotifications_Call struct {
	*mock.Call
}

// UpdateNotifications is a helper method to define mock.On call
//   - tx *sqlx.Tx
//   - notifications []*query.Notification
func (_e *DB_Expecter) UpdateNotifications(tx interface{}, notifications interface{}) *DB_UpdateNotifications_Call {
	return &DB_UpdateNotifications_Call{Call: _e.mock.On("UpdateNotifications", tx, notifications)}
}

func (_c *DB_UpdateNotifications_Call) Run(run func(tx *sqlx.Tx, notifications []*query.Notification)) *DB_UpdateNotifications_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*sqlx.Tx), args[1].([]*query.Notification))
	})
	return _c
}

func (_c *DB_UpdateNotifications_Call) Return(_a0 error) *DB_UpdateNotifications_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DB_UpdateNotifications_Call) RunAndReturn(run func(*sqlx.Tx, []*query.Notification) error) *DB_UpdateNotifications_Call {
	_c.Call.Return(run)
	return _c
}

// NewDB creates a new instance of DB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDB(t interface {
	mock.TestingT
	Cleanup(func())
}) *DB {
	mock := &DB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
