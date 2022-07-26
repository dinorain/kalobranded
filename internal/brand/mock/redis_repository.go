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

// MockBrandRedisRepository is a mock of BrandRedisRepository interface.
type MockBrandRedisRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBrandRedisRepositoryMockRecorder
}

// MockBrandRedisRepositoryMockRecorder is the mock recorder for MockBrandRedisRepository.
type MockBrandRedisRepositoryMockRecorder struct {
	mock *MockBrandRedisRepository
}

// NewMockBrandRedisRepository creates a new mock instance.
func NewMockBrandRedisRepository(ctrl *gomock.Controller) *MockBrandRedisRepository {
	mock := &MockBrandRedisRepository{ctrl: ctrl}
	mock.recorder = &MockBrandRedisRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBrandRedisRepository) EXPECT() *MockBrandRedisRepositoryMockRecorder {
	return m.recorder
}

// DeleteBrandCtx mocks base method.
func (m *MockBrandRedisRepository) DeleteBrandCtx(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBrandCtx", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBrandCtx indicates an expected call of DeleteBrandCtx.
func (mr *MockBrandRedisRepositoryMockRecorder) DeleteBrandCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBrandCtx", reflect.TypeOf((*MockBrandRedisRepository)(nil).DeleteBrandCtx), ctx, key)
}

// GetByIdCtx mocks base method.
func (m *MockBrandRedisRepository) GetByIdCtx(ctx context.Context, key string) (*models.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIdCtx", ctx, key)
	ret0, _ := ret[0].(*models.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIdCtx indicates an expected call of GetByIdCtx.
func (mr *MockBrandRedisRepositoryMockRecorder) GetByIdCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIdCtx", reflect.TypeOf((*MockBrandRedisRepository)(nil).GetByIdCtx), ctx, key)
}

// SetBrandCtx mocks base method.
func (m *MockBrandRedisRepository) SetBrandCtx(ctx context.Context, key string, seconds int, user *models.Brand) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBrandCtx", ctx, key, seconds, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetBrandCtx indicates an expected call of SetBrandCtx.
func (mr *MockBrandRedisRepositoryMockRecorder) SetBrandCtx(ctx, key, seconds, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBrandCtx", reflect.TypeOf((*MockBrandRedisRepository)(nil).SetBrandCtx), ctx, key, seconds, user)
}