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
        // Use JSON encoder for production
        config.Encoding = "json"
        config.DisableStacktrace = true
    } else {
        config = zap.NewDevelopmentConfig()
        // Use console encoder with better formatting for development
        config.Encoding = "console"
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
        config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
        config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
        config.DisableStacktrace = true
    }
    
    // Customize output format
    config.DisableCaller = false
    config.EncoderConfig.CallerKey = "caller"
    config.EncoderConfig.MessageKey = "msg"
    config.EncoderConfig.LevelKey = "level"
    config.EncoderConfig.TimeKey = "time"
    
    // Improve readability
    config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
    
    // Create logger with appropriate skip level
    logger, _ := config.Build(zap.AddCallerSkip(1))
    
    // Conditionally enable stack traces only for errors and fatal logs
    if environment != "production" {
        logger = logger.WithOptions(
            zap.AddStacktrace(zapcore.ErrorLevel),
        )
    }
    
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