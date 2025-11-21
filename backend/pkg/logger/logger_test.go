package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "Valid console logger",
			config: Config{
				Level:      "info",
				OutputPath: "stdout",
				Format:     "console",
			},
			wantErr: false,
		},
		{
			name: "Valid JSON logger",
			config: Config{
				Level:      "debug",
				OutputPath: "stdout",
				Format:     "json",
			},
			wantErr: false,
		},
		{
			name: "Valid error level",
			config: Config{
				Level:      "error",
				OutputPath: "stdout",
				Format:     "console",
			},
			wantErr: false,
		},
		{
			name: "Invalid output path",
			config: Config{
				Level:      "info",
				OutputPath: "/invalid/path/that/does/not/exist/test.log",
				Format:     "console",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("New() returned nil logger without error")
			}
			if logger != nil {
				defer logger.Sync()
			}
		})
	}
}

func TestNewDefault(t *testing.T) {
	logger := NewDefault()
	if logger == nil {
		t.Error("NewDefault() returned nil")
	}
	defer logger.Sync()
}

func TestNewProduction(t *testing.T) {
	logger, err := NewProduction()
	if err != nil {
		t.Errorf("NewProduction() error = %v", err)
	}
	if logger == nil {
		t.Error("NewProduction() returned nil")
	}
	defer logger.Sync()
}

func TestNewDevelopment(t *testing.T) {
	logger, err := NewDevelopment()
	if err != nil {
		t.Errorf("NewDevelopment() error = %v", err)
	}
	if logger == nil {
		t.Error("NewDevelopment() returned nil")
	}
	defer logger.Sync()
}

func TestWithField(t *testing.T) {
	logger := NewDefault()
	defer logger.Sync()

	fieldLogger := logger.WithField("test_key", "test_value")
	if fieldLogger == nil {
		t.Error("WithField() returned nil")
	}
}

func TestWithFields(t *testing.T) {
	logger := NewDefault()
	defer logger.Sync()

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	fieldLogger := logger.WithFields(fields)
	if fieldLogger == nil {
		t.Error("WithFields() returned nil")
	}
}

func TestWithError(t *testing.T) {
	logger := NewDefault()
	defer logger.Sync()

	err := os.ErrNotExist
	errLogger := logger.WithError(err)
	if errLogger == nil {
		t.Error("WithError() returned nil")
	}
}

func TestContext(t *testing.T) {
	logger := NewDefault()
	defer logger.Sync()

	fields := []zap.Field{
		zap.String("context_key", "context_value"),
		zap.Int("count", 42),
	}

	contextLogger := logger.Context(fields...)
	if contextLogger == nil {
		t.Error("Context() returned nil")
	}
}

func TestLoggingLevels(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	logger, err := New(Config{
		Level:      "debug",
		OutputPath: tmpFile.Name(),
		Format:     "json",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Test different log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")

	// Sync to ensure all logs are written
	logger.Sync()

	// Read log file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "debug message") {
		t.Error("Debug message not found in logs")
	}
	if !strings.Contains(logContent, "info message") {
		t.Error("Info message not found in logs")
	}
	if !strings.Contains(logContent, "warn message") {
		t.Error("Warn message not found in logs")
	}
}

func TestInfoWithFields(t *testing.T) {
	var buf bytes.Buffer
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	logger, err := New(Config{
		Level:      "info",
		OutputPath: tmpFile.Name(),
		Format:     "json",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	fields := map[string]interface{}{
		"user_id":   "123",
		"action":    "login",
		"timestamp": "2024-01-01T12:00:00Z",
	}

	logger.InfoWithFields("User action", fields)
	logger.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Parse JSON log
	var logEntry map[string]interface{}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) > 0 {
		if err := json.Unmarshal([]byte(lines[0]), &logEntry); err != nil {
			t.Fatalf("Failed to parse JSON log: %v", err)
		}

		if logEntry["message"] != "User action" {
			t.Errorf("Expected message 'User action', got '%v'", logEntry["message"])
		}
		if logEntry["user_id"] != "123" {
			t.Errorf("Expected user_id '123', got '%v'", logEntry["user_id"])
		}
	}

	_ = buf
}

func TestErrorWithFields(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	logger, err := New(Config{
		Level:      "error",
		OutputPath: tmpFile.Name(),
		Format:     "json",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	fields := map[string]interface{}{
		"error_code": 500,
		"endpoint":   "/api/test",
	}

	logger.ErrorWithFields("Request failed", fields)
	logger.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "Request failed") {
		t.Error("Error message not found in logs")
	}
	if !strings.Contains(string(content), "error_code") {
		t.Error("error_code field not found in logs")
	}
}

func TestWarnWithFields(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	logger, err := New(Config{
		Level:      "warn",
		OutputPath: tmpFile.Name(),
		Format:     "json",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	fields := map[string]interface{}{
		"threshold": 80,
		"current":   85,
	}

	logger.WarnWithFields("Memory usage high", fields)
	logger.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "Memory usage high") {
		t.Error("Warning message not found in logs")
	}
}

func TestDebugWithFields(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	logger, err := New(Config{
		Level:      "debug",
		OutputPath: tmpFile.Name(),
		Format:     "json",
	})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	fields := map[string]interface{}{
		"function": "TestDebugWithFields",
		"line":     123,
	}

	logger.DebugWithFields("Debug info", fields)
	logger.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "Debug info") {
		t.Error("Debug message not found in logs")
	}
}
