// Code generated by MockGen. DO NOT EDIT.
// Source: redis_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/dinorain/kalobranded/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockOrderRedisRepository is a mock of OrderRedisRepository interface.
type MockOrderRedisRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRedisRepositoryMockRecorder
}

// MockOrderRedisRepositoryMockRecorder is the mock recorder for MockOrderRedisRepository.
type MockOrderRedisRepositoryMockRecorder struct {
	mock *MockOrderRedisRepository
}

// NewMockOrderRedisRepository creates a new mock instance.
func NewMockOrderRedisRepository(ctrl *gomock.Controller) *MockOrderRedisRepository {
	mock := &MockOrderRedisRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRedisRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRedisRepository) EXPECT() *MockOrderRedisRepositoryMockRecorder {
	return m.recorder
}

// DeleteOrderCtx mocks base method.
func (m *MockOrderRedisRepository) DeleteOrderCtx(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrderCtx", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrderCtx indicates an expected call of DeleteOrderCtx.
func (mr *MockOrderRedisRepositoryMockRecorder) DeleteOrderCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrderCtx", reflect.TypeOf((*MockOrderRedisRepository)(nil).DeleteOrderCtx), ctx, key)
}

// GetByIdCtx mocks base method.
func (m *MockOrderRedisRepository) GetByIdCtx(ctx context.Context, key string) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIdCtx", ctx, key)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIdCtx indicates an expected call of GetByIdCtx.
func (mr *MockOrderRedisRepositoryMockRecorder) GetByIdCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIdCtx", reflect.TypeOf((*MockOrderRedisRepository)(nil).GetByIdCtx), ctx, key)
}

// SetOrderCtx mocks base method.
func (m *MockOrderRedisRepository) SetOrderCtx(ctx context.Context, key string, seconds int, user *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOrderCtx", ctx, key, seconds, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOrderCtx indicates an expected call of SetOrderCtx.
func (mr *MockOrderRedisRepositoryMockRecorder) SetOrderCtx(ctx, key, seconds, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOrderCtx", reflect.TypeOf((*MockOrderRedisRepository)(nil).SetOrderCtx), ctx, key, seconds, user)
}
