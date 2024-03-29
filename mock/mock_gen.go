// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fox-one/holder/core (interfaces: AssetService,UserStore,UserService,MessageStore,MessageService,Notifier,PoolStore,VaultStore,WalletStore)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	core "github.com/fox-one/holder/core"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAssetService is a mock of AssetService interface
type MockAssetService struct {
	ctrl     *gomock.Controller
	recorder *MockAssetServiceMockRecorder
}

// MockAssetServiceMockRecorder is the mock recorder for MockAssetService
type MockAssetServiceMockRecorder struct {
	mock *MockAssetService
}

// NewMockAssetService creates a new mock instance
func NewMockAssetService(ctrl *gomock.Controller) *MockAssetService {
	mock := &MockAssetService{ctrl: ctrl}
	mock.recorder = &MockAssetServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAssetService) EXPECT() *MockAssetServiceMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockAssetService) Find(arg0 context.Context, arg1 string) (*core.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockAssetServiceMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockAssetService)(nil).Find), arg0, arg1)
}

// List mocks base method
func (m *MockAssetService) List(arg0 context.Context) ([]*core.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*core.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockAssetServiceMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAssetService)(nil).List), arg0)
}

// MockUserStore is a mock of UserStore interface
type MockUserStore struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoreMockRecorder
}

// MockUserStoreMockRecorder is the mock recorder for MockUserStore
type MockUserStoreMockRecorder struct {
	mock *MockUserStore
}

// NewMockUserStore creates a new mock instance
func NewMockUserStore(ctrl *gomock.Controller) *MockUserStore {
	mock := &MockUserStore{ctrl: ctrl}
	mock.recorder = &MockUserStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserStore) EXPECT() *MockUserStoreMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockUserStore) Find(arg0 context.Context, arg1 string) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockUserStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserStore)(nil).Find), arg0, arg1)
}

// Save mocks base method
func (m *MockUserStore) Save(arg0 context.Context, arg1 *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockUserStoreMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserStore)(nil).Save), arg0, arg1)
}

// MockUserService is a mock of UserService interface
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// Auth mocks base method
func (m *MockUserService) Auth(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Auth indicates an expected call of Auth
func (mr *MockUserServiceMockRecorder) Auth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockUserService)(nil).Auth), arg0, arg1)
}

// Find mocks base method
func (m *MockUserService) Find(arg0 context.Context, arg1 string) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockUserServiceMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUserService)(nil).Find), arg0, arg1)
}

// Login mocks base method
func (m *MockUserService) Login(arg0 context.Context, arg1 string) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login
func (mr *MockUserServiceMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), arg0, arg1)
}

// MockMessageStore is a mock of MessageStore interface
type MockMessageStore struct {
	ctrl     *gomock.Controller
	recorder *MockMessageStoreMockRecorder
}

// MockMessageStoreMockRecorder is the mock recorder for MockMessageStore
type MockMessageStoreMockRecorder struct {
	mock *MockMessageStore
}

// NewMockMessageStore creates a new mock instance
func NewMockMessageStore(ctrl *gomock.Controller) *MockMessageStore {
	mock := &MockMessageStore{ctrl: ctrl}
	mock.recorder = &MockMessageStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessageStore) EXPECT() *MockMessageStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockMessageStore) Create(arg0 context.Context, arg1 []*core.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockMessageStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMessageStore)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockMessageStore) Delete(arg0 context.Context, arg1 []*core.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockMessageStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMessageStore)(nil).Delete), arg0, arg1)
}

// List mocks base method
func (m *MockMessageStore) List(arg0 context.Context, arg1 int) ([]*core.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*core.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockMessageStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockMessageStore)(nil).List), arg0, arg1)
}

// MockMessageService is a mock of MessageService interface
type MockMessageService struct {
	ctrl     *gomock.Controller
	recorder *MockMessageServiceMockRecorder
}

// MockMessageServiceMockRecorder is the mock recorder for MockMessageService
type MockMessageServiceMockRecorder struct {
	mock *MockMessageService
}

