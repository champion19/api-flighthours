package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/champion19/api-flighthours/tools/utils"
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
	if os.MkdirAll(logDir, 0755) != nil {

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		logger := slog.New(handler)
		slog.SetDefault(logger)
		return &SlogLogger{}
	}

	logFileName := filepath.Join(logDir, time.Now().Format("backend-20060102.log"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		logger := slog.New(handler)
		slog.SetDefault(logger)
		return &SlogLogger{}
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &SlogLogger{logFile: logFile}
}

func (s *SlogLogger) WithTraceID(traceID string) Logger {
	return &SlogLogger{
		logFile: s.logFile,
		traceID: traceID,
	}
}

func (s *SlogLogger) enrichWithContext(args ...any) []any {
	if s.traceID != "" {
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
	slog.Info("[SUCCESS] "+msg, enrichedArgs...)
}

func (s *SlogLogger) Warn(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Warn(msg, enrichedArgs...)
}

func (s *SlogLogger) Fatal(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Error("[FATAL] "+msg, enrichedArgs...)
}

func (s *SlogLogger) Panic(msg string, args ...any) {
	enrichedArgs := s.enrichWithContext(args...)
	slog.Error("[PANIC] "+msg, enrichedArgs...)
}
