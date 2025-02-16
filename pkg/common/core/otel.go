package core

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// opentelemetry instrumentation
// please see: https://opentelemetry.io/docs/languages/go/instrumentation/
//             https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#general-sdk-configuration

var (
	ServiceInstanceId = "Default"
	SpanName          = "RunService"
)

func NewOtelExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	if os.Getenv("OTEL_SDK_DISABLED") == "true" {
		slog.DebugContext(ctx, "opentelemetry is disabled")
		return tracetest.NewNoopExporter(), nil
	}
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
	// NOTE: https://opentelemetry.io/docs/concepts/instrumentation/libraries/#error-handling
	// we can ignore the errors returned, since we don't even get one
	exporter, _ := NewOtelExporter(ctx)
	resource, _ := NewOtelResource(serviceInstanceId)
	return NewOtelTraceProvider(exporter, resource)
}
