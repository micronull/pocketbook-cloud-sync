// Code generated by MockGen. DO NOT EDIT.
// Source: command.go
//
// Generated by this command:
//
//	mockgen -source command.go -destination mocks/command.go -package mocks -mock_names command=Command
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Command is a mock of command interface.
type Command struct {
	ctrl     *gomock.Controller
	recorder *CommandMockRecorder
	isgomock struct{}
}

// CommandMockRecorder is the mock recorder for Command.
type CommandMockRecorder struct {
	mock *Command
}

// NewCommand creates a new mock instance.
func NewCommand(ctrl *gomock.Controller) *Command {
	mock := &Command{ctrl: ctrl}
	mock.recorder = &CommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Command) EXPECT() *CommandMockRecorder {
	return m.recorder
}

// Description mocks base method.
func (m *Command) Description() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Description")
	ret0, _ := ret[0].(string)
	return ret0
}

// Description indicates an expected call of Description.
func (mr *CommandMockRecorder) Description() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Description", reflect.TypeOf((*Command)(nil).Description))
}

// Help mocks base method.
func (m *Command) Help() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Help")
	ret0, _ := ret[0].(string)
	return ret0
}

// Help indicates an expected call of Help.
func (mr *CommandMockRecorder) Help() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Help", reflect.TypeOf((*Command)(nil).Help))
}

// Run mocks base method.
func (m *Command) Run(args []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", args)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *CommandMockRecorder) Run(args any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*Command)(nil).Run), args)
}