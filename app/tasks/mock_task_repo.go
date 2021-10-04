// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/shaheerhas/todo-list/app/tasks (interfaces: TaskRepo)

// Package mocks is a generated GoMock package.
package tasks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTaskRepo is a mock of TaskRepo interface.
type MockTaskRepo struct {
	ctrl     *gomock.Controller
	recorder *MockTaskRepoMockRecorder
}

// MockTaskRepoMockRecorder is the mock recorder for MockTaskRepo.
type MockTaskRepoMockRecorder struct {
	mock *MockTaskRepo
}

// NewMockTaskRepo creates a new mock instance.
func NewMockTaskRepo(ctrl *gomock.Controller) *MockTaskRepo {
	mock := &MockTaskRepo{ctrl: ctrl}
	mock.recorder = &MockTaskRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskRepo) EXPECT() *MockTaskRepoMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTaskRepo) CreateTask(arg0 TaskApp, arg1 Task) (Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", arg0, arg1)
	ret0, _ := ret[0].(Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskRepoMockRecorder) CreateTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskRepo)(nil).CreateTask), arg0, arg1)
}

// UpdateTask mocks base method.
func (m *MockTaskRepo) UpdateTask(arg0 TaskApp, arg1 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskRepoMockRecorder) UpdateTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskRepo)(nil).UpdateTask), arg0, arg1)
}