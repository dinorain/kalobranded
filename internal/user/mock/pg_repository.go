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

// MockUserPGRepository is a mock of UserPGRepository interface.
type MockUserPGRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserPGRepositoryMockRecorder
}

// MockUserPGRepositoryMockRecorder is the mock recorder for MockUserPGRepository.
type MockUserPGRepositoryMockRecorder struct {
	mock *MockUserPGRepository
}

// NewMockUserPGRepository creates a new mock instance.
func NewMockUserPGRepository(ctrl *gomock.Controller) *MockUserPGRepository {
	mock := &MockUserPGRepository{ctrl: ctrl}
	mock.recorder = &MockUserPGRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserPGRepository) EXPECT() *MockUserPGRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserPGRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserPGRepositoryMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserPGRepository)(nil).Create), ctx, user)
}

// DeleteById mocks base method.
func (m *MockUserPGRepository) DeleteById(ctx context.Context, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockUserPGRepositoryMockRecorder) DeleteById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockUserPGRepository)(nil).DeleteById), ctx, userID)
}

// FindAll mocks base method.
func (m *MockUserPGRepository) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockUserPGRepositoryMockRecorder) FindAll(ctx, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockUserPGRepository)(nil).FindAll), ctx, pagination)
}

// FindByEmail mocks base method.
func (m *MockUserPGRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserPGRepositoryMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserPGRepository)(nil).FindByEmail), ctx, email)
}

// FindById mocks base method.
func (m *MockUserPGRepository) FindById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockUserPGRepositoryMockRecorder) FindById(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockUserPGRepository)(nil).FindById), ctx, userID)
}

// UpdateById mocks base method.
func (m *MockUserPGRepository) UpdateById(ctx context.Context, user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockUserPGRepositoryMockRecorder) UpdateById(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockUserPGRepository)(nil).UpdateById), ctx, user)
}