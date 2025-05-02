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

func (l *ZapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}