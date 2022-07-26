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

// MockBrandPGRepository is a mock of BrandPGRepository interface.
type MockBrandPGRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBrandPGRepositoryMockRecorder
}

// MockBrandPGRepositoryMockRecorder is the mock recorder for MockBrandPGRepository.
type MockBrandPGRepositoryMockRecorder struct {
	mock *MockBrandPGRepository
}

// NewMockBrandPGRepository creates a new mock instance.
func NewMockBrandPGRepository(ctrl *gomock.Controller) *MockBrandPGRepository {
	mock := &MockBrandPGRepository{ctrl: ctrl}
	mock.recorder = &MockBrandPGRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBrandPGRepository) EXPECT() *MockBrandPGRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockBrandPGRepository) Create(ctx context.Context, user *models.Brand) (*models.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(*models.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockBrandPGRepositoryMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBrandPGRepository)(nil).Create), ctx, user)
}

// DeleteById mocks base method.
func (m *MockBrandPGRepository) DeleteById(ctx context.Context, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockBrandPGRepositoryMockRecorder) DeleteById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockBrandPGRepository)(nil).DeleteById), ctx, userID)
}

// FindAll mocks base method.
func (m *MockBrandPGRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination)
	ret0, _ := ret[0].([]models.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockBrandPGRepositoryMockRecorder) FindAll(ctx, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockBrandPGRepository)(nil).FindAll), ctx, pagination)
}

// FindById mocks base method.
func (m *MockBrandPGRepository) FindById(ctx context.Context, userID uuid.UUID) (*models.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, userID)
	ret0, _ := ret[0].(*models.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockBrandPGRepositoryMockRecorder) FindById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockBrandPGRepository)(nil).FindById), ctx, userID)
}

// UpdateById mocks base method.
func (m *MockBrandPGRepository) UpdateById(ctx context.Context, user *models.Brand) (*models.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", ctx, user)
	ret0, _ := ret[0].(*models.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockBrandPGRepositoryMockRecorder) UpdateById(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockBrandPGRepository)(nil).UpdateById), ctx, user)
}