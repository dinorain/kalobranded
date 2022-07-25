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

// MockSellerRedisRepository is a mock of SellerRedisRepository interface.
type MockSellerRedisRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSellerRedisRepositoryMockRecorder
}

// MockSellerRedisRepositoryMockRecorder is the mock recorder for MockSellerRedisRepository.
type MockSellerRedisRepositoryMockRecorder struct {
	mock *MockSellerRedisRepository
}

// NewMockSellerRedisRepository creates a new mock instance.
func NewMockSellerRedisRepository(ctrl *gomock.Controller) *MockSellerRedisRepository {
	mock := &MockSellerRedisRepository{ctrl: ctrl}
	mock.recorder = &MockSellerRedisRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSellerRedisRepository) EXPECT() *MockSellerRedisRepositoryMockRecorder {
	return m.recorder
}

// DeleteSellerCtx mocks base method.
func (m *MockSellerRedisRepository) DeleteSellerCtx(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSellerCtx", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSellerCtx indicates an expected call of DeleteSellerCtx.
func (mr *MockSellerRedisRepositoryMockRecorder) DeleteSellerCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSellerCtx", reflect.TypeOf((*MockSellerRedisRepository)(nil).DeleteSellerCtx), ctx, key)
}

// GetByIdCtx mocks base method.
func (m *MockSellerRedisRepository) GetByIdCtx(ctx context.Context, key string) (*models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIdCtx", ctx, key)
	ret0, _ := ret[0].(*models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIdCtx indicates an expected call of GetByIdCtx.
func (mr *MockSellerRedisRepositoryMockRecorder) GetByIdCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIdCtx", reflect.TypeOf((*MockSellerRedisRepository)(nil).GetByIdCtx), ctx, key)
}

// SetSellerCtx mocks base method.
func (m *MockSellerRedisRepository) SetSellerCtx(ctx context.Context, key string, seconds int, user *models.Seller) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSellerCtx", ctx, key, seconds, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSellerCtx indicates an expected call of SetSellerCtx.
func (mr *MockSellerRedisRepositoryMockRecorder) SetSellerCtx(ctx, key, seconds, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSellerCtx", reflect.TypeOf((*MockSellerRedisRepository)(nil).SetSellerCtx), ctx, key, seconds, user)
}