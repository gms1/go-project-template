package common

import (
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitRuntimeNoEcs(t *testing.T) {
	givenMaxProcs := runtime.GOMAXPROCS(0)
	givenMemLimit := "500MiB"
	t.Setenv("GOMAXPROCS", strconv.Itoa(givenMaxProcs))
	t.Setenv("GOMEMLIMIT", givenMemLimit)
	t.Setenv("ECS_CONTAINER_METADATA_URI_V4", "")
	InitRuntime()
	assert.Equal(t, givenMaxProcs, runtime.GOMAXPROCS(0))
	assert.Equal(t, givenMemLimit, os.Getenv("GOMEMLIMIT"))
}

func TestInitRuntimeEcs(t *testing.T) {
	givenMaxProcs := runtime.GOMAXPROCS(0)
	givenMemLimit := "500MiB"
	t.Setenv("GOMAXPROCS", strconv.Itoa(givenMaxProcs))
	t.Setenv("GOMEMLIMIT", givenMemLimit)
	t.Setenv("ECS_CONTAINER_METADATA_URI_V4", "xxx")
	InitRuntime()
	assert.Equal(t, givenMaxProcs, runtime.GOMAXPROCS(0))
	assert.Equal(t, givenMemLimit, os.Getenv("GOMEMLIMIT"))
}
