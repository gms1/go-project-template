package core

import (
	"context"
	"errors"
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestRunServiceOk(t *testing.T) {
	defer AssertNoSignalHandler(t)
	_, found := os.LookupEnv("OTEL_SDK_DISABLED")
	if !found {
		t.Setenv("OTEL_SDK_DISABLED", "true")
	}

	assert.NoError(t, RunService(
		context.Background(),
		func(ctx context.Context, cancel context.CancelFunc, span trace.Span) error {
			return nil
		},
		nil,
		t.Name(),
	))
}

func TestRunServiceErrors(t *testing.T) {
	defer AssertNoSignalHandler(t)
	_, found := os.LookupEnv("OTEL_SDK_DISABLED")
	if !found {
		t.Setenv("OTEL_SDK_DISABLED", "true")
	}

	givenError := errors.New("test main error")

	assert.Equal(t, givenError, ErrorRootCause(RunService(
		context.Background(),
		func(ctx context.Context, cancel context.CancelFunc, span trace.Span) error {
			return givenError
		},
		nil,
		t.Name(),
	)))
}

func TestRunServicePanics(t *testing.T) {
	defer AssertNoSignalHandler(t)
	_, found := os.LookupEnv("OTEL_SDK_DISABLED")
	if !found {
		t.Setenv("OTEL_SDK_DISABLED", "true")
	}

	givenErrorText := "test main panic"
	err := RunService(
		context.Background(),
		func(ctx context.Context, cancel context.CancelFunc, span trace.Span) error {
			panic(givenErrorText)
		},
		nil,
		t.Name(),
	)
	assert.Contains(t, givenErrorText, err.Error())
}

func TestInitServiceRuntime(t *testing.T) {
	defer AssertNoSignalHandler(t)
	testCases := []struct {
		name    string
		withEcs bool
	}{
		{"with ECS", true},
		{"without ECS", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			givenMaxProcs := runtime.GOMAXPROCS(0)
			givenMemLimit := "500MiB"
			t.Setenv("GOMAXPROCS", strconv.Itoa(givenMaxProcs))
			t.Setenv("GOMEMLIMIT", givenMemLimit)
			if testCase.withEcs {
				t.Setenv("ECS_CONTAINER_METADATA_URI_V4", "xxx")
			} else {
				t.Setenv("ECS_CONTAINER_METADATA_URI_V4", "")
			}
			InitServiceRuntime(context.Background())
			assert.Equal(t, givenMaxProcs, runtime.GOMAXPROCS(0))
			assert.Equal(t, givenMemLimit, os.Getenv("GOMEMLIMIT"))
		})
	}
}