// NewMockMessageService creates a new mock instance
func NewMockMessageService(ctrl *gomock.Controller) *MockMessageService {
	mock := &MockMessageService{ctrl: ctrl}
	mock.recorder = &MockMessageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessageService) EXPECT() *MockMessageServiceMockRecorder {
	return m.recorder
}

// Meet mocks base method
func (m *MockMessageService) Meet(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Meet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Meet indicates an expected call of Meet
func (mr *MockMessageServiceMockRecorder) Meet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Meet", reflect.TypeOf((*MockMessageService)(nil).Meet), arg0, arg1)
}

// Send mocks base method
func (m *MockMessageService) Send(arg0 context.Context, arg1 []*core.Message, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockMessageServiceMockRecorder) Send(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockMessageService)(nil).Send), arg0, arg1, arg2)
}

// MockNotifier is a mock of Notifier interface
type MockNotifier struct {
	ctrl     *gomock.Controller
	recorder *MockNotifierMockRecorder
}

// MockNotifierMockRecorder is the mock recorder for MockNotifier
type MockNotifierMockRecorder struct {
	mock *MockNotifier
}

// NewMockNotifier creates a new mock instance
func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
	mock := &MockNotifier{ctrl: ctrl}
	mock.recorder = &MockNotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNotifier) EXPECT() *MockNotifierMockRecorder {
	return m.recorder
}

// Auth mocks base method
func (m *MockNotifier) Auth(arg0 context.Context, arg1 *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Auth indicates an expected call of Auth
func (mr *MockNotifierMockRecorder) Auth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockNotifier)(nil).Auth), arg0, arg1)
}

// LockDone mocks base method
func (m *MockNotifier) LockDone(arg0 context.Context, arg1 *core.Pool, arg2 *core.Vault) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockDone", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// LockDone indicates an expected call of LockDone
func (mr *MockNotifierMockRecorder) LockDone(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockDone", reflect.TypeOf((*MockNotifier)(nil).LockDone), arg0, arg1, arg2)
}

// Snapshot mocks base method
func (m *MockNotifier) Snapshot(arg0 context.Context, arg1 *core.Transfer, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Snapshot", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Snapshot indicates an expected call of Snapshot
func (mr *MockNotifierMockRecorder) Snapshot(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Snapshot", reflect.TypeOf((*MockNotifier)(nil).Snapshot), arg0, arg1, arg2)
}

// Transaction mocks base method
func (m *MockNotifier) Transaction(arg0 context.Context, arg1 *core.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction
func (mr *MockNotifierMockRecorder) Transaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockNotifier)(nil).Transaction), arg0, arg1)
}

// MockPoolStore is a mock of PoolStore interface
type MockPoolStore struct {
	ctrl     *gomock.Controller
	recorder *MockPoolStoreMockRecorder
}

// MockPoolStoreMockRecorder is the mock recorder for MockPoolStore
type MockPoolStoreMockRecorder struct {
	mock *MockPoolStore
}

// NewMockPoolStore creates a new mock instance
func NewMockPoolStore(ctrl *gomock.Controller) *MockPoolStore {
	mock := &MockPoolStore{ctrl: ctrl}
	mock.recorder = &MockPoolStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPoolStore) EXPECT() *MockPoolStoreMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockPoolStore) Find(arg0 context.Context, arg1 string) (*core.Pool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Pool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockPoolStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockPoolStore)(nil).Find), arg0, arg1)
}

// List mocks base method
func (m *MockPoolStore) List(arg0 context.Context) ([]*core.Pool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*core.Pool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockPoolStoreMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPoolStore)(nil).List), arg0)
}

// Save mocks base method
func (m *MockPoolStore) Save(arg0 context.Context, arg1 *core.Pool, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockPoolStoreMockRecorder) Save(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPoolStore)(nil).Save), arg0, arg1, arg2)
}

// UpdateInfo mocks base method
func (m *MockPoolStore) UpdateInfo(arg0 context.Context, arg1 *core.Pool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInfo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInfo indicates an expected call of UpdateInfo
func (mr *MockPoolStoreMockRecorder) UpdateInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInfo", reflect.TypeOf((*MockPoolStore)(nil).UpdateInfo), arg0, arg1)
}

