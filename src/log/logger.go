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

func (l *ZapLogger) Error(ctx context.Context, newRelicEnabled bool, i ...interface{}) {
	if newRelicEnabled {
		l.WithCtxValue(ctx).Sugar().Error(fmt.Sprint(i...))
	} else {
		l.zap.Sugar().Error(fmt.Sprint(i...))
	}
}

func (l *ZapLogger) Errorf(ctx context.Context, newRelicEnabled bool, s string, i ...interface{}) {
	if newRelicEnabled {
		l.WithCtxValue(ctx).Sugar().Errorf(fmt.Sprint(i...))
	} else {
		l.zap.Sugar().Errorf(fmt.Sprint(i...))
	}
}

func (l *ZapLogger) Fatal(ctx context.Context, newRelicEnabled bool, i ...interface{}) {
	if newRelicEnabled {
		l.WithCtxValue(ctx).Fatal(fmt.Sprint(i...))
	} else {
		l.zap.Sugar().Fatal(fmt.Sprint(i...))
	}
}

func (l *ZapLogger) Fatalf(ctx context.Context, newRelicEnabled bool, s string, i ...interface{}) {
	if newRelicEnabled {
		l.WithCtxValue(ctx).Error(fmt.Sprintf(s, i...))
	} else {
		l.zap.Sugar().Fatalf(fmt.Sprintf(s, i...))
	}
}

func (l *ZapLogger) Info(ctx context.Context, newRelicEnabled bool, i ...interface{}) {
	if newRelicEnabled {
		l.WithCtxValue(ctx).Info(fmt.Sprint(i...))
	} else {
		l.zap.Sugar().Info(fmt.Sprint(i...))
	}
}

func (l *ZapLogger) Infof(ctx context.Context, newRelicEnabled bool, s string, i ...interface{}) {
	if newRelicEnabled {
		l.WithCtxValue(ctx).Info(fmt.Sprintf(s, i...))
	} else {
		l.zap.Sugar().Infof(fmt.Sprintf(s, i...))
	}
}

func (l *ZapLogger) Printf(ctx context.Context, newRelicEnabled bool, s string, i ...interface{}) {
	// Kafka writers need a Printf on loggers
	l.WithCtxValue(ctx).Info(fmt.Sprintf(s, i...))
}

func (l *ZapLogger) Warn(i ...interface{}) {
	l.zap.Warn(fmt.Sprint(i...))
}

func (l *ZapLogger) Warnf(s string, i ...interface{}) {
	l.zap.Warn(fmt.Sprintf(s, i...))
}
