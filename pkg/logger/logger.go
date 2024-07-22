package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

type ContextKey string

var (
	contextKeys = map[ContextKey]string{
		TraceIDKey: "trace_id",
	}
	TraceIDKey = ContextKey("trace_id")
)

// NewLogger creates a new Logger instance
func NewLogger(mode string) (*Logger, error) {
	var config zap.Config
	if mode == "prod" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger: zLogger,
		sugar:  zLogger.Sugar(),
	}, nil
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() {
	err := l.logger.Sync()
	if err != nil {
		l.logger.Error("failed to sync logger", zap.Error(err))
	}
	err = l.sugar.Sync()
	if err != nil {
		l.logger.Error("failed to sync logger", zap.Error(err))
	}
}

// Debug logs a message at Debug level
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs a message at Info level
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a message at Warn level
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs a message at Error level
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a message at Fatal level
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	logger := l.logger
	var fields []zap.Field

	for key, value := range contextKeys {
		if val, ok := ctx.Value(key).(string); ok {
			fields = append(fields, zap.String(value, val))
		}
	}

	return &Logger{
		logger: logger.With(fields...),
	}
}
