// Code generated by MockGen. DO NOT EDIT.
// Source: git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo (interfaces: CategoryApi,ImageApi)

// Package images is a generated GoMock package.
package images

import (
	piwigo "git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/pkg/piwigo"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCategoryApi is a mock of CategoryApi interface
type MockCategoryApi struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryApiMockRecorder
}

// MockCategoryApiMockRecorder is the mock recorder for MockCategoryApi
type MockCategoryApiMockRecorder struct {
	mock *MockCategoryApi
}

// NewMockCategoryApi creates a new mock instance
func NewMockCategoryApi(ctrl *gomock.Controller) *MockCategoryApi {
	mock := &MockCategoryApi{ctrl: ctrl}
	mock.recorder = &MockCategoryApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCategoryApi) EXPECT() *MockCategoryApiMockRecorder {
	return m.recorder
}

// CreateCategory mocks base method
func (m *MockCategoryApi) CreateCategory(arg0 int, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory
func (mr *MockCategoryApiMockRecorder) CreateCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockCategoryApi)(nil).CreateCategory), arg0, arg1)
}

// GetAllCategories mocks base method
func (m *MockCategoryApi) GetAllCategories() (map[string]*piwigo.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCategories")
	ret0, _ := ret[0].(map[string]*piwigo.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCategories indicates an expected call of GetAllCategories
func (mr *MockCategoryApiMockRecorder) GetAllCategories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCategories", reflect.TypeOf((*MockCategoryApi)(nil).GetAllCategories))
}

// MockImageApi is a mock of ImageApi interface
type MockImageApi struct {
	ctrl     *gomock.Controller
	recorder *MockImageApiMockRecorder
}

// MockImageApiMockRecorder is the mock recorder for MockImageApi
type MockImageApiMockRecorder struct {
	mock *MockImageApi
}

// NewMockImageApi creates a new mock instance
func NewMockImageApi(ctrl *gomock.Controller) *MockImageApi {
	mock := &MockImageApi{ctrl: ctrl}
	mock.recorder = &MockImageApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImageApi) EXPECT() *MockImageApiMockRecorder {
	return m.recorder
}

// DeleteImages mocks base method
func (m *MockImageApi) DeleteImages(arg0 []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImages", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImages indicates an expected call of DeleteImages
func (mr *MockImageApiMockRecorder) DeleteImages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImages", reflect.TypeOf((*MockImageApi)(nil).DeleteImages), arg0)
}

// ImageCheckFile mocks base method
func (m *MockImageApi) ImageCheckFile(arg0 int, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageCheckFile", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageCheckFile indicates an expected call of ImageCheckFile
func (mr *MockImageApiMockRecorder) ImageCheckFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageCheckFile", reflect.TypeOf((*MockImageApi)(nil).ImageCheckFile), arg0, arg1)
}

// ImagesExistOnPiwigo mocks base method
func (m *MockImageApi) ImagesExistOnPiwigo(arg0 []string) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImagesExistOnPiwigo", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImagesExistOnPiwigo indicates an expected call of ImagesExistOnPiwigo
func (mr *MockImageApiMockRecorder) ImagesExistOnPiwigo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImagesExistOnPiwigo", reflect.TypeOf((*MockImageApi)(nil).ImagesExistOnPiwigo), arg0)
}

// UploadImage mocks base method
func (m *MockImageApi) UploadImage(arg0 int, arg1, arg2 string, arg3 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage
func (mr *MockImageApiMockRecorder) UploadImage(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockImageApi)(nil).UploadImage), arg0, arg1, arg2, arg3)
}
