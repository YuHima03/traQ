// Code generated by MockGen. DO NOT EDIT.
// Source: pin.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	model "github.com/traPtitech/traQ/model"
	reflect "reflect"
)

// MockPinRepository is a mock of PinRepository interface
type MockPinRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPinRepositoryMockRecorder
}

// MockPinRepositoryMockRecorder is the mock recorder for MockPinRepository
type MockPinRepositoryMockRecorder struct {
	mock *MockPinRepository
}

// NewMockPinRepository creates a new mock instance
func NewMockPinRepository(ctrl *gomock.Controller) *MockPinRepository {
	mock := &MockPinRepository{ctrl: ctrl}
	mock.recorder = &MockPinRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPinRepository) EXPECT() *MockPinRepositoryMockRecorder {
	return m.recorder
}

// PinMessage mocks base method
func (m *MockPinRepository) PinMessage(messageID, userID uuid.UUID) (*model.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PinMessage", messageID, userID)
	ret0, _ := ret[0].(*model.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PinMessage indicates an expected call of PinMessage
func (mr *MockPinRepositoryMockRecorder) PinMessage(messageID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PinMessage", reflect.TypeOf((*MockPinRepository)(nil).PinMessage), messageID, userID)
}

// UnpinMessage mocks base method
func (m *MockPinRepository) UnpinMessage(messageID uuid.UUID) (*model.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnpinMessage", messageID)
	ret0, _ := ret[0].(*model.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnpinMessage indicates an expected call of UnpinMessage
func (mr *MockPinRepositoryMockRecorder) UnpinMessage(messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpinMessage", reflect.TypeOf((*MockPinRepository)(nil).UnpinMessage), messageID)
}

// GetPinnedMessageByChannelID mocks base method
func (m *MockPinRepository) GetPinnedMessageByChannelID(channelID uuid.UUID) ([]*model.Pin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPinnedMessageByChannelID", channelID)
	ret0, _ := ret[0].([]*model.Pin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPinnedMessageByChannelID indicates an expected call of GetPinnedMessageByChannelID
func (mr *MockPinRepositoryMockRecorder) GetPinnedMessageByChannelID(channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPinnedMessageByChannelID", reflect.TypeOf((*MockPinRepository)(nil).GetPinnedMessageByChannelID), channelID)
}
