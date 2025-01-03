package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// go.opentelemetry.io/otel/sdk/trace/tracetest

func TestNewOtelExporter(t *testing.T) {
	ctx := context.Background()
	_, err := NewOtelExporter(ctx)
	assert.NoError(t, err)
}

func TestNewOtelResource(t *testing.T) {
	_, err := NewOtelResource("")
	assert.NoError(t, err)
}

func TestNewOtelTraceProvider(t *testing.T) {
	ctx := context.Background()
	exporter := tracetest.NewNoopExporter()
	resource, _ := NewOtelResource("")
	traceProvider := NewOtelTraceProvider(exporter, resource)
	assert.NotNil(t, traceProvider)
	defer func() { _ = traceProvider.Shutdown(ctx) }()
}

func TestNewOtelDefaultTraceProvider(t *testing.T) {
	ctx := context.Background()
	traceProvider := NewOtelDefaultTraceProvider(ctx, "")
	assert.NotNil(t, traceProvider)
	defer func() { _ = traceProvider.Shutdown(ctx) }()
}
