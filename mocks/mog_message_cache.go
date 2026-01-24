package mocks

import (
	"context"

	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
	"github.com/stretchr/testify/mock"
)

// MockMessageCache is a mock implementation of MessageCache
type MockMessageCache struct {
	mock.Mock
}

func (m *MockMessageCache) LoadMessages(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageCache) ReloadMessages(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageCache) StartAutoRefresh(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockMessageCache) StopAutoRefresh() {
	m.Called()
}

func (m *MockMessageCache) GetMessage(code string) *cachetypes.CachedMessage {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*cachetypes.CachedMessage)
}

func (m *MockMessageCache) GetMessageResponse(code string, params ...string) *cachetypes.MessageResponse {
	args := m.Called(code, params)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*cachetypes.MessageResponse)
}

func (m *MockMessageCache) GetHTTPStatus(code string) int {
	args := m.Called(code)
	return args.Int(0)
}

func (m *MockMessageCache) MessageCount() int {
	args := m.Called()
	return args.Int(0)
}
