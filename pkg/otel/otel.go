package otel

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type OTel struct {
	tp            *sdktrace.TracerProvider
	shutdownFuncs []func(context.Context) error
}

func NewOTel(ctx context.Context, endpoint string) (*OTel, error) {
	shutdownFuncs := []func(context.Context) error{}

	exp, err := newOtlpHttpExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	tp, err := newTracerProvider(exp)
	if err != nil {
		return nil, err
	}

	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)

	otel.SetTracerProvider(tp)

	return &OTel{
		tp:            tp,
		shutdownFuncs: shutdownFuncs,
	}, nil
}

func (o *OTel) TracerProvider() *sdktrace.TracerProvider {
	return o.tp
}

func (o *OTel) Shutdown(ctx context.Context) (err error) {
	for _, f := range o.shutdownFuncs {
		err = errors.Join(err, f(ctx))
	}
	o.shutdownFuncs = nil
	return
}

func newOtlpHttpExporter(ctx context.Context, endpoint string) (*otlptrace.Exporter, error) {
	return otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(endpoint),
	)
}

func newTracerProvider(exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("user-service"),
		),
	)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	), nil
}
