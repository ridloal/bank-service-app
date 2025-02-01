package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger() *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "" // Disable stacktrace

	// Create custom level enabler
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Create console encoder
	consoleEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)

	// Create console core
	consoleCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), lowPriority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), highPriority),
	)

	logger := zap.New(consoleCore,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return &Logger{
		Logger: logger,
	}
}

// Methods untuk logging dengan context
func (l *Logger) InfoWithContext(message string, fields ...zapcore.Field) {
	l.Info(message, fields...)
}

func (l *Logger) ErrorWithContext(message string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.Error(message, fields...)
}

func (l *Logger) WarnWithContext(message string, fields ...zapcore.Field) {
	l.Warn(message, fields...)
}

func (l *Logger) DebugWithContext(message string, fields ...zapcore.Field) {
	l.Debug(message, fields...)
}
