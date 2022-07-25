// Code generated by MockGen. DO NOT EDIT.
// Source: redis_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/dinorain/checkoutaja/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockSessRepository is a mock of SessRepository interface.
type MockSessRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessRepositoryMockRecorder
}

// MockSessRepositoryMockRecorder is the mock recorder for MockSessRepository.
type MockSessRepositoryMockRecorder struct {
	mock *MockSessRepository
}

// NewMockSessRepository creates a new mock instance.
func NewMockSessRepository(ctrl *gomock.Controller) *MockSessRepository {
	mock := &MockSessRepository{ctrl: ctrl}
	mock.recorder = &MockSessRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessRepository) EXPECT() *MockSessRepositoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessRepository) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, session, expire)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessRepositoryMockRecorder) CreateSession(ctx, session, expire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessRepository)(nil).CreateSession), ctx, session, expire)
}

// DeleteById mocks base method.
func (m *MockSessRepository) DeleteById(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockSessRepositoryMockRecorder) DeleteById(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockSessRepository)(nil).DeleteById), ctx, sessionID)
}

// GetSessionById mocks base method.
func (m *MockSessRepository) GetSessionById(ctx context.Context, sessionID string) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionById", ctx, sessionID)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionById indicates an expected call of GetSessionById.
func (mr *MockSessRepositoryMockRecorder) GetSessionById(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionById", reflect.TypeOf((*MockSessRepository)(nil).GetSessionById), ctx, sessionID)
}