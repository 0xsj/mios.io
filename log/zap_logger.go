// log/zap_logger.go
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(environment string) Logger {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, _ := config.Build()
	return &ZapLogger{
		logger: logger.Sugar(),
	}
}

func (l *ZapLogger) Info(args ...any) {
	l.logger.Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}

func (l *ZapLogger) Debug(args ...any) {
	l.logger.Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}

func (l *ZapLogger) Error(args ...any) {
	l.logger.Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}

func (l *ZapLogger) Fatal(args ...any) {
	l.logger.Fatal(args...)
}

func (l *ZapLogger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}

func (l *ZapLogger) Warn(args ...any) {
	l.logger.Warn(args...)
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	l.logger.Warnf(format, args...)
}
