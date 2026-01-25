package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/champion19/flighthours-api/tools/utils"
)

type SlogLogger struct {
	logFile *os.File
	traceID string
}

func NewSlogLogger() *SlogLogger {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		if root, err := utils.FindModuleRoot(); err == nil {
			logDir = filepath.Join(root, "logs")
		} else {
			logDir = filepath.Join("/var/log/flighthours")
		}
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// If we can't create the directory, just log to stdout
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		logger := slog.New(handler)
		slog.SetDefault(logger)
		return &SlogLogger{}
	}

	// Create log file with date in name
	logFileName := filepath.Join(logDir, time.Now().Format("backend-20060102.log"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// If we can't open the file, just log to stdout
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		logger := slog.New(handler)
		slog.SetDefault(logger)
		return &SlogLogger{}
	}

	// Create multi-writer to write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Configure JSON handler for structured logging
	// This format is easily parsed by Loki/Promtail
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // Allow all log levels
		AddSource: true,            // Include file, line, and function info
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &SlogLogger{logFile: logFile}
}

// WithTraceID returns a new Logger instance with the traceID pre-configured.
// All logs from this logger will include traceID=<value> for correlation in Loki.
func (s *SlogLogger) WithTraceID(traceID string) Logger {
	return &SlogLogger{
		logFile: s.logFile,
		traceID: traceID,
	}
}

// enrichWithContext prepends traceID to the args if it's set.
// This ensures all logs include the traceID for Loki correlation.
func (s *SlogLogger) enrichWithContext(args ...any) []any {
	if s.traceID != "" {
		// Prepend traceID=<value> to match Loki's derivedField regex: traceID=(\w+)
		return append([]any{"traceID", s.traceID}, args...)
	}
	return args
}

func (s *SlogLogger) Info(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Info(msg, enrichedArgs...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Error(msg, enrichedArgs...)
}

func (s *SlogLogger) Debug(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Debug(msg, enrichedArgs...)
}

func (s *SlogLogger) Success(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Info(msg, enrichedArgs...)
}

func (s *SlogLogger) Warn(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Warn(msg, enrichedArgs...)
}

func (s *SlogLogger) Fatal(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Error(msg, enrichedArgs...)
}

func (s *SlogLogger) Panic(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Error(msg, enrichedArgs...)
}
