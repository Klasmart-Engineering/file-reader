package log

import (
	"context"
	"fmt"

	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"
)

const (
	traceIDKey    = "trace.id"
	spanIDKey     = "span.id"
	entityGUIDKey = "entity.guid"
	entityNameKey = "entity.name"
	entityTypeKey = "entity.type"
	hostnameKey   = "hostname"
)

type ZapLogger struct {
	zap *zap.Logger
}

func (l *ZapLogger) WithCtxValue(ctx context.Context) *zap.Logger {
	return l.zap.With(l.keyAndValueFromContext(ctx)...)
}

func (l *ZapLogger) keyAndValueFromContext(ctx context.Context) []zap.Field {
	if txn := newrelic.FromContext(ctx); nil != txn {
		metadata := txn.GetLinkingMetadata()
		return []zap.Field{
			zap.String(traceIDKey, metadata.TraceID),
			zap.String(spanIDKey, metadata.SpanID),
			zap.String(entityGUIDKey, metadata.EntityGUID),
			zap.String(entityNameKey, metadata.EntityName),
			zap.String(entityTypeKey, metadata.EntityType),
			zap.String(hostnameKey, metadata.Hostname),
		}
	}
	return nil

}

func (l *ZapLogger) Named(name string) *zap.Logger {
	return l.zap.Named(name)
}

func Wrap(l *zap.Logger) *ZapLogger {
	return &ZapLogger{zap: l}
}

func (l *ZapLogger) Error(ctx context.Context, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Error(fmt.Sprint(i...))
}

func (l *ZapLogger) Errorf(ctx context.Context, s string, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Error(fmt.Sprintf(s, i...))
}

func (l *ZapLogger) Fatal(ctx context.Context, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Fatal(fmt.Sprint(i...))
}

func (l *ZapLogger) Fatalf(ctx context.Context, s string, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Fatal(fmt.Sprintf(s, i...))
}

func (l *ZapLogger) Debug(ctx context.Context, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Debug(fmt.Sprint(i...))
}

func (l *ZapLogger) Debugf(ctx context.Context, s string, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Debugf(fmt.Sprintf(s, i...))
}

func (l *ZapLogger) Info(ctx context.Context, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Info(fmt.Sprint(i...))
}

func (l *ZapLogger) Infof(ctx context.Context, s string, i ...interface{}) {
	l.WithCtxValue(ctx).Sugar().Info(fmt.Sprintf(s, i...))
}

func (l *ZapLogger) Printf(s string, i ...interface{}) {
	// Kafka-go writers need a Printf on loggers
	l.zap.Info(fmt.Sprintf(s, i...))
}

func (l *ZapLogger) Warn(i ...interface{}) {
	l.zap.Sugar().Warn(fmt.Sprint(i...))
}

func (l *ZapLogger) Warnf(s string, i ...interface{}) {
	l.zap.Sugar().Warn(fmt.Sprintf(s, i...))
}
