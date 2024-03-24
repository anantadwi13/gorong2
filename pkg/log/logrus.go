package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func InitLogrusLogger(l *logrus.Logger) {
	Init(&LogrusLogger{l})
}

func (l *LogrusLogger) getLogger(ctx context.Context, level Level) logrus.Ext1FieldLogger {
	if ctx == nil {
		return l.logger
	}
	switch level {
	case TraceLevel, DebugLevel:
		// handle fields
		return l.logger.WithFields(nil)
	default:
		return l.logger
	}
}

func (l *LogrusLogger) Trace(ctx context.Context, args ...interface{}) {
	l.getLogger(ctx, TraceLevel).Trace(args)
}

func (l *LogrusLogger) Debug(ctx context.Context, args ...interface{}) {
	l.getLogger(ctx, DebugLevel).Debug(args)
}

func (l *LogrusLogger) Info(ctx context.Context, args ...interface{}) {
	l.getLogger(ctx, InfoLevel).Info(args)
}

func (l *LogrusLogger) Warn(ctx context.Context, args ...interface{}) {
	l.getLogger(ctx, WarnLevel).Warn(args)
}

func (l *LogrusLogger) Error(ctx context.Context, args ...interface{}) {
	l.getLogger(ctx, ErrorLevel).Error(args)
}

func (l *LogrusLogger) Panic(ctx context.Context, args ...interface{}) {
	l.getLogger(ctx, PanicLevel).Panic(args)
}

func (l *LogrusLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.getLogger(ctx, FatalLevel).Fatal(args)
}

func (l *LogrusLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx, TraceLevel).Tracef(format, args)
}

func (l *LogrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx, DebugLevel).Debugf(format, args)
}

func (l *LogrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx, InfoLevel).Infof(format, args)
}

func (l *LogrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx, WarnLevel).Warnf(format, args)
}

func (l *LogrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx, ErrorLevel).Errorf(format, args)
}

func (l *LogrusLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx, PanicLevel).Panicf(format, args)
}

func (l *LogrusLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.getLogger(ctx, FatalLevel).Fatalf(format, args)
}
