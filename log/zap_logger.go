// log/zap_logger.go
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
    logger *zap.SugaredLogger
    layer  string // Added layer field for context
}

func NewZapLogger(environment string) Logger {
    var config zap.Config
    if environment == "production" {
        config = zap.NewProductionConfig()
    } else {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    // Enable caller information
    config.DisableCaller = false
    config.EncoderConfig.CallerKey = "caller"
    config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
    
    logger, _ := config.Build(zap.AddCallerSkip(1)) // Skip one level to show actual caller
    return &ZapLogger{
        logger: logger.Sugar(),
    }
}

// WithLayer creates a new logger with layer information
func (l *ZapLogger) WithLayer(layer string) Logger {
    return &ZapLogger{
        logger: l.logger.With("layer", layer),
        layer:  layer,
    }
}

// WithField adds a custom field to the logger
func (l *ZapLogger) WithField(key string, value interface{}) Logger {
    return &ZapLogger{
        logger: l.logger.With(key, value),
        layer:  l.layer,
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