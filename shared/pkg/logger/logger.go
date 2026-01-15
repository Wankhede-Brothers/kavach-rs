// Package logger provides structured logging for kavach using slog.
// DACE: Uses Go 1.21+ standard library slog - no external dependencies.
// Reference: https://go.dev/blog/slog
package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

// Get returns the singleton logger instance.
func Get() *slog.Logger {
	once.Do(func() {
		defaultLogger = initLogger()
	})
	return defaultLogger
}

func initLogger() *slog.Logger {
	logDir := filepath.Join(os.Getenv("HOME"), ".local", "shared", "shared-ai", "memory", "logs")
	os.MkdirAll(logDir, 0755)

	logPath := filepath.Join(logDir, "kavach.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Fallback to stderr if can't open log file
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// Debug logs a debug message.
func Debug(component, msg string, args ...any) {
	Get().Debug(msg, append([]any{"component", component}, args...)...)
}

// Info logs an info message.
func Info(component, msg string, args ...any) {
	Get().Info(msg, append([]any{"component", component}, args...)...)
}

// Warn logs a warning message.
func Warn(component, msg string, args ...any) {
	Get().Warn(msg, append([]any{"component", component}, args...)...)
}

// Error logs an error message.
func Error(component, msg string, args ...any) {
	Get().Error(msg, append([]any{"component", component}, args...)...)
}

// WarnIfErr logs a warning if err is not nil, returns true if error occurred.
func WarnIfErr(component string, err error, msg string) bool {
	if err != nil {
		Warn(component, msg, "error", err.Error())
		return true
	}
	return false
}

// ErrorIfErr logs an error if err is not nil, returns true if error occurred.
func ErrorIfErr(component string, err error, msg string) bool {
	if err != nil {
		Error(component, msg, "error", err.Error())
		return true
	}
	return false
}
