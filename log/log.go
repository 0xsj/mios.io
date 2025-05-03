// old log.go
package log

type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type EmptyLogger struct{}

func (l *EmptyLogger) Info(args ...any)                  {}
func (l *EmptyLogger) Infof(format string, args ...any)  {}
func (l *EmptyLogger) Debug(args ...any)                 {}
func (l *EmptyLogger) Debugf(format string, args ...any) {}
func (l *EmptyLogger) Error(args ...any)                 {}
func (l *EmptyLogger) Errorf(format string, args ...any) {}
func (l *EmptyLogger) Fatal(args ...any)                 {}
func (l *EmptyLogger) Fatalf(format string, args ...any) {}
func (l *EmptyLogger) Warn(args ...any)                  {}
func (l *EmptyLogger) Warnf(format string, args ...any)  {}
