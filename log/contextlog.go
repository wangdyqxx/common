package log

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/wangdyqxx/common/tracer"
)

var (
	emptyTrace = contextKV{tracer.ContextKeyTraceID: jaeger.TraceID{High: 0, Low: 0}}
	emptyHead  = contextKV{tracer.ContextKeyHeadUid: int64(0)}
)

var errorTraceIDNotFound = errors.New("traceID not found")
var errorHeadKVNotFound = errors.New("valid context head not found")

type contextKV map[string]interface{}

func newContextKV() contextKV {
	return contextKV{}
}

func (ckv contextKV) String() string {
	if v, ok := ckv[tracer.ContextKeyTraceID]; ok {
		return fmt.Sprintf("%v", v)
	}

	var parts []string
	if v, ok := ckv[tracer.ContextKeyHeadUid]; ok {
		if uid, uok := v.(int64); uok {
			parts = append(parts, fmt.Sprintf("%d", uid))
		}
	}

	var restParts []string
	for k, v := range ckv {
		if k != tracer.ContextKeyHeadUid && k != tracer.ContextKeyTraceID {
			restParts = append(restParts, fmt.Sprintf("%s:%v", k, v))
		}
	}
	if len(restParts) > 0 {
		parts = append(parts, strings.Join(restParts, " "))
		return strings.Join(parts, "\t")
	} else {
		return strings.Join(parts, "\t") + "\t"
	}
}

func extractTraceID(ctx context.Context) (error, contextKV) {
	ckv := newContextKV()
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			ckv[tracer.ContextKeyTraceID] = sc.TraceID()
			return nil, ckv
		}
	}
	return errorTraceIDNotFound, nil
}

func extractHead(ctx context.Context, fullHead bool) (error, contextKV) {
	head := ctx.Value(tracer.ContextKeyHead)
	if chd, ok := head.(tracer.ContextHeader); ok {
		kv := chd.ToKV()
		if fullHead {
			return nil, contextKV(chd.ToKV())
		}
		return nil, contextKV(map[string]interface{}{tracer.ContextKeyHeadUid: kv[tracer.ContextKeyHeadUid]})
	}
	return errorHeadKVNotFound, nil
}

func extractContext(ctx context.Context, fullHead bool) (v []interface{}) {
	if ctx == nil {
		return
	}

	if err, ckv := extractTraceID(ctx); err == nil {
		v = append(v, ckv)
	} else {
		v = append(v, emptyTrace)
	}

	if err, ckv := extractHead(ctx, fullHead); err == nil {
		v = append(v, ckv)
	} else {
		v = append(v, emptyHead)
	}

	return
}

func extractContextAsString(ctx context.Context, fullHead bool) (s string) {
	var parts []string
	for _, kv := range extractContext(ctx, fullHead) {
		parts = append(parts, fmt.Sprint(kv))
	}
	return strings.Join(parts, "  ")
}
