// Code generated by MockGen. DO NOT EDIT.
// Source: lexer.go

// Package mock_lexer is a generated GoMock package.
package mock_lexer

import (
	gomock "github.com/golang/mock/gomock"
	token "github.com/xujinzheng/monkey/token"
	reflect "reflect"
)

// MockLexer is a mock of Lexer interface
type MockLexer struct {
	ctrl     *gomock.Controller
	recorder *MockLexerMockRecorder
}

// MockLexerMockRecorder is the mock recorder for MockLexer
type MockLexerMockRecorder struct {
	mock *MockLexer
}

// NewMockLexer creates a new mock instance
func NewMockLexer(ctrl *gomock.Controller) *MockLexer {
	mock := &MockLexer{ctrl: ctrl}
	mock.recorder = &MockLexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLexer) EXPECT() *MockLexerMockRecorder {
	return m.recorder
}

// NextToken mocks base method
func (m *MockLexer) NextToken() token.Token {
	ret := m.ctrl.Call(m, "NextToken")
	ret0, _ := ret[0].(token.Token)
	return ret0
}

// NextToken indicates an expected call of NextToken
func (mr *MockLexerMockRecorder) NextToken() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextToken", reflect.TypeOf((*MockLexer)(nil).NextToken))
}
