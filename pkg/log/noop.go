package log

import (
	"context"
)

var noop = &noopLogger{}

type noopLogger struct {
}

func (n noopLogger) Trace(ctx context.Context, args ...interface{}) {}

func (n noopLogger) Debug(ctx context.Context, args ...interface{}) {}

func (n noopLogger) Info(ctx context.Context, args ...interface{}) {}

func (n noopLogger) Warn(ctx context.Context, args ...interface{}) {}

func (n noopLogger) Error(ctx context.Context, args ...interface{}) {}

func (n noopLogger) Panic(ctx context.Context, args ...interface{}) {}

func (n noopLogger) Fatal(ctx context.Context, args ...interface{}) {}

func (n noopLogger) Tracef(ctx context.Context, format string, args ...interface{}) {}

func (n noopLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}

func (n noopLogger) Infof(ctx context.Context, format string, args ...interface{}) {}

func (n noopLogger) Warnf(ctx context.Context, format string, args ...interface{}) {}

func (n noopLogger) Errorf(ctx context.Context, format string, args ...interface{}) {}

func (n noopLogger) Panicf(ctx context.Context, format string, args ...interface{}) {}

func (n noopLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {}
