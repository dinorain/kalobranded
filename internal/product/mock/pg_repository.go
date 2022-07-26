// Code generated by MockGen. DO NOT EDIT.
// Source: pg_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/dinorain/kalobranded/internal/models"
	utils "github.com/dinorain/kalobranded/pkg/utils"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockProductPGRepository is a mock of ProductPGRepository interface.
type MockProductPGRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductPGRepositoryMockRecorder
}

// MockProductPGRepositoryMockRecorder is the mock recorder for MockProductPGRepository.
type MockProductPGRepositoryMockRecorder struct {
	mock *MockProductPGRepository
}

// NewMockProductPGRepository creates a new mock instance.
func NewMockProductPGRepository(ctrl *gomock.Controller) *MockProductPGRepository {
	mock := &MockProductPGRepository{ctrl: ctrl}
	mock.recorder = &MockProductPGRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductPGRepository) EXPECT() *MockProductPGRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductPGRepository) Create(ctx context.Context, user *models.Product) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProductPGRepositoryMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductPGRepository)(nil).Create), ctx, user)
}

// DeleteById mocks base method.
func (m *MockProductPGRepository) DeleteById(ctx context.Context, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockProductPGRepositoryMockRecorder) DeleteById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockProductPGRepository)(nil).DeleteById), ctx, userID)
}

// FindAll mocks base method.
func (m *MockProductPGRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination)
	ret0, _ := ret[0].([]models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockProductPGRepositoryMockRecorder) FindAll(ctx, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockProductPGRepository)(nil).FindAll), ctx, pagination)
}

// FindAllByBrandId mocks base method.
func (m *MockProductPGRepository) FindAllByBrandId(ctx context.Context, brandID uuid.UUID, pagination *utils.Pagination) ([]models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByBrandId", ctx, brandID, pagination)
	ret0, _ := ret[0].([]models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByBrandId indicates an expected call of FindAllByBrandId.
func (mr *MockProductPGRepositoryMockRecorder) FindAllByBrandId(ctx, brandID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByBrandId", reflect.TypeOf((*MockProductPGRepository)(nil).FindAllByBrandId), ctx, brandID, pagination)
}

// FindById mocks base method.
func (m *MockProductPGRepository) FindById(ctx context.Context, userID uuid.UUID) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, userID)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockProductPGRepositoryMockRecorder) FindById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockProductPGRepository)(nil).FindById), ctx, userID)
}

// UpdateById mocks base method.
func (m *MockProductPGRepository) UpdateById(ctx context.Context, user *models.Product) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", ctx, user)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockProductPGRepositoryMockRecorder) UpdateById(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockProductPGRepository)(nil).UpdateById), ctx, user)
}
