// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/psql/auth.go
//
// Generated by this command:
//
//	mockgen -source=internal/repository/psql/auth.go -destination=internal/repository/psql/mocks/mock_auth_repo.go
//

// Package mock_psql is a generated GoMock package.
package mock_psql

import (
	context "context"
	reflect "reflect"

	core "github.com/Cheasezz/anSpace/backend/internal/core"
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

// CreateUser mocks base method.
func (m *MockAuth) CreateUser(ctx context.Context, signUp core.AuthCredentials) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, signUp)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthMockRecorder) CreateUser(ctx, signUp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuth)(nil).CreateUser), ctx, signUp)
}

// DeletePassResetCode mocks base method.
func (m *MockAuth) DeletePassResetCode(ctx context.Context, code core.CodeCredentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePassResetCode", ctx, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePassResetCode indicates an expected call of DeletePassResetCode.
func (mr *MockAuthMockRecorder) DeletePassResetCode(ctx, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePassResetCode", reflect.TypeOf((*MockAuth)(nil).DeletePassResetCode), ctx, code)
}

// GetUserByEmail mocks base method.
func (m *MockAuth) GetUserByEmail(ctx context.Context, email string) (core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockAuthMockRecorder) GetUserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockAuth)(nil).GetUserByEmail), ctx, email)
}

// GetUserById mocks base method.
func (m *MockAuth) GetUserById(ctx context.Context, userId uuid.UUID) (core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", ctx, userId)
	ret0, _ := ret[0].(core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockAuthMockRecorder) GetUserById(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockAuth)(nil).GetUserById), ctx, userId)
}

// GetUserIdByLogPas mocks base method.
func (m *MockAuth) GetUserIdByLogPas(ctx context.Context, signIn core.AuthCredentials) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIdByLogPas", ctx, signIn)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIdByLogPas indicates an expected call of GetUserIdByLogPas.
func (mr *MockAuthMockRecorder) GetUserIdByLogPas(ctx, signIn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIdByLogPas", reflect.TypeOf((*MockAuth)(nil).GetUserIdByLogPas), ctx, signIn)
}

// GetUserSessionByRefreshToken mocks base method.
func (m *MockAuth) GetUserSessionByRefreshToken(ctx context.Context, refreshToken string) (core.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSessionByRefreshToken", ctx, refreshToken)
	ret0, _ := ret[0].(core.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSessionByRefreshToken indicates an expected call of GetUserSessionByRefreshToken.
func (mr *MockAuthMockRecorder) GetUserSessionByRefreshToken(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSessionByRefreshToken", reflect.TypeOf((*MockAuth)(nil).GetUserSessionByRefreshToken), ctx, refreshToken)
}

// SetPassResetCode mocks base method.
func (m *MockAuth) SetPassResetCode(ctx context.Context, code core.CodeCredentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPassResetCode", ctx, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPassResetCode indicates an expected call of SetPassResetCode.
func (mr *MockAuthMockRecorder) SetPassResetCode(ctx, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPassResetCode", reflect.TypeOf((*MockAuth)(nil).SetPassResetCode), ctx, code)
}

// SetSession mocks base method.
func (m *MockAuth) SetSession(ctx context.Context, session core.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSession", ctx, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSession indicates an expected call of SetSession.
func (mr *MockAuthMockRecorder) SetSession(ctx, session any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSession", reflect.TypeOf((*MockAuth)(nil).SetSession), ctx, session)
}
