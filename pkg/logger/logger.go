package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

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

// DPanic logs a message at DPanic level
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

// Panic logs a message at Panic level
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

// Fatal logs a message at Fatal level
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Debugf logs a message at Debug level with a template
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Infof logs a message at Info level with a template
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warnf logs a message at Warn level with a template
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Errorf logs a message at Error level with a template
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// DPanicf logs a message at DPanic level with a template
func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.sugar.DPanicf(template, args...)
}

// Panicf logs a message at Panic level with a template
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

// Fatalf logs a message at Fatal level with a template
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}
