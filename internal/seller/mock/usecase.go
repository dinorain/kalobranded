// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

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

// MockSellerUseCase is a mock of SellerUseCase interface.
type MockSellerUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSellerUseCaseMockRecorder
}

// MockSellerUseCaseMockRecorder is the mock recorder for MockSellerUseCase.
type MockSellerUseCaseMockRecorder struct {
	mock *MockSellerUseCase
}

// NewMockSellerUseCase creates a new mock instance.
func NewMockSellerUseCase(ctrl *gomock.Controller) *MockSellerUseCase {
	mock := &MockSellerUseCase{ctrl: ctrl}
	mock.recorder = &MockSellerUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSellerUseCase) EXPECT() *MockSellerUseCaseMockRecorder {
	return m.recorder
}

// CachedFindById mocks base method.
func (m *MockSellerUseCase) CachedFindById(ctx context.Context, sellerID uuid.UUID) (*models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CachedFindById", ctx, sellerID)
	ret0, _ := ret[0].(*models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CachedFindById indicates an expected call of CachedFindById.
func (mr *MockSellerUseCaseMockRecorder) CachedFindById(ctx, sellerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CachedFindById", reflect.TypeOf((*MockSellerUseCase)(nil).CachedFindById), ctx, sellerID)
}

// DeleteById mocks base method.
func (m *MockSellerUseCase) DeleteById(ctx context.Context, sellerID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, sellerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockSellerUseCaseMockRecorder) DeleteById(ctx, sellerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockSellerUseCase)(nil).DeleteById), ctx, sellerID)
}

// FindAll mocks base method.
func (m *MockSellerUseCase) FindAll(ctx context.Context, pagination *utils.Pagination) ([]models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination)
	ret0, _ := ret[0].([]models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockSellerUseCaseMockRecorder) FindAll(ctx, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockSellerUseCase)(nil).FindAll), ctx, pagination)
}

// FindByEmail mocks base method.
func (m *MockSellerUseCase) FindByEmail(ctx context.Context, email string) (*models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(*models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockSellerUseCaseMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockSellerUseCase)(nil).FindByEmail), ctx, email)
}

// FindById mocks base method.
func (m *MockSellerUseCase) FindById(ctx context.Context, sellerID uuid.UUID) (*models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, sellerID)
	ret0, _ := ret[0].(*models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockSellerUseCaseMockRecorder) FindById(ctx, sellerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockSellerUseCase)(nil).FindById), ctx, sellerID)
}

// GenerateTokenPair mocks base method.
func (m *MockSellerUseCase) GenerateTokenPair(seller *models.Seller, sessionID string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokenPair", seller, sessionID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateTokenPair indicates an expected call of GenerateTokenPair.
func (mr *MockSellerUseCaseMockRecorder) GenerateTokenPair(seller, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokenPair", reflect.TypeOf((*MockSellerUseCase)(nil).GenerateTokenPair), seller, sessionID)
}

// Login mocks base method.
func (m *MockSellerUseCase) Login(ctx context.Context, email, password string) (*models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(*models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockSellerUseCaseMockRecorder) Login(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockSellerUseCase)(nil).Login), ctx, email, password)
}

// Register mocks base method.
func (m *MockSellerUseCase) Register(ctx context.Context, seller *models.Seller) (*models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, seller)
	ret0, _ := ret[0].(*models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockSellerUseCaseMockRecorder) Register(ctx, seller interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockSellerUseCase)(nil).Register), ctx, seller)
}

// UpdateById mocks base method.
func (m *MockSellerUseCase) UpdateById(ctx context.Context, seller *models.Seller) (*models.Seller, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", ctx, seller)
	ret0, _ := ret[0].(*models.Seller)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockSellerUseCaseMockRecorder) UpdateById(ctx, seller interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockSellerUseCase)(nil).UpdateById), ctx, seller)
}
