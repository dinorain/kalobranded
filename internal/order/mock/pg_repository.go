// Code generated by MockGen. DO NOT EDIT.
// Source: pg_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/dinorain/checkoutaja/internal/models"
	utils "github.com/dinorain/checkoutaja/pkg/utils"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockOrderPGRepository is a mock of OrderPGRepository interface.
type MockOrderPGRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderPGRepositoryMockRecorder
}

// MockOrderPGRepositoryMockRecorder is the mock recorder for MockOrderPGRepository.
type MockOrderPGRepositoryMockRecorder struct {
	mock *MockOrderPGRepository
}

// NewMockOrderPGRepository creates a new mock instance.
func NewMockOrderPGRepository(ctrl *gomock.Controller) *MockOrderPGRepository {
	mock := &MockOrderPGRepository{ctrl: ctrl}
	mock.recorder = &MockOrderPGRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderPGRepository) EXPECT() *MockOrderPGRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockOrderPGRepository) Create(ctx context.Context, user *models.Order) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockOrderPGRepositoryMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderPGRepository)(nil).Create), ctx, user)
}

// DeleteById mocks base method.
func (m *MockOrderPGRepository) DeleteById(ctx context.Context, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockOrderPGRepositoryMockRecorder) DeleteById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockOrderPGRepository)(nil).DeleteById), ctx, userID)
}

// FindAll mocks base method.
func (m *MockOrderPGRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockOrderPGRepositoryMockRecorder) FindAll(ctx, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockOrderPGRepository)(nil).FindAll), ctx, pagination)
}

// FindAllBySellerId mocks base method.
func (m *MockOrderPGRepository) FindAllBySellerId(ctx context.Context, sellerID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllBySellerId", ctx, sellerID, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllBySellerId indicates an expected call of FindAllBySellerId.
func (mr *MockOrderPGRepositoryMockRecorder) FindAllBySellerId(ctx, sellerID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllBySellerId", reflect.TypeOf((*MockOrderPGRepository)(nil).FindAllBySellerId), ctx, sellerID, pagination)
}

// FindAllByUserId mocks base method.
func (m *MockOrderPGRepository) FindAllByUserId(ctx context.Context, userID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserId", ctx, userID, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserId indicates an expected call of FindAllByUserId.
func (mr *MockOrderPGRepositoryMockRecorder) FindAllByUserId(ctx, userID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserId", reflect.TypeOf((*MockOrderPGRepository)(nil).FindAllByUserId), ctx, userID, pagination)
}

// FindAllByUserIdSellerId mocks base method.
func (m *MockOrderPGRepository) FindAllByUserIdSellerId(ctx context.Context, userID, sellerID uuid.UUID, pagination *utils.Pagination) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserIdSellerId", ctx, userID, sellerID, pagination)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserIdSellerId indicates an expected call of FindAllByUserIdSellerId.
func (mr *MockOrderPGRepositoryMockRecorder) FindAllByUserIdSellerId(ctx, userID, sellerID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserIdSellerId", reflect.TypeOf((*MockOrderPGRepository)(nil).FindAllByUserIdSellerId), ctx, userID, sellerID, pagination)
}

// FindById mocks base method.
func (m *MockOrderPGRepository) FindById(ctx context.Context, userID uuid.UUID) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, userID)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockOrderPGRepositoryMockRecorder) FindById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockOrderPGRepository)(nil).FindById), ctx, userID)
}

// UpdateById mocks base method.
func (m *MockOrderPGRepository) UpdateById(ctx context.Context, user *models.Order) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", ctx, user)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockOrderPGRepositoryMockRecorder) UpdateById(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockOrderPGRepository)(nil).UpdateById), ctx, user)
}