// MockVaultStore is a mock of VaultStore interface
type MockVaultStore struct {
	ctrl     *gomock.Controller
	recorder *MockVaultStoreMockRecorder
}

// MockVaultStoreMockRecorder is the mock recorder for MockVaultStore
type MockVaultStoreMockRecorder struct {
	mock *MockVaultStore
}

// NewMockVaultStore creates a new mock instance
func NewMockVaultStore(ctrl *gomock.Controller) *MockVaultStore {
	mock := &MockVaultStore{ctrl: ctrl}
	mock.recorder = &MockVaultStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVaultStore) EXPECT() *MockVaultStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockVaultStore) Create(arg0 context.Context, arg1 *core.Vault) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockVaultStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVaultStore)(nil).Create), arg0, arg1)
}

// Find mocks base method
func (m *MockVaultStore) Find(arg0 context.Context, arg1 string) (*core.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockVaultStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockVaultStore)(nil).Find), arg0, arg1)
}

// List mocks base method
func (m *MockVaultStore) List(arg0 context.Context, arg1 int64, arg2 int) ([]*core.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*core.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockVaultStoreMockRecorder) List(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockVaultStore)(nil).List), arg0, arg1, arg2)
}

// ListUser mocks base method
func (m *MockVaultStore) ListUser(arg0 context.Context, arg1 string) ([]*core.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUser", arg0, arg1)
	ret0, _ := ret[0].([]*core.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUser indicates an expected call of ListUser
func (mr *MockVaultStoreMockRecorder) ListUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUser", reflect.TypeOf((*MockVaultStore)(nil).ListUser), arg0, arg1)
}

// Update mocks base method
func (m *MockVaultStore) Update(arg0 context.Context, arg1 *core.Vault, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockVaultStoreMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVaultStore)(nil).Update), arg0, arg1, arg2)
}

// MockWalletStore is a mock of WalletStore interface
type MockWalletStore struct {
	ctrl     *gomock.Controller
	recorder *MockWalletStoreMockRecorder
}

// MockWalletStoreMockRecorder is the mock recorder for MockWalletStore
type MockWalletStoreMockRecorder struct {
	mock *MockWalletStore
}

// NewMockWalletStore creates a new mock instance
func NewMockWalletStore(ctrl *gomock.Controller) *MockWalletStore {
	mock := &MockWalletStore{ctrl: ctrl}
	mock.recorder = &MockWalletStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWalletStore) EXPECT() *MockWalletStoreMockRecorder {
	return m.recorder
}

// Assign mocks base method
func (m *MockWalletStore) Assign(arg0 context.Context, arg1 []*core.Output, arg2 *core.Transfer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Assign", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Assign indicates an expected call of Assign
func (mr *MockWalletStoreMockRecorder) Assign(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Assign", reflect.TypeOf((*MockWalletStore)(nil).Assign), arg0, arg1, arg2)
}

// CountOutputs mocks base method
func (m *MockWalletStore) CountOutputs(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountOutputs", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountOutputs indicates an expected call of CountOutputs
func (mr *MockWalletStoreMockRecorder) CountOutputs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountOutputs", reflect.TypeOf((*MockWalletStore)(nil).CountOutputs), arg0)
}

// CountUnhandledTransfers mocks base method
func (m *MockWalletStore) CountUnhandledTransfers(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUnhandledTransfers", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUnhandledTransfers indicates an expected call of CountUnhandledTransfers
func (mr *MockWalletStoreMockRecorder) CountUnhandledTransfers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUnhandledTransfers", reflect.TypeOf((*MockWalletStore)(nil).CountUnhandledTransfers), arg0)
}

// CreateRawTransaction mocks base method
func (m *MockWalletStore) CreateRawTransaction(arg0 context.Context, arg1 *core.RawTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRawTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRawTransaction indicates an expected call of CreateRawTransaction
func (mr *MockWalletStoreMockRecorder) CreateRawTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRawTransaction", reflect.TypeOf((*MockWalletStore)(nil).CreateRawTransaction), arg0, arg1)
}

// CreateTransfers mocks base method
func (m *MockWalletStore) CreateTransfers(arg0 context.Context, arg1 []*core.Transfer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransfers", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTransfers indicates an expected call of CreateTransfers
func (mr *MockWalletStoreMockRecorder) CreateTransfers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransfers", reflect.TypeOf((*MockWalletStore)(nil).CreateTransfers), arg0, arg1)
}

// ExpireRawTransaction mocks base method
func (m *MockWalletStore) ExpireRawTransaction(arg0 context.Context, arg1 *core.RawTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpireRawTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExpireRawTransaction indicates an expected call of ExpireRawTransaction
func (mr *MockWalletStoreMockRecorder) ExpireRawTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpireRawTransaction", reflect.TypeOf((*MockWalletStore)(nil).ExpireRawTransaction), arg0, arg1)
}

// FindSpentBy mocks base method
func (m *MockWalletStore) FindSpentBy(arg0 context.Context, arg1, arg2 string) (*core.Output, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSpentBy", arg0, arg1, arg2)
	ret0, _ := ret[0].(*core.Output)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSpentBy indicates an expected call of FindSpentBy
func (mr *MockWalletStoreMockRecorder) FindSpentBy(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSpentBy", reflect.TypeOf((*MockWalletStore)(nil).FindSpentBy), arg0, arg1, arg2)
}

// List mocks base method
func (m *MockWalletStore) List(arg0 context.Context, arg1 int64, arg2 int) ([]*core.Output, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*core.Output)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockWalletStoreMockRecorder) List(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockWalletStore)(nil).List), arg0, arg1, arg2)
}

