package contract

import (
	"context"
)

const TraceIDKey = "trace_id"

type (
	traceIDKey struct{}
	userIDKey  struct{}
)

// NewTraceIDContext 创建跟踪ID上下文
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// FromTraceIDContext 从上下文中获取跟踪ID
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewUserIDContext 创建用户ID上下文
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// FromUserIDContext 从上下文中获取用户ID
func FromUserIDContext(ctx context.Context) string {
	v := ctx.Value(userIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

//// WithContext Use context create entry
//func WithContext(ctx context.Context) *Entry {
//	if ctx == nil {
//		ctx = context.Background()
//	}
//
//	fields := map[string]interface{}{}
//
//	if v := FromTraceIDContext(ctx); v != "" {
//		fields[TraceIDKey] = v
//	}
//
//	return logrus.WithContext(ctx).WithFields(fields)
//}
