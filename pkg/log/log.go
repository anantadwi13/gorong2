package log

import (
	"context"
	"sync"
)

type Level uint8

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger interface {
	Trace(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Panic(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})

	Tracef(ctx context.Context, format string, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
}

var (
	logger Logger
	mu     sync.RWMutex
	once   sync.Once
)

func Init(l Logger) {
	if l == nil {
		return
	}
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		logger = l
	})
}

func GetLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return logger
}

func getLogger() Logger {
	l := GetLogger()
	if l != nil {
		return l
	}
	return noop
}

func Trace(ctx context.Context, args ...interface{}) {
	getLogger().Trace(ctx, args)
}

func Debug(ctx context.Context, args ...interface{}) {
	getLogger().Debug(ctx, args)
}

func Info(ctx context.Context, args ...interface{}) {
	getLogger().Info(ctx, args)
}

func Warn(ctx context.Context, args ...interface{}) {
	getLogger().Warn(ctx, args)
}

func Error(ctx context.Context, args ...interface{}) {
	getLogger().Error(ctx, args)
}

func Panic(ctx context.Context, args ...interface{}) {
	getLogger().Panic(ctx, args)
}

func Fatal(ctx context.Context, args ...interface{}) {
	getLogger().Fatal(ctx, args)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	getLogger().Tracef(ctx, format, args)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLogger().Debugf(ctx, format, args)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	getLogger().Infof(ctx, format, args)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	getLogger().Warnf(ctx, format, args)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLogger().Errorf(ctx, format, args)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	getLogger().Panicf(ctx, format, args)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	getLogger().Fatalf(ctx, format, args)
}
