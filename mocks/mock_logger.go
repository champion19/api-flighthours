package mocks

import (
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of logger.Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Debug(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Success(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Fatal(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Panic(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) WithTraceID(traceID string) logger.Logger {
	args := m.Called(traceID)
	if args.Get(0) == nil {
		return m // Return self if no mock setup
	}
	return args.Get(0).(logger.Logger)
}
