package common

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// opentelemetry instrumentation
// please see: https://opentelemetry.io/docs/languages/go/instrumentation/
//             https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#general-sdk-configuration

func NewOtelExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(ctx)
}

func NewOtelResource(serviceInstanceId string) (*sdkresource.Resource, error) {
	return sdkresource.Merge(
		sdkresource.Default(),
		sdkresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(Program),
			semconv.ServiceVersion(Version),
			semconv.ServiceNamespace(Package),
			semconv.ServiceInstanceID(serviceInstanceId),
		),
	)
}

func NewOtelTraceProvider(exporter sdktrace.SpanExporter, resource *sdkresource.Resource) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
}

func NewOtelDefaultTraceProvider(ctx context.Context, serviceInstanceId string) *sdktrace.TracerProvider {
	exporter, _ := NewOtelExporter(ctx)
	resource, _ := NewOtelResource(serviceInstanceId)
	return NewOtelTraceProvider(exporter, resource)
}
