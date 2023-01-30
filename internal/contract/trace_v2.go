package contract

/*


import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"time"
)

const TraceKey = "traceId"

func NewTraceId(tag string) string {
	now := time.Now()
	return fmt.Sprintf("%d.%d.%s", now.Unix(), now.Nanosecond(), tag)
}

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	// 从Context里面取
	traceInfo := GetTraceIdFromContext(ctx)

	if traceInfo == "" {
		traceInfo = GetTraceIdFromGRPCMeta(ctx)
	}

	return traceInfo
}

func GetTraceIdFromGRPCMeta(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if traceHeader, inMap := md[meta.TraceIdKey]; inMap {
			return traceHeader[0]
		}
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if traceHeader, inMap := md[meta.TraceIdKey]; inMap {
			return traceHeader[0]
		}
	}

	return ""
}

func GetTraceIdFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	traceId, ok := ctx.Value(TraceKey).(string)

	if !ok {
		return ""
	}

	return traceId
}

func SetTraceIdToContext(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, TraceKey, traceId)
}
*/
