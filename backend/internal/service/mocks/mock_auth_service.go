// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/auth.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/auth.go -destination=internal/service/mocks/mock_auth_service.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	core "github.com/Cheasezz/anSpace/backend/internal/core"
	auth "github.com/Cheasezz/anSpace/backend/pkg/auth"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// GetUser mocks base method.
func (m *MockAuth) GetUser(ctx context.Context, userId uuid.UUID) (core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, userId)
	ret0, _ := ret[0].(core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAuthMockRecorder) GetUser(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuth)(nil).GetUser), ctx, userId)
}

// LogOut mocks base method.
func (m *MockAuth) LogOut(ctx context.Context, refreshToken string) (auth.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogOut", ctx, refreshToken)
	ret0, _ := ret[0].(auth.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogOut indicates an expected call of LogOut.
func (mr *MockAuthMockRecorder) LogOut(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogOut", reflect.TypeOf((*MockAuth)(nil).LogOut), ctx, refreshToken)
}

// RefreshAccessToken mocks base method.
func (m *MockAuth) RefreshAccessToken(ctx context.Context, refreshToken string) (auth.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshAccessToken", ctx, refreshToken)
	ret0, _ := ret[0].(auth.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshAccessToken indicates an expected call of RefreshAccessToken.
func (mr *MockAuthMockRecorder) RefreshAccessToken(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshAccessToken", reflect.TypeOf((*MockAuth)(nil).RefreshAccessToken), ctx, refreshToken)
}

// SignIn mocks base method.
func (m *MockAuth) SignIn(ctx context.Context, signIn core.AuthCredentials) (auth.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, signIn)
	ret0, _ := ret[0].(auth.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthMockRecorder) SignIn(ctx, signIn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuth)(nil).SignIn), ctx, signIn)
}

// SignUp mocks base method.
func (m *MockAuth) SignUp(ctx context.Context, signUp core.AuthCredentials) (auth.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, signUp)
	ret0, _ := ret[0].(auth.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthMockRecorder) SignUp(ctx, signUp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuth)(nil).SignUp), ctx, signUp)
}
