package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type Loggerer struct {
	logger *zap.Logger
}

func NewLogger(level zapcore.Level) (Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	log, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Loggerer{logger: log}, nil
}
func (l *Loggerer) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}
func (l *Loggerer) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}
func (l *Loggerer) Warn(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}
func (l *Loggerer) Error(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}
func (l *Loggerer) Fatal(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}
