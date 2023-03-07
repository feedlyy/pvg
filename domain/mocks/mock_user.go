// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	domain "pvg/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserRepository)(nil).Delete), ctx, id)
}

// Fetch mocks base method.
func (m *MockUserRepository) Fetch(ctx context.Context) ([]domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", ctx)
	ret0, _ := ret[0].([]domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockUserRepositoryMockRecorder) Fetch(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockUserRepository)(nil).Fetch), ctx)
}

// GetById mocks base method.
func (m *MockUserRepository) GetById(ctx context.Context, id int) (domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUserRepositoryMockRecorder) GetById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUserRepository)(nil).GetById), ctx, id)
}

// GetByUsername mocks base method.
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", ctx, username)
	ret0, _ := ret[0].(domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername.
func (mr *MockUserRepositoryMockRecorder) GetByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockUserRepository)(nil).GetByUsername), ctx, username)
}

// GetByUsrPhoneEmail mocks base method.
func (m *MockUserRepository) GetByUsrPhoneEmail(ctx context.Context, user domain.Users) (domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsrPhoneEmail", ctx, user)
	ret0, _ := ret[0].(domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsrPhoneEmail indicates an expected call of GetByUsrPhoneEmail.
func (mr *MockUserRepositoryMockRecorder) GetByUsrPhoneEmail(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsrPhoneEmail", reflect.TypeOf((*MockUserRepository)(nil).GetByUsrPhoneEmail), ctx, user)
}

// Insert mocks base method.
func (m *MockUserRepository) Insert(ctx context.Context, usr domain.Users) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, usr)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockUserRepositoryMockRecorder) Insert(ctx, usr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserRepository)(nil).Insert), ctx, usr)
}

// Update mocks base method.
func (m *MockUserRepository) Update(ctx context.Context, usr domain.Users) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, usr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(ctx, usr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), ctx, usr)
}

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// ActivateUser mocks base method.
func (m *MockUserService) ActivateUser(ctx context.Context, username string, code int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateUser", ctx, username, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateUser indicates an expected call of ActivateUser.
func (mr *MockUserServiceMockRecorder) ActivateUser(ctx, username, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateUser", reflect.TypeOf((*MockUserService)(nil).ActivateUser), ctx, username, code)
}

// CreateUser mocks base method.
func (m *MockUserService) CreateUser(ctx context.Context, usr domain.Users) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, usr)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserServiceMockRecorder) CreateUser(ctx, usr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserService)(nil).CreateUser), ctx, usr)
}

// DeleteUser mocks base method.
func (m *MockUserService) DeleteUser(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserServiceMockRecorder) DeleteUser(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserService)(nil).DeleteUser), ctx, id)
}

// GetAllUser mocks base method.
func (m *MockUserService) GetAllUser(ctx context.Context) ([]domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUser", ctx)
	ret0, _ := ret[0].([]domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUser indicates an expected call of GetAllUser.
func (mr *MockUserServiceMockRecorder) GetAllUser(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUser", reflect.TypeOf((*MockUserService)(nil).GetAllUser), ctx)
}

// GetUser mocks base method.
func (m *MockUserService) GetUser(ctx context.Context, username string) (domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, username)
	ret0, _ := ret[0].(domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserServiceMockRecorder) GetUser(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserService)(nil).GetUser), ctx, username)
}

// RequestActivationCode mocks base method.
func (m *MockUserService) RequestActivationCode(ctx context.Context, username string) (domain.ActivationCodes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestActivationCode", ctx, username)
	ret0, _ := ret[0].(domain.ActivationCodes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RequestActivationCode indicates an expected call of RequestActivationCode.
func (mr *MockUserServiceMockRecorder) RequestActivationCode(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestActivationCode", reflect.TypeOf((*MockUserService)(nil).RequestActivationCode), ctx, username)
}

// UpdateUser mocks base method.
func (m *MockUserService) UpdateUser(ctx context.Context, usr domain.Users) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, usr)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserServiceMockRecorder) UpdateUser(ctx, usr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserService)(nil).UpdateUser), ctx, usr)
}
