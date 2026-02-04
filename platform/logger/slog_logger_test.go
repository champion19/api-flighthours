package logger

import (
	"os"
	"testing"
)

func TestNewSlogLogger(t *testing.T) {
	t.Run("creates logger instance", func(t *testing.T) {
		logger := NewSlogLogger()
		if logger == nil {
			t.Error("expected non-nil logger")
		}
	})
}

func TestNewSlogLogger_WithLogDir(t *testing.T) {
	t.Run("creates logger with LOG_DIR env set", func(t *testing.T) {
		// Create a temp directory
		tmpDir := t.TempDir()
		os.Setenv("LOG_DIR", tmpDir)
		defer os.Unsetenv("LOG_DIR")

		logger := NewSlogLogger()
		if logger == nil {
			t.Error("expected non-nil logger")
		}
	})

	t.Run("creates logger when LOG_DIR is invalid", func(t *testing.T) {
		// Set LOG_DIR to a path that cannot be created
		os.Setenv("LOG_DIR", "/this/path/should/not/exist/or/be/creatable/by/test")
		defer os.Unsetenv("LOG_DIR")

		logger := NewSlogLogger()
		// Should still return a logger (console-only fallback)
		if logger == nil {
			t.Error("expected non-nil logger even with invalid LOG_DIR")
		}
	})
}

func TestWithTraceID(t *testing.T) {
	logger := NewSlogLogger()

	t.Run("creates new logger with traceID", func(t *testing.T) {
		traceID := "test-trace-123"
		loggerWithTrace := logger.WithTraceID(traceID)
		if loggerWithTrace == nil {
			t.Error("expected non-nil logger")
		}
		// Type assertion to verify it's SlogLogger
		if _, ok := loggerWithTrace.(*SlogLogger); !ok {
			t.Error("expected SlogLogger type")
		}
	})
}

func TestLogMethods(t *testing.T) {
	logger := NewSlogLogger()

	// These tests just ensure no panics occur
	t.Run("Info does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Info panicked: %v", r)
			}
		}()
		logger.Info("test message", "key", "value")
	})

	t.Run("Error does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Error panicked: %v", r)
			}
		}()
		logger.Error("test error", "error", "some error")
	})

	t.Run("Debug does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Debug panicked: %v", r)
			}
		}()
		logger.Debug("debug message", "detail", "value")
	})

	t.Run("Success does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Success panicked: %v", r)
			}
		}()
		logger.Success("success message")
	})

	t.Run("Warn does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Warn panicked: %v", r)
			}
		}()
		logger.Warn("warning message", "level", "high")
	})

	t.Run("Fatal does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Fatal panicked: %v", r)
			}
		}()
		logger.Fatal("fatal message")
	})

	t.Run("Panic method does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic method panicked: %v", r)
			}
		}()
		logger.Panic("panic message")
	})
}

func TestLogMethodsWithTraceID(t *testing.T) {
	logger := NewSlogLogger()
	loggerWithTrace := logger.WithTraceID("trace-123")

	t.Run("Info with traceID", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Info panicked: %v", r)
			}
		}()
		loggerWithTrace.Info("message with trace")
	})

	t.Run("Error with traceID", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Error panicked: %v", r)
			}
		}()
		loggerWithTrace.Error("error with trace")
	})
}
