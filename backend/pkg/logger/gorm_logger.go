package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger is a custom GORM logger that uses Zap for structured logging
type GormLogger struct {
	ZapLogger                 *Logger
	SlowThreshold             time.Duration
	SkipErrRecordNotFound     bool
	ParameterizedQueries      bool
	LogLevel                  gormlogger.LogLevel
}

// NewGormLogger creates a new GORM logger that uses Zap
func NewGormLogger(zapLogger *Logger) *GormLogger {
	return &GormLogger{
		ZapLogger:             zapLogger,
		SlowThreshold:         200 * time.Millisecond,
		SkipErrRecordNotFound: true,
		ParameterizedQueries:  false,
		LogLevel:              gormlogger.Info,
	}
}

// LogMode sets the log level
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Info(fmt.Sprintf(msg, data...),
			zap.String("component", "gorm"),
		)
	}
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Warn(fmt.Sprintf(msg, data...),
			zap.String("component", "gorm"),
		)
	}
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.ZapLogger.Error(fmt.Sprintf(msg, data...),
			zap.String("component", "gorm"),
		)
	}
}

// Trace logs SQL queries with execution time
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("component", "gorm"),
		zap.Duration("duration", elapsed),
		zap.Int64("rows", rows),
	}

	// Only log SQL if not using parameterized queries
	if !l.ParameterizedQueries {
		fields = append(fields, zap.String("sql", sql))
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound):
		fields = append(fields, zap.Error(err))
		l.ZapLogger.Error("SQL query error", fields...)

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		fields = append(fields, zap.Duration("slow_threshold", l.SlowThreshold))
		l.ZapLogger.Warn("Slow SQL query", fields...)

	case l.LogLevel >= gormlogger.Info:
		l.ZapLogger.Debug("SQL query executed", fields...)
	}
}
