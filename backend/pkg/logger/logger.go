package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
	config Config
}

// Config holds logger configuration
type Config struct {
	Level      string // debug, info, warn, error
	OutputPath string // stdout, stderr, or file path
	Format     string // json or console
	EnableFile bool   // Enable file logging in addition to console
}

// New creates a new logger instance
func New(cfg Config) (*Logger, error) {
	// Parse log level
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "message"

	// Create encoder
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output
	var writer zapcore.WriteSyncer
	if cfg.OutputPath == "" || cfg.OutputPath == "stdout" {
		writer = zapcore.AddSync(os.Stdout)
	} else if cfg.OutputPath == "stderr" {
		writer = zapcore.AddSync(os.Stderr)
	} else {
		file, err := os.OpenFile(cfg.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writer = zapcore.AddSync(file)
	}

	// Create core
	core := zapcore.NewCore(encoder, writer, level)

	// Create logger with caller information
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		Logger: zapLogger,
		config: cfg,
	}, nil
}

// NewDefault creates a logger with default configuration
func NewDefault() *Logger {
	logger, _ := New(Config{
		Level:      "info",
		OutputPath: "stdout",
		Format:     "console",
	})
	return logger
}

// NewProduction creates a production logger (JSON format, info level)
func NewProduction() (*Logger, error) {
	return New(Config{
		Level:      "info",
		OutputPath: "stdout",
		Format:     "json",
	})
}

// NewDevelopment creates a development logger (console format, debug level)
func NewDevelopment() (*Logger, error) {
	return New(Config{
		Level:      "debug",
		OutputPath: "stdout",
		Format:     "console",
	})
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.With(zap.Any(key, value)),
		config: l.config,
	}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{
		Logger: l.With(zapFields...),
		config: l.config,
	}
}

// WithError adds error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.With(zap.Error(err)),
		config: l.config,
	}
}

// Context returns a logger with context fields
func (l *Logger) Context(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.With(fields...),
		config: l.config,
	}
}

// InfoWithFields logs an info message with structured fields
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Info(msg, zapFields...)
}

// ErrorWithFields logs an error message with structured fields
func (l *Logger) ErrorWithFields(msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Error(msg, zapFields...)
}

// WarnWithFields logs a warning message with structured fields
func (l *Logger) WarnWithFields(msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Warn(msg, zapFields...)
}

// DebugWithFields logs a debug message with structured fields
func (l *Logger) DebugWithFields(msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.Debug(msg, zapFields...)
}