// ListPendingRawTransactions mocks base method
func (m *MockWalletStore) ListPendingRawTransactions(arg0 context.Context, arg1 int) ([]*core.RawTransaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPendingRawTransactions", arg0, arg1)
	ret0, _ := ret[0].([]*core.RawTransaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPendingRawTransactions indicates an expected call of ListPendingRawTransactions
func (mr *MockWalletStoreMockRecorder) ListPendingRawTransactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPendingRawTransactions", reflect.TypeOf((*MockWalletStore)(nil).ListPendingRawTransactions), arg0, arg1)
}

// ListSpentBy mocks base method
func (m *MockWalletStore) ListSpentBy(arg0 context.Context, arg1, arg2 string) ([]*core.Output, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSpentBy", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*core.Output)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSpentBy indicates an expected call of ListSpentBy
func (mr *MockWalletStoreMockRecorder) ListSpentBy(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSpentBy", reflect.TypeOf((*MockWalletStore)(nil).ListSpentBy), arg0, arg1, arg2)
}

// ListTransfers mocks base method
func (m *MockWalletStore) ListTransfers(arg0 context.Context, arg1 core.TransferStatus, arg2 int) ([]*core.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTransfers", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*core.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTransfers indicates an expected call of ListTransfers
func (mr *MockWalletStoreMockRecorder) ListTransfers(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTransfers", reflect.TypeOf((*MockWalletStore)(nil).ListTransfers), arg0, arg1, arg2)
}

// ListUnspent mocks base method
func (m *MockWalletStore) ListUnspent(arg0 context.Context, arg1 string, arg2 int) ([]*core.Output, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUnspent", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*core.Output)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUnspent indicates an expected call of ListUnspent
func (mr *MockWalletStoreMockRecorder) ListUnspent(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUnspent", reflect.TypeOf((*MockWalletStore)(nil).ListUnspent), arg0, arg1, arg2)
}

// Save mocks base method
func (m *MockWalletStore) Save(arg0 context.Context, arg1 []*core.Output, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockWalletStoreMockRecorder) Save(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockWalletStore)(nil).Save), arg0, arg1, arg2)
}

// UpdateTransfer mocks base method
func (m *MockWalletStore) UpdateTransfer(arg0 context.Context, arg1 *core.Transfer, arg2 core.TransferStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransfer", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTransfer indicates an expected call of UpdateTransfer
func (mr *MockWalletStoreMockRecorder) UpdateTransfer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransfer", reflect.TypeOf((*MockWalletStore)(nil).UpdateTransfer), arg0, arg1, arg2)
}
