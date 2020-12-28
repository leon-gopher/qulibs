package gorm

import (
	"github.com/opentracing/opentracing-go"
	gogorm "gorm.io/gorm"
)

const (
	OpentracingContextKey     = "gorm:opentracing_context"
	OpentracingSpanContextKey = "gorm:opentracing_spanner"
)

type TraceClient struct {
	*gogorm.DB

	trace *TraceContext
}

type TraceContext struct {
	tracer  opentracing.Tracer
	spanCtx opentracing.SpanContext
}

func (ctx *TraceContext) StartSpan(name string) opentracing.Span {
	if ctx.spanCtx == nil {
		return ctx.tracer.StartSpan(name)
	}

	return ctx.tracer.StartSpan(name, opentracing.ChildOf(ctx.spanCtx))
}
