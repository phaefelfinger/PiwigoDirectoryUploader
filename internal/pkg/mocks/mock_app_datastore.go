// Code generated by MockGen. DO NOT EDIT.
// Source: git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/app (interfaces: ImageMetadataProvider)

// Package mocks is a generated GoMock package.
package mocks

import (
	app "git.haefelfinger.net/piwigo/PiwigoDirectoryUploader/internal/app"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockImageMetadataProvider is a mock of ImageMetadataProvider interface
type MockImageMetadataProvider struct {
	ctrl     *gomock.Controller
	recorder *MockImageMetadataProviderMockRecorder
}

// MockImageMetadataProviderMockRecorder is the mock recorder for MockImageMetadataProvider
type MockImageMetadataProviderMockRecorder struct {
	mock *MockImageMetadataProvider
}

// NewMockImageMetadataProvider creates a new mock instance
func NewMockImageMetadataProvider(ctrl *gomock.Controller) *MockImageMetadataProvider {
	mock := &MockImageMetadataProvider{ctrl: ctrl}
	mock.recorder = &MockImageMetadataProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImageMetadataProvider) EXPECT() *MockImageMetadataProviderMockRecorder {
	return m.recorder
}

// ImageMetadata mocks base method
func (m *MockImageMetadataProvider) ImageMetadata(arg0 string) (app.ImageMetaData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageMetadata", arg0)
	ret0, _ := ret[0].(app.ImageMetaData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageMetadata indicates an expected call of ImageMetadata
func (mr *MockImageMetadataProviderMockRecorder) ImageMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageMetadata", reflect.TypeOf((*MockImageMetadataProvider)(nil).ImageMetadata), arg0)
}

// ImageMetadataToUpload mocks base method
func (m *MockImageMetadataProvider) ImageMetadataToUpload() ([]app.ImageMetaData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageMetadataToUpload")
	ret0, _ := ret[0].([]app.ImageMetaData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageMetadataToUpload indicates an expected call of ImageMetadataToUpload
func (mr *MockImageMetadataProviderMockRecorder) ImageMetadataToUpload() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageMetadataToUpload", reflect.TypeOf((*MockImageMetadataProvider)(nil).ImageMetadataToUpload))
}

// SaveImageMetadata mocks base method
func (m *MockImageMetadataProvider) SaveImageMetadata(arg0 app.ImageMetaData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveImageMetadata", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveImageMetadata indicates an expected call of SaveImageMetadata
func (mr *MockImageMetadataProviderMockRecorder) SaveImageMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveImageMetadata", reflect.TypeOf((*MockImageMetadataProvider)(nil).SaveImageMetadata), arg0)
}

// SavePiwigoIdAndUpdateUploadFlag mocks base method
func (m *MockImageMetadataProvider) SavePiwigoIdAndUpdateUploadFlag(arg0 string, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePiwigoIdAndUpdateUploadFlag", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePiwigoIdAndUpdateUploadFlag indicates an expected call of SavePiwigoIdAndUpdateUploadFlag
func (mr *MockImageMetadataProviderMockRecorder) SavePiwigoIdAndUpdateUploadFlag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePiwigoIdAndUpdateUploadFlag", reflect.TypeOf((*MockImageMetadataProvider)(nil).SavePiwigoIdAndUpdateUploadFlag), arg0, arg1)
}